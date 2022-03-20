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

// integer tries to convert v into an integer
func integer(v reflect.Value) (int64, error) {
	k := v.Kind()

	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint()), nil
	}
	return 0, CannotCastToIntegerError
}

// findSizeForSlice solves the problem that since slice is variable length, so deserialize of slice needs a size to determine length to read.
func findSizeForSlice(f reflect.StructField, v reflect.Value) (int64, error) {
	// get size tag
	sizeVarName, _ := f.Tag.Lookup("size")
	if len(sizeVarName) > 0 {
		// try get size field and recover
		sizeF := v.FieldByName(sizeVarName)
		panic := recover()
		if panic != nil {
			return 0, SliceFiledSizeTagNotFound
		}
		return integer(sizeF)
	}
	return 0, SliceFiledSizeTagNotFound
}
