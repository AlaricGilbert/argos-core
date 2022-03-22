package bitcoin

import (
	"encoding/binary"

	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/cloudwego/netpoll"
)

// NetworkMagic is the type of magic value which indicates message origin network.
type NetworkMagic uint32

// VarInt is an integer stores as an uint64 but serializes & deserializes using variable length.
type VarInt uint64

// VarString is a string which contains VarInt and variable length of bytes string data
type VarString string

// InventoryType identifies the object type linked to this inventory
type InventoryType uint32

// ServiceType is set of bitfield of features to be enabled for some connection
type ServiceType uint64

// Serves checks the ServicesType contains the reference ServiceType (which should be set one of: NODE_NETWORK, NODE_GETUTXO, NODE_BLOOM, NODE_WITNESS, NODE_XTHIN, NODE_COMPACT_FILTERS or NODE_NETWORK_LIMITED)
func (s ServiceType) Serves(reference ServiceType) bool {
	return (uint64(s) & uint64(reference)) != 0
}

// Serialize serialize the given variable into writer
func (v *VarInt) Serialize(w netpoll.Writer, order binary.ByteOrder) (int, error) {
	i := uint64(*v)
	var err error
	var bytes int
	if i < 0xFD {
		return 1, w.WriteByte(byte(i))
	} else if i < 0xFFFF {
		if err = w.WriteByte(0xFD); err != nil {
			return 0, err
		}
		bytes, err = serialization.SerializeWithEndian(w, uint16(i), order)
		return bytes + 1, err
	} else if i < 0xFFFFFFFF {
		if err = w.WriteByte(0xFE); err != nil {
			return 0, err
		}
		bytes, err = serialization.SerializeWithEndian(w, uint32(i), order)
		return bytes + 1, err
	} else {
		if err = w.WriteByte(0xFF); err != nil {
			return 0, err
		}
		bytes, err = serialization.SerializeWithEndian(w, uint64(i), order)
		return bytes + 1, err
	}
}

func (str *VarString) Serialize(w netpoll.Writer, order binary.ByteOrder) (int, error) {
	var err error
	var bytes, n int
	if bytes, err = serialization.Serialize(w, VarInt(len(*str))); err != nil {
		return bytes, err
	}
	n, err = serialization.Serialize(w, []byte(*str))
	return bytes + n, err
}
