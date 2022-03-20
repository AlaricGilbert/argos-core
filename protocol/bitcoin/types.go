package bitcoin

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
