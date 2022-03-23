package bitcoin

const (
	MagicMain     NetworkMagic = 0xD9B4BEF9
	MagicTestnet  NetworkMagic = 0xDAB5BFFA
	MagicTestnet3 NetworkMagic = 0x0709110B
	MagicSignet   NetworkMagic = 0x40CF030A
	MagicNamecoin NetworkMagic = 0xFEB4BEF9
)

const (
	MessageHeaderLength = 24
)

var MagicSeeker = [][]byte{
	{0xF9, 0xBE, 0xB4, 0xD9},
	{0xFA, 0xBF, 0xB5, 0xDA},
	{0x0B, 0x11, 0x09, 0x07},
	{0x0A, 0x03, 0xCF, 0x40},
	{0xF9, 0xBE, 0xB4, 0xFE},
}

var magics = map[string]NetworkMagic{
	"main":     MagicMain,
	"testnet":  MagicTestnet,
	"regtest":  MagicTestnet,
	"testnet3": MagicTestnet3,
	"signet":   MagicSignet,
	"default":  MagicSignet,
	"namecoin": MagicNamecoin,
}

var magicNames = map[NetworkMagic]string{
	MagicMain:     "main",
	MagicTestnet:  "testnet",
	MagicTestnet3: "testnet3",
	MagicSignet:   "signet",
	MagicNamecoin: "namecoin",
}

const (
	MSG_VALIDATION_MASK InventoryType = 0b10111111111111111111111111111100
	// MSG_WITNESS_FLAG is NOT a valid InventoryType value but a FLAG to
	MSG_WITNESS_FLAG InventoryType = 1 << 30
	// MSG_TX means hash is related to a transaction
	MSG_TX InventoryType = 1
	// MSG_BLOCK means hash is related to a data block
	MSG_BLOCK InventoryType = 2
	// MSG_FILTERED_BLOCK means hash of a block header; identical to MSG_BLOCK. Only to be used in getdata message. Indicates the reply should be a merkleblock message rather than a block message; this only works if a bloom filter has been set. See BIP 37 for more info.
	MSG_FILTERED_BLOCK InventoryType = 3
	// MSG_CMPCT_BLOCK means hash of a block header; identical to MSG_BLOCK. Only to be used in getdata message. Indicates the reply should be a cmpctblock message. See BIP 152 for more info.
	MSG_CMPCT_BLOCK InventoryType = 4
	// MSG_WITNESS_TX means hash of a transaction with witness data. See BIP 144 for more info.
	MSG_WITNESS_TX InventoryType = MSG_TX | MSG_WITNESS_FLAG
	// MSG_WITNESS_BLOCK means hash of a block with witness data. See BIP 144 for more info.
	MSG_WITNESS_BLOCK InventoryType = MSG_BLOCK | MSG_WITNESS_FLAG
	// MSG_WITNESS_FILTERED_BLOCK means hash of a block with witness data. Only to be used in getdata message. Indicates the reply should be a merkleblock message rather than a block message; this only works if a bloom filter has been set. See BIP 144 for more info.
	MSG_WITNESS_FILTERED_BLOCK InventoryType = MSG_FILTERED_BLOCK | MSG_WITNESS_FLAG
)

const (
	NODE_NETWORK         ServiceType = 1 // This node can be asked for full blocks instead of just headers.
	NODE_GETUTXO         ServiceType = 2
	NODE_BLOOM           ServiceType = 4
	NODE_WITNESS         ServiceType = 8
	NODE_XTHIN           ServiceType = 16 // Never formally proposed (as a BIP), and discontinued. Was historically sporadically seen on the network.
	NODE_COMPACT_FILTERS ServiceType = 64
	NODE_NETWORK_LIMITED ServiceType = 1024
)

const (
	REJECT_MALFORMED       = 0x01
	REJECT_INVALID         = 0x10
	REJECT_OBSOLETE        = 0x11
	REJECT_DUPLICATE       = 0x12
	REJECT_NONSTANDARD     = 0x40
	REJECT_DUST            = 0x41
	REJECT_INSUFFICIENTFEE = 0x42
	REJECT_CHECKPOINT      = 0x43
)

const BitcoinMessageMaxLength = 4096 * 1024 // 4096 kilobytes

const (
	CommandReject  = "reject"
	CommandVersion = "version"
	CommandVerack  = "verack"
)

const (
	UserAgent = "/Argos:0.1/"
)
