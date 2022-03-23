package bitcoin

import (
	"encoding/binary"
	"fmt"

	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/cloudwego/netpoll"
)

// NetworkMagic is the type of magic value which indicates message origin network.
type NetworkMagic uint32

// String implements fmt.Stringer
func (nm NetworkMagic) String() string {
	if str, ok := magicNames[nm]; ok {
		return str
	}
	return "INVALID"
}

// VarInt is an integer stores as an uint64 but serializes & deserializes using variable length.
type VarInt uint64

// VarString is a string which contains VarInt and variable length of bytes string data
type VarString string

// InventoryType identifies the object type linked to this inventory
type InventoryType uint32

// String implements fmt.Stringer
func (i InventoryType) String() string {
	switch i {
	case MSG_TX:
		return "TX"
	case MSG_BLOCK:
		return "BLOCK"
	case MSG_FILTERED_BLOCK:
		return "FILTERED_BLOCK"
	case MSG_CMPCT_BLOCK:
		return "CMPCT_BLOCK"
	case MSG_WITNESS_TX:
		return "WITNESS_TX"
	case MSG_WITNESS_BLOCK:
		return "WITNESS_BLOCK"
	case MSG_WITNESS_FILTERED_BLOCK:
		return "WITNESS_FILTERED_BLOCK"
	default:
		return "INVALID"
	}
}

// Tx checks the inventory has TX flag
func (i InventoryType) Tx() bool {
	return i.Basic() == MSG_TX
}

// Block checks the inventory has Block flag
func (i InventoryType) Block() bool {
	return i.Basic() == MSG_BLOCK
}

// FilteredBlock checks the inventory has FilteredBlock flag
func (i InventoryType) FilteredBlock() bool {
	return i.Basic() == MSG_FILTERED_BLOCK
}

// CmpctBlock checks the inventory has CmpctBlock flag
func (i InventoryType) CmpctBlock() bool {
	return i.Basic() == MSG_CMPCT_BLOCK
}

// Basic gets inventory type without witness flag
func (i InventoryType) Basic() InventoryType {
	return i & (^MSG_WITNESS_FLAG)
}

// Witness checks whether the inventory type marks as containing witness data
func (i InventoryType) Witness() bool {
	return (i & MSG_WITNESS_FLAG) == MSG_WITNESS_FLAG
}

// Valid checks whether the inventory type valids
func (i InventoryType) Valid() bool {
	return ((i & MSG_VALIDATION_MASK) == 0) && !(i.Witness() && (i.Basic() == MSG_CMPCT_BLOCK))
}

// ServiceType is set of bitfield of features to be enabled for some connection
type ServiceType uint64

// String implements fmt.Stringer
func (s ServiceType) String() string {
	return fmt.Sprintf("{NODE_NETWORK: %t, NODE_GETUTXO: %t, NODE_BLOOM: %t, NODE_WITNESS: %t, NODE_XTHIN: %t, NODE_COMPACT_FILTERS: %t, NODE_NETWORK_LIMITED: %t}",
		s.Serves(NODE_NETWORK),
		s.Serves(NODE_GETUTXO),
		s.Serves(NODE_BLOOM),
		s.Serves(NODE_WITNESS),
		s.Serves(NODE_XTHIN),
		s.Serves(NODE_COMPACT_FILTERS),
		s.Serves(NODE_NETWORK_LIMITED),
	)
}

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
