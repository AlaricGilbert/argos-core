package serialization

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/cloudwego/netpoll"
)

// intDataSize returns the size of the data required to represent the data when encoded.
// It returns zero if the type cannot be implemented by the fast path in Read or Write.
func intDataSize(data any) int {
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case []bool:
		return len(data)
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, uint16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, uint32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, uint64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	case float32, *float32:
		return 4
	case float64, *float64:
		return 8
	case []float32:
		return 4 * len(data)
	case []float64:
		return 8 * len(data)
	}
	return 0
}

func Deserialize(r netpoll.Reader, data any) (int, error) {
	return DeserializeWithEndian(r, data, binary.LittleEndian)
}

func DeserializeWithEndian(r netpoll.Reader, data any, order binary.ByteOrder) (int, error) {
	// bytes read
	var bs []byte
	// error
	var err error
	// count of total bytes have been read
	var bytes = 0
	// tmp count of bytes have been read
	var n = 0

	// Fast path for basic types and slices.
	if bytes = intDataSize(data); bytes != 0 {
		if bs, err = r.Next(bytes); err != nil {
			return 0, err
		}
		switch data := data.(type) {
		case *bool:
			*data = bs[0] != 0
		case *int8:
			*data = int8(bs[0])
		case *uint8:
			*data = bs[0]
		case *int16:
			*data = int16(order.Uint16(bs))
		case *uint16:
			*data = order.Uint16(bs)
		case *int32:
			*data = int32(order.Uint32(bs))
		case *uint32:
			*data = order.Uint32(bs)
		case *int64:
			*data = int64(order.Uint64(bs))
		case *uint64:
			*data = order.Uint64(bs)
		case *float32:
			*data = math.Float32frombits(order.Uint32(bs))
		case *float64:
			*data = math.Float64frombits(order.Uint64(bs))
		case []bool:
			for i, x := range bs { // Easier to loop over the input for 8-bit values.
				data[i] = x != 0
			}
		case []int8:
			for i, x := range bs {
				data[i] = int8(x)
			}
		case []uint8:
			copy(data, bs)
		case []int16:
			for i := range data {
				data[i] = int16(order.Uint16(bs[2*i:]))
			}
		case []uint16:
			for i := range data {
				data[i] = order.Uint16(bs[2*i:])
			}
		case []int32:
			for i := range data {
				data[i] = int32(order.Uint32(bs[4*i:]))
			}
		case []uint32:
			for i := range data {
				data[i] = order.Uint32(bs[4*i:])
			}
		case []int64:
			for i := range data {
				data[i] = int64(order.Uint64(bs[8*i:]))
			}
		case []uint64:
			for i := range data {
				data[i] = order.Uint64(bs[8*i:])
			}
		case []float32:
			for i := range data {
				data[i] = math.Float32frombits(order.Uint32(bs[4*i:]))
			}
		case []float64:
			for i := range data {
				data[i] = math.Float64frombits(order.Uint64(bs[8*i:]))
			}
		default:
			bytes = 0 // fast path doesn't apply
		}
		if bytes != 0 {
			return bytes, nil
		}
	}

	// fall into customized type deserialize
	for _, d := range deserializers {
		if n, err = d(r, data, order); err != DeserializeTypeDismatchError {
			return bytes + n, err
		}
		bytes += n
	}

	// Fallback to reflect-based decoding.
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		l := v.Len()

		for i := 0; i < l; i++ {
			if n, err = DeserializeWithEndian(r, v.Index(i).Addr().Interface(), order); err != nil {
				return bytes + n, err
			}
			bytes += n
		}
		return bytes, nil
	case reflect.Struct:
		l := v.NumField()
		typ := v.Type()
		var size int64
		for i := 0; i < l; i++ {
			fieldTyp := typ.Field(i)
			fieldValue := v.Field(i)
			order := Order(fieldTyp)
			if v.Field(i).Kind() == reflect.Slice {
				size, err = findSizeForSlice(fieldTyp, v)
				if err != nil {
					return bytes, SliceFiledSizeTagNotFound
				}
				fieldValue.Set(reflect.MakeSlice(fieldTyp.Type, int(size), int(size)))
			}

			if Omit(fieldTyp) {
				if i != l-1 {
					return bytes, NonLastFieldContainsOmitOption
				}
				if r.Len() == 0 {
					return bytes, nil
				}
			}
			if n, err = DeserializeWithEndian(r, fieldValue.Addr().Interface(), order); err != nil {
				return bytes + n, err
			}
			bytes += n
		}
		return bytes, nil
	}
	return bytes, errors.New(fmt.Sprintf("read failed: invalid type %v", v.Type()))
}
