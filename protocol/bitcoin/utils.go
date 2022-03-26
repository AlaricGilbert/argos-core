package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"strings"
)

// hash returns the hash result of sha256(sha256(data))
func hash(data []byte) [32]byte {
	h := sha256.Sum256(data)
	hh := sha256.Sum256(h[:])

	return hh
}

func checksum(data []byte) ([32]byte, uint32) {
	hh := hash(data)
	return hh, binary.LittleEndian.Uint32(hh[:])
}

// Index finds the index of given element e in array s returns index and succeed result of the process
func Index[T comparable](s []T, e T) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return -1, false
}

// SliceToString converts a byte slice into a string ends with '\0'
func SliceToString(data []byte) string {
	n := bytes.IndexByte(data, 0)
	if n != -1 {
		data = data[:n]
	}
	return string(data)
}

func FmtSlice[T any](slice []T, dataFmt func(T) string) string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	length := len(slice)
	for i := 0; i < length; i++ {
		builder.WriteString(dataFmt(slice[i]))
		if i+1 != length {
			builder.WriteString(", ")
		}
	}
	builder.WriteRune(']')
	return builder.String()
}
