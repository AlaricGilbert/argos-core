package bitcoin

import (
	"bytes"
	"crypto/sha256"
)

// hash returns the hash result of sha256(sha256(data))
func hash(data []byte) []byte {
	h := sha256.New()
	hh := sha256.New()
	return hh.Sum(h.Sum(data))
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
