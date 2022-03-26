package bitcoin

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	h := hash([]byte("hello"))
	assert.Equal(t, "9595c9df90075148eb06860365df33584b75bff782a510c6cd4883a419833d50", hex.EncodeToString(h[:]))
}

func TestIndex(t *testing.T) {
	var arr = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, b := range arr {
		id, ok := Index(arr, b)
		assert.True(t, ok)
		assert.Equal(t, i, id)
	}
}

func TestSliceToString(t *testing.T) {
	var arr = []byte{'H', 'e', 'l', 'l', '0', ',', ' ', 'w', '0', 'r', 'l', 'd', '!', '\000'}
	assert.Equal(t, "Hell0, w0rld!", SliceToString(arr))
}
