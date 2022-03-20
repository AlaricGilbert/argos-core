package serialization

import (
	"encoding/binary"

	"github.com/cloudwego/netpoll"
)

type CustomDeserializer func(r netpoll.Reader, data any, order binary.ByteOrder) (int, error)

var deserializers []CustomDeserializer

func RegisterDeserializer(deserializer CustomDeserializer) {
	deserializers = append(deserializers, deserializer)
}
