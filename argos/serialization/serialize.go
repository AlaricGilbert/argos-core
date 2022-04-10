package serialization

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/cloudwego/netpoll"
)

// Serialize using builtin basic type memory binary representations and little endian order
// to serialize the given `data` object pointer into binary reading stream.
// It should be noticed the `data` must be a pointer when a slice's corresponding size field
// was not set to the slice's length.
func Serialize(w netpoll.Writer, data any) (int, error) {
	return SerializeWithEndian(w, data, binary.LittleEndian)
}

// Serialize using builtin basic type memory binary representations and given `order``
// to serialize the given `data` object pointer into binary reading stream.
// It should be noticed the `data` must be a pointer when a slice's corresponding size field
// was not set to the slice's length.
func SerializeWithEndian(w netpoll.Writer, data any, order binary.ByteOrder) (int, error) {
	if data == nil {
		return 0, nil
	}
	// error
	var err error
	// count of total bytes have been read
	var bytes int
	// tmp count of bytes have been read
	var n int

	// Fast path for basic types and slices.
	if n := intDataSize(data); n != 0 {
		bs := make([]byte, n)
		switch v := data.(type) {
		case *bool:
			if *v {
				bs[0] = 1
			} else {
				bs[0] = 0
			}
		case bool:
			if v {
				bs[0] = 1
			} else {
				bs[0] = 0
			}
		case []bool:
			for i, x := range v {
				if x {
					bs[i] = 1
				} else {
					bs[i] = 0
				}
			}
		case *int8:
			bs[0] = byte(*v)
		case int8:
			bs[0] = byte(v)
		case []int8:
			for i, x := range v {
				bs[i] = byte(x)
			}
		case *uint8:
			bs[0] = *v
		case uint8:
			bs[0] = v
		case []uint8:
			bs = v
		case *int16:
			order.PutUint16(bs, uint16(*v))
		case int16:
			order.PutUint16(bs, uint16(v))
		case []int16:
			for i, x := range v {
				order.PutUint16(bs[2*i:], uint16(x))
			}
		case *uint16:
			order.PutUint16(bs, *v)
		case uint16:
			order.PutUint16(bs, v)
		case []uint16:
			for i, x := range v {
				order.PutUint16(bs[2*i:], x)
			}
		case *int32:
			order.PutUint32(bs, uint32(*v))
		case int32:
			order.PutUint32(bs, uint32(v))
		case []int32:
			for i, x := range v {
				order.PutUint32(bs[4*i:], uint32(x))
			}
		case *uint32:
			order.PutUint32(bs, *v)
		case uint32:
			order.PutUint32(bs, v)
		case []uint32:
			for i, x := range v {
				order.PutUint32(bs[4*i:], x)
			}
		case *int64:
			order.PutUint64(bs, uint64(*v))
		case int64:
			order.PutUint64(bs, uint64(v))
		case []int64:
			for i, x := range v {
				order.PutUint64(bs[8*i:], uint64(x))
			}
		case *uint64:
			order.PutUint64(bs, *v)
		case uint64:
			order.PutUint64(bs, v)
		case []uint64:
			for i, x := range v {
				order.PutUint64(bs[8*i:], x)
			}
		case *float32:
			order.PutUint32(bs, math.Float32bits(*v))
		case float32:
			order.PutUint32(bs, math.Float32bits(v))
		case []float32:
			for i, x := range v {
				order.PutUint32(bs[4*i:], math.Float32bits(x))
			}
		case *float64:
			order.PutUint64(bs, math.Float64bits(*v))
		case float64:
			order.PutUint64(bs, math.Float64bits(v))
		case []float64:
			for i, x := range v {
				order.PutUint64(bs[8*i:], math.Float64bits(x))
			}
		}
		return w.WriteBinary(bs)
	}

	// fall into customized type deserialize
	for _, s := range serializers {
		if n, err = s.Serialize(w, data, order); err != ErrSerializeTypeDismatch {
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
			if n, err = SerializeWithEndian(w, v.Index(i).Addr().Interface(), order); err != nil {
				return bytes + n, err
			}
			bytes += n
		}
		return bytes, nil
	case reflect.Struct:
		l := v.NumField()
		typ := v.Type()

		// first run sets all the slice size variable
		for i := 0; i < l; i++ {
			fieldTyp := typ.Field(i)
			field := v.Field(i)

			if field.Kind() == reflect.Slice {
				sizeF, err := findSizeForSlice(fieldTyp, v)
				if err != nil {
					return bytes, err
				}
				len := field.Len()
				if sizeF.CanInt() {
					size := sizeF.Int()
					if size != int64(field.Len()) {
						if !sizeF.CanSet() {
							return bytes, ErrUsingUnaddressableValue
						}
						sizeF.SetInt(int64(len))
					}
				} else if sizeF.CanUint() {
					size := sizeF.Uint()
					if size != uint64(field.Len()) {
						if !sizeF.CanSet() {
							return bytes, ErrUsingUnaddressableValue
						}
						sizeF.SetUint(uint64(len))
					}
				} else {
					return bytes, ErrSliceFiledSizeTagNotFound
				}
			}
		}

		for i := 0; i < l; i++ {
			order := Order(typ.Field(i))
			field := v.Field(i)

			if n, err = SerializeWithEndian(w, field.Addr().Interface(), order); err != nil {
				return bytes + n, err
			}
			bytes += n
		}
		return bytes, nil
	}
	return bytes, errors.New(fmt.Sprintf("write failed: invalid type %v", v.Type()))
}
