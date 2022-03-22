package serialization

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/cloudwego/netpoll"
)

// Deserialize using builtin basic type memory binary representations and little endian order
// to deserialize the given binary reading stream into given `data` object pointer.
// It should be noticed the `data` must be a pointer because value pass will cause unaddressable
// panic
func Deserialize(r netpoll.Reader, data any) (int, error) {
	return DeserializeWithEndian(r, data, binary.LittleEndian)
}

// DeserializeWithEndian using builtin basic type memory binary representations and given `order`
// to deserialize the given binary reading stream into given `data` object pointer.
// It should be noticed the `data` must be a pointer because value pass will cause unaddressable
// panic
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
	for _, s := range serializers {
		if n, err = s.Deserialize(r, data, order); err != SerializeTypeDismatchError {
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
		var size int
		for i := 0; i < l; i++ {
			fieldTyp := typ.Field(i)
			fieldValue := v.Field(i)
			order := Order(fieldTyp)
			if v.Field(i).Kind() == reflect.Slice {

				if sizeF, err := findSizeForSlice(fieldTyp, v); err != nil {
					return bytes, err
				} else {
					if sizeF.CanInt() {
						size = int(sizeF.Int())
					} else if sizeF.CanUint() {
						size = int(sizeF.Uint())
					}
				}
				fieldValue.Set(reflect.MakeSlice(fieldTyp.Type, size, size))
			}

			if Omit(fieldTyp) {
				if i != l-1 {
					return bytes, NonLastFieldContainsOmitOptionError
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
