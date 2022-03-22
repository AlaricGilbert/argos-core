package serialization

import (
	"encoding/binary"

	"github.com/cloudwego/netpoll"
)

// Serializer describes a object can handle user-defined struct serialization behavours.
type Serializer interface {
	Serialize(w netpoll.Writer, data any, order binary.ByteOrder) (int, error)
	Deserialize(r netpoll.Reader, data any, order binary.ByteOrder) (int, error)
}

var serializers map[string]Serializer = make(map[string]Serializer)

// RegisterSerializer registers the Serializer with id.
// After registeration, the user-defined struct serialization will be called before recursivly
// reflect based serialize and deserialize.
func RegisterSerializer(id string, s Serializer) error {
	if _, ok := serializers[id]; ok {
		return SerializerAlreadyExistsError
	}
	serializers[id] = s
	return nil
}
