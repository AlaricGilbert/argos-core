package serialization

import (
	"encoding/binary"
	"reflect"
	"strings"
)

// Order finds the 'order' tag in given field and returns binary.BigEndian if value of 'order' set to "network" or "big".
func Order(f reflect.StructField) binary.ByteOrder {
	tag, _ := f.Tag.Lookup("order")
	if tag == "network" || tag == "big" {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func Omit(f reflect.StructField) bool {
	tag, _ := f.Tag.Lookup("deserialize")
	deserializeOpts := strings.Split(tag, ",")
	for _, v := range deserializeOpts {
		if strings.TrimSpace(v) == "omit" {
			return true
		}
	}
	return false
}

// findSizeForSlice solves the problem that since slice is variable length, so deserialize of slice needs a size to determine length to read.
func findSizeForSlice(f reflect.StructField, v reflect.Value) (reflect.Value, error) {
	// get size tag
	sizeVarName, _ := f.Tag.Lookup("size")
	if len(sizeVarName) > 0 {
		// try get size field and recover
		sizeF := v.FieldByName(sizeVarName)
		panic := recover()
		if panic != nil {
			return v, SliceFiledSizeTagNotFound
		}
		return sizeF, nil
	}
	return v, SliceFiledSizeTagNotFound
}

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
