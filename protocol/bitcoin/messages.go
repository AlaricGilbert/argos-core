package bitcoin

import (
	"net"
	"net/netip"
)

type MessageHeader struct {
	magic    NetworkMagic // Magic value indicating message origin network, and used to seek to next message when stream state is unknown
	command  [12]byte     // ASCII string identifying the packet content, NULL padded (non-NULL padding results in packet rejected)
	length   uint32       // Length of payload in number of bytes
	checksum uint32       // First 4 bytes of sha256(sha256(payload))
}

// NetworkAddress is used when a network address is needed somewhere
type NetworkAddress struct {
	Time     uint32   // the Time (version >= 31402). Not present in version message.
	Services uint64   // same service(s) listed in version
	IP       [16]byte // IPv6 address. Network byte order. The original client only supported IPv4 and only Read the last 4 bytes to get the IPv4 address. However, the IPv4 address is written into the message as a 16 byte IPv4-mapped IPv6 address
	Port     uint16   `order:"network"` // port number, network byte order
}

func (a *NetworkAddress) AddrPort() netip.AddrPort {
	return netip.AddrPortFrom(netip.AddrFrom16(a.IP), a.Port)
}

func (a *NetworkAddress) TCPAddr() *net.TCPAddr {
	return net.TCPAddrFromAddrPort(a.AddrPort())
}

// networkAddress is a substitution for NetworkAddress in Version message since NetworkAddress.Time was not present in version message
type networkAddress struct {
	Services uint64   // same service(s) listed in version
	IP       [16]byte // IPv6 address. Network byte order. The original client only supported IPv4 and only Read the last 4 bytes to get the IPv4 address. However, the IPv4 address is written into the message as a 16 byte IPv4-mapped IPv6 address
	Port     uint16   `order:"network"` // port number, network byte order
}

func (a *networkAddress) AddrPort() netip.AddrPort {
	return netip.AddrPortFrom(netip.AddrFrom16(a.IP), a.Port)
}

func (a *networkAddress) TCPAddr() *net.TCPAddr {
	return net.TCPAddrFromAddrPort(a.AddrPort())
}

// Inventory vectors are used for notifying other nodes about objects they have or data which is being requested.
type Inventory struct {
	Type InventoryType // Identifies the object type linked to this inventory
	Hash [32]byte      // Hash of the object
}

// Version will immediately be advertised when a node creates an outgoing connection.
// The remote node will respond with its Version.
// No further communication is possible until both peers have exchanged their Version.
type Version struct {
	Version      int32          // Identifies protocol version being used by the node
	Services     ServiceType    // Bitfield of features to be enabled for this connection
	Timestamp    int64          // Standard UNIX Timestamp in seconds
	AddrReceived networkAddress // The network address of the node receiving this message
	AddrFrom     networkAddress // Field can be ignored. This used to be the network address of the node emitting this message, but most P2P implementations send 26 dummy bytes. The "services" field of the address would also be redundant with the second field of the version message.
	Nonce        uint64         // Node random nonce, randomly generated every time a version packet is sent. This nonce is used to detect connections to self.
	UserAgent    VarString      // User Agent (0x00 if string is 0 bytes long)
	StartHeight  int32          // The last block received by the emitting node
	Relay        bool           // Whether the remote peer should announce relayed transactions or not, see BIP 0037
}

// Addr provides information on known nodes of the network. Non-advertised nodes should be forgotten after typically 3 hours
type Addr struct {
	Count    VarInt           // Number of address entries (max: 1000)
	AddrList []NetworkAddress `size:"Count"` // Address of other nodes on the network. version < 209 will only read the first one.
}

type BlockLocator struct {
	Version            uint32     // the protocol version
	HashCount          VarInt     // number of block locator hash entries
	BlockLocatorHashes [][32]byte `size:"HashCount"` // block locator object; newest back to genesis block (dense to start, but then sparse)
	HashStop           [32]byte   // hash of the last desired block; set to zero to get as many blocks as possible (500)
}

type OutPoint struct {
	Hash  [32]byte // The hash of the referenced transaction.
	Index uint32   // The index of the specific output in the transaction. The first output is 0, etc.
}

// TransactionIn means Transaction inputs or sources for coins
type TransactionIn struct {
	PreviousOutput  OutPoint // The previous output transaction reference, as an OutPoint structure
	ScriptLength    VarInt   // The length of the signature script
	SignatureSctipt []byte   `size:"ScriptLength"` // Computational Script for confirming transaction authorization
	Sequence        uint32   // Transaction version as defined by the sender. Intended for "replacement" of transactions when information is updated before inclusion into a block.
}

// TransactionOut means Transaction outputs or destinations for coins
type TransactionOut struct {
	Value          int64  // Transaction Value
	PKScriptLength VarInt // Length of the pk_script
	PKScript       []byte `size:"PKScriptLength"` // Usually contains the public key as a Bitcoin script setting up conditions to claim this output.
}

type TransactionWitnesss struct {
	WitnessCount VarInt
	WitnessData  []byte `size:"WitnessCount"`
}

// Transaction describes a bitcoin transaction, in reply to getdata. When a bloom filter is applied tx objects are sent automatically for matching transactions following the merkleblock.
type Transaction struct {
	Version        uint32                // Transaction data format version
	Flag           [2]uint8              // If present, always 0001, and indicates the presence of witness data
	TxInCount      VarInt                // Number of Transaction inputs (never zero)
	TxIn           []TransactionIn       `size:"TxInCount"` // A list of 1 or more transaction inputs or sources for coins
	TxOutCount     VarInt                // Number of Transaction outputs
	TxOut          []TransactionOut      `size:"TxOutCount"` // A list of 1 or more transaction outputs or destinations for coins
	TxWitnessCount VarInt                // Number of TxWitness
	TxWitness      []TransactionWitnesss `size:"TxWitnessCount"` // A list of witnesses, one for each input; omitted if flag is omitted above
	// The block number or timestamp at which this transaction is unlocked:
	// Value         Description
	// 0                Not locked
	// < 500000000      Block number at which this transaction is unlocked
	// >= 500000000     UNIX timestamp at which this transaction is unlocked
	// If all TxIn inputs have final (0xffffffff) sequence numbers then lock_time is irrelevant. Otherwise, the transaction may not be added to a block until after lock_time (see NLockTime).
	LockTime uint32
}

// Block message is sent in response to a getdata message which requests transaction information from a block hash.
type Block struct {
	Version    int32         // Block version information (note, this is signed)
	PrevBlock  [32]byte      // The hash value of the previous block this particular block references
	MerkleRoot [32]byte      // The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp  uint32        // A timestamp recording when this block was created (Will overflow in 2106[2])
	Bits       uint32        // The calculated difficulty target being used for this block
	Nonce      uint32        // The nonce used to generate this block… to allow variations of the header and compute different hashes
	TxCount    VarInt        // Number of transaction entries
	Txs        []Transaction `size:"TxCount"` // Block transactions, in format of "tx" command
}

// BlockHeader is sent in a header packet in response to a getheaders message.
type BlockHeader struct {
	Version          int32    // Block version information (note, this is signed)
	PrevBlock        [32]byte // The hash value of the previous block this particular block references
	MerkleRoot       [32]byte // The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp        uint32   // A timestamp recording when this block was created (Will overflow in 2106[2])
	Bits             uint32   // The calculated difficulty target being used for this block
	Nonce            uint32   // The nonce used to generate this block… to allow variations of the header and compute different hashes
	TransactionCount VarInt   // Number of transaction entries, this value is always 0
}

// Headers packet returns block headers in response to a getheaders packet.
type Headers struct {
	Count   VarInt
	Headers []BlockHeader `size:"Count"`
}

// Ping message is sent primarily to confirm that the TCP/IP connection is still valid. An error in transmission is presumed to be a closed connection and the address is removed as a current peer.
type Ping struct {
	Nonce uint64 // random nonce
}

// Pong message is sent in response to a ping message. In modern protocol versions, a pong response is generated using a nonce included in the ping.
type Pong struct {
	Nonce uint64 // nonce from ping
}

// Reject message is sent when messages are rejected.
type Reject struct {
	Message VarString // type of message rejected
	CCode   byte      // code relating to rejected message
	Reason  VarString // text version of reason for rejection
	Data    [32]byte  `deserialize:"oimt"` // Optional extra data provided by some errors. Currently, all errors which provide this field fill it with the TXID or block header hash of the object being rejected, so the field is 32 bytes.
}

// FilterLoad
type FilterLoad struct {
	count      VarInt
	Filter     []byte `size:"count"` // The filter itself is simply a bit field of arbitrary byte-aligned size. The maximum size is 36,000 bytes
	NHashFuncs uint32 // The number of hash functions to use in this filter. The maximum value allowed in this field is 50.
	NTweak     uint32 // A random value to add to the seed value in the hash function used by the bloom filter.
	NFlags     uint8  // A set of flags that control how matched items are added to the filter.
}

// FilterAdd
type FilterAdd struct {
	count VarInt
	Data  []byte `size:"count"` // The data element to add to the current filter.
}

// MerkleBlock
type MerkleBlock struct {
	Version           int32      // Block version information, based upon the software version creating this block (note, this is signed)
	PrevBlock         [32]byte   // The hash value of the previous block this particular block references
	MerkleRoot        [32]byte   // The reference to a Merkle tree collection which is a hash of all transactions related to this block
	Timestamp         uint32     // A timestamp recording when this block was created (Will overflow in 2106[2])
	Bits              uint32     // The calculated difficulty target being used for this block
	Nonce             uint32     // The nonce used to generate this block… to allow variations of the header and compute different hashes
	TotalTransactions uint32     // Number of transactions in the block (including unmatched ones)
	HashCount         VarInt     // The number of hashes to follow
	Hashes            [][32]byte `size:"HashCount"` // Hashes in depth-first order
	FlagBytes         VarInt     // The size of flags (in bytes) to follow
	Flags             []byte     `size:"FlagBytes"` // Flag bits, packed per 8 in a byte, least significant bit first. Extra 0 bits are padded on to reach full byte size.
}

type HeaderAndShortIDs struct {
	Header       BlockHeader
	Nonce        uint64
	shortIDCount VarInt
	ShortIDs     [][6]byte
}
