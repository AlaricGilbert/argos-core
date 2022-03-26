package bitcoin

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/netip"
	"time"
)

// MessageHeader is the header of all messages, contains a magic number which used for identify the network and locate the message start in network stream
type MessageHeader struct {
	Magic    NetworkMagic // Magic value indicating message origin network, and used to seek to next message when stream state is unknown
	Command  [12]byte     // ASCII string identifying the packet content, NULL padded (non-NULL padding results in packet rejected)
	Length   uint32       // Length of payload in number of bytes
	Checksum uint32       // First 4 bytes of sha256(sha256(payload))
}

// String implements fmt.Stringer
func (h MessageHeader) String() string {
	return fmt.Sprintf("{Magic: %v, Command: \"%s\", Length: %d, Checksum: 0x%x}", h.Magic, SliceToString(h.Command[:]), h.Length, h.Checksum)
}

// NetworkAddress is used when a network address is needed somewhere
type NetworkAddress struct {
	Time     uint32      // the Time (version >= 31402). Not present in version message.
	Services ServiceType // same service(s) listed in version
	IP       [16]byte    // IPv6 address. Network byte order. The original client only supported IPv4 and only Read the last 4 bytes to get the IPv4 address. However, the IPv4 address is written into the message as a 16 byte IPv4-mapped IPv6 address
	Port     uint16      `order:"network"` // port number, network byte order
}

// String implements fmt.Stringer
func (addr NetworkAddress) String() string {
	return fmt.Sprintf("{Timestamp: %d, Services: %v, IP: %v, Port: %d}", addr.Time, addr.Services, net.IP(addr.IP[:]), addr.Port)
}

// NewNetworkAddress creates a NetworkAddress from given service and address
func NewNetworkAddress(services ServiceType, addr *net.TCPAddr) *NetworkAddress {
	netAddr := &NetworkAddress{
		Time:     uint32(time.Now().Unix()),
		Services: services,
		Port:     uint16(addr.Port),
	}

	copy(netAddr.IP[:], addr.IP)
	return netAddr
}

// AddrPort converts NetworkAddress into netip.AddrPort
func (a *NetworkAddress) AddrPort() netip.AddrPort {
	return netip.AddrPortFrom(netip.AddrFrom16(a.IP), a.Port)
}

// AddrPort converts NetworkAddress into *net.TCPAddr
func (a *NetworkAddress) TCPAddr() *net.TCPAddr {
	return net.TCPAddrFromAddrPort(a.AddrPort())
}

// networkAddress is a substitution for NetworkAddress in Version message since NetworkAddress.Time was not present in version message
type networkAddress struct {
	Services ServiceType // same service(s) listed in version
	IP       [16]byte    // IPv6 address. Network byte order. The original client only supported IPv4 and only Read the last 4 bytes to get the IPv4 address. However, the IPv4 address is written into the message as a 16 byte IPv4-mapped IPv6 address
	Port     uint16      `order:"network"` // port number, network byte order
}

// String implements fmt.Stringer
func (addr networkAddress) String() string {
	return fmt.Sprintf("{Services: %v, IP: %v, Port: %d}", addr.Services, net.IP(addr.IP[:]), addr.Port)
}

// newNetworkAddress creates a networkAddress from given service and address
func newNetworkAddress(services ServiceType, addr *net.TCPAddr) *networkAddress {
	netAddr := &networkAddress{
		Services: services,
		Port:     uint16(addr.Port),
	}

	copy(netAddr.IP[:], addr.IP)
	return netAddr
}

// AddrPort converts networkAddress into netip.AddrPort
func (a *networkAddress) AddrPort() netip.AddrPort {
	return netip.AddrPortFrom(netip.AddrFrom16(a.IP), a.Port)
}

// AddrPort converts networkAddress into *net.TCPAddr
func (a *networkAddress) TCPAddr() *net.TCPAddr {
	return net.TCPAddrFromAddrPort(a.AddrPort())
}

// Inv allows a node to advertise its knowledge of one or more objects. It can be received unsolicited, or in reply to getblocks.
type Inv struct {
	Count     VarInt
	Inventory []Inventory `size:"Count"`
}

// Getdata is used in response to inv, to retrieve the content of a specific object, and is usually sent after receiving an inv packet, after filtering known elements. It can be used to retrieve transactions, but only if they are in the memory pool or relay set - arbitrary access to transactions in the chain is not allowed to avoid having clients start to depend on nodes having full transaction indexes (which modern nodes do not).
type GetData struct {
	Count     VarInt
	Inventory []Inventory `size:"Count"`
}

// NotFound is a response to a getdata, sent if any requested data items could not be relayed, for example, because the requested transaction was not in the memory pool or relay set.
type NotFound struct {
	Count     VarInt
	Inventory []Inventory `size:"Count"`
}

// Inventory vectors are used for notifying other nodes about objects they have or data which is being requested.
type Inventory struct {
	Type InventoryType // Identifies the object type linked to this inventory
	Hash [32]byte      // Hash of the object
}

// String implements fmt.Stringer
func (i Inventory) String() string {
	return fmt.Sprintf("{Type: %v, Hash: %s}", i.Type, hex.EncodeToString(i.Hash[:]))
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

// String implements fmt.Stringer
func (v Version) String() string {
	return fmt.Sprintf("{Version: %d, Services: %v, Timestamp: %v, AddrReceived: %v, AddrFrom: %v, Nonce: 0x%x, UserAgent: \"%s\", StartHeight: %d, Relay: %t}",
		v.Version,
		v.Services,
		v.Timestamp,
		v.AddrReceived,
		v.AddrFrom,
		v.Nonce,
		v.UserAgent,
		v.StartHeight,
		v.Relay,
	)

}

// Addr provides information on known nodes of the network. Non-advertised nodes should be forgotten after typically 3 hours
type Addr struct {
	Count    VarInt           // Number of address entries (max: 1000)
	AddrList []NetworkAddress `size:"Count"` // Address of other nodes on the network. version < 209 will only read the first one.
}

// String implements fmt.Stringer
func (addr Addr) String() string {
	return FmtSlice(addr.AddrList, func(t NetworkAddress) string {
		return t.String()
	})
}

type OutPoint struct {
	Hash  [32]byte // The hash of the referenced transaction.
	Index uint32   // The index of the specific output in the transaction. The first output is 0, etc.
}

// String implements fmt.Stringer
func (o OutPoint) String() string {
	return fmt.Sprintf("{Hash: %s, Index: %d}", hex.EncodeToString(o.Hash[:]), o.Index)
}

// TransactionIn means Transaction inputs or sources for coins
type TransactionIn struct {
	PreviousOutput  OutPoint // The previous output transaction reference, as an OutPoint structure
	ScriptLength    VarInt   // The length of the signature script
	SignatureSctipt []byte   `size:"ScriptLength"` // Computational Script for confirming transaction authorization
	Sequence        uint32   // Transaction version as defined by the sender. Intended for "replacement" of transactions when information is updated before inclusion into a block.
}

// String implements fmt.Stringer
func (t TransactionIn) String() string {
	return fmt.Sprintf("{PreviousOutput: %v, ScriptLength: %d, SignatureSctipt: %s, Sequence: %d}",
		t.PreviousOutput,
		uint64(t.ScriptLength),
		hex.EncodeToString([]byte(t.SignatureSctipt)),
		t.Sequence,
	)
}

// TransactionOut means Transaction outputs or destinations for coins
type TransactionOut struct {
	Value          int64  // Transaction Value
	PKScriptLength VarInt // Length of the pk_script
	PKScript       []byte `size:"PKScriptLength"` // Usually contains the public key as a Bitcoin script setting up conditions to claim this output.
}

// String implements fmt.Stringer
func (t TransactionOut) String() string {
	return fmt.Sprintf("{Value: %d, PKScriptLength: %d, PKScript: %s}",
		t.Value,
		uint64(t.PKScriptLength),
		hex.EncodeToString([]byte(t.PKScript)),
	)
}

type TransactionWitnesss struct {
	WitnessCount VarInt
	WitnessData  []byte `size:"WitnessCount"`
}

// String implements fmt.Stringer
func (t TransactionWitnesss) String() string {
	return FmtSlice(t.WitnessData, func(t byte) string {
		return fmt.Sprint(t)
	})
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

// String implements fmt.Stringer
func (tx Transaction) String() string {
	return fmt.Sprintf("{Version: %d, Flag: [%d, %d], TxIn: %v, TxOut: %v, TxWitness: %s, LockTime: %d}",
		tx.Version,
		tx.Flag[0], tx.Flag[1],
		FmtSlice(tx.TxIn, func(t TransactionIn) string {
			return t.String()
		}),
		FmtSlice(tx.TxOut, func(t TransactionOut) string {
			return t.String()
		}),
		FmtSlice(tx.TxWitness, func(t TransactionWitnesss) string {
			return t.String()
		}),
		tx.LockTime,
	)
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

// String implements fmt.Stringer
func (b Block) String() string {
	return fmt.Sprintf("{Version: %d, PrevBlock: %s, MerkleRoot: %s, Timestamp: %d, Bits: %d, Nonce: 0x%x, Txs: %s}",
		b.Version,
		hex.EncodeToString(b.PrevBlock[:]),
		hex.EncodeToString(b.MerkleRoot[:]),
		b.Timestamp,
		b.Bits,
		b.Nonce,
		FmtSlice(b.Txs, func(t Transaction) string {
			return t.String()
		}),
	)
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

// String implements fmt.Stringer
func (b BlockHeader) String() string {
	return fmt.Sprintf("{Version: %d, PrevBlock: %s, MerkleRoot: %s, Timestamp: %d, Bits: %d, Nonce: 0x%x}",
		b.Version,
		hex.EncodeToString(b.PrevBlock[:]),
		hex.EncodeToString(b.MerkleRoot[:]),
		b.Timestamp,
		b.Bits,
		b.Nonce,
	)
}

// Headers packet returns block headers in response to a getheaders packet.
type Headers struct {
	Count   VarInt
	Headers []BlockHeader `size:"Count"`
}

// String implements fmt.Stringer
func (h Headers) String() string {
	return FmtSlice(h.Headers, func(b BlockHeader) string {
		return b.String()
	})
}

// Ping message is sent primarily to confirm that the TCP/IP connection is still valid. An error in transmission is presumed to be a closed connection and the address is removed as a current peer.
type Ping struct {
	Nonce uint64 // random nonce
}

// String implements fmt.Stringer
func (p Ping) String() string {
	return fmt.Sprintf("{Nonce: 0x%x}", p.Nonce)
}

// Pong message is sent in response to a ping message. In modern protocol versions, a pong response is generated using a nonce included in the ping.
type Pong struct {
	Nonce uint64 // nonce from ping
}

// String implements fmt.Stringer
func (p Pong) String() string {
	return fmt.Sprintf("{Nonce: 0x%x}", p.Nonce)
}

// Reject message is sent when messages are rejected.
type Reject struct {
	Message VarString // type of message rejected
	CCode   byte      // code relating to rejected message
	Reason  VarString // text version of reason for rejection
	Data    [32]byte  `deserialize:"oimt"` // Optional extra data provided by some errors. Currently, all errors which provide this field fill it with the TXID or block header hash of the object being rejected, so the field is 32 bytes.
}

// String implements fmt.Stringer
func (r Reject) String() string {
	return fmt.Sprintf("{Message: %s, CCode: 0x%x, Reason: %s, Data: %s}",
		r.Message,
		r.CCode,
		r.Reason,
		hex.EncodeToString(r.Data[:]),
	)
}

// FilterLoad
type FilterLoad struct {
	count      VarInt
	Filter     []byte `size:"count"` // The filter itself is simply a bit field of arbitrary byte-aligned size. The maximum size is 36,000 bytes
	NHashFuncs uint32 // The number of hash functions to use in this filter. The maximum value allowed in this field is 50.
	NTweak     uint32 // A random value to add to the seed value in the hash function used by the bloom filter.
	NFlags     uint8  // A set of flags that control how matched items are added to the filter.
}

// String implements fmt.Stringer
func (f FilterLoad) String() string {
	return fmt.Sprintf("{Filter: %s, NHashFuncs: %d, NTweak: 0x%x, NFlags: %d}",
		FmtSlice(f.Filter, func(b byte) string { return fmt.Sprint(b) }),
		f.NHashFuncs,
		f.NTweak,
		f.NFlags,
	)
}

// FilterAdd
type FilterAdd struct {
	count VarInt
	Data  []byte `size:"count"` // The data element to add to the current filter.
}

// String implements fmt.Stringer
func (f FilterAdd) String() string {
	return FmtSlice(f.Data, func(b byte) string {
		return fmt.Sprint(b)
	})
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

// String implements fmt.Stringer
func (b MerkleBlock) String() string {
	return fmt.Sprintf("{Version: %d, PrevBlock: %s, MerkleRoot: %s, Timestamp: %d, Bits: %d, Nonce: 0x%x, TotalTransactions: %d, Hashes: %s, Flags: %s}",
		b.Version,
		hex.EncodeToString(b.PrevBlock[:]),
		hex.EncodeToString(b.MerkleRoot[:]),
		b.Timestamp,
		b.Bits,
		b.Nonce,
		b.TotalTransactions,
		FmtSlice(b.Hashes, func(t [32]byte) string {
			return hex.EncodeToString(t[:])
		}),
		FmtSlice(b.Flags, func(t byte) string {
			return fmt.Sprint(t)
		}),
	)
}

// FeeFilter serves to instruct peers not to send "inv"'s to the node for transactions with fees below the specified fee rate.
type FeeFilter int64

type GetHeaders struct {
	Version            uint32     // the protocol version
	HashCount          VarInt     // number of block locator hash entries
	BlockLocatorHashes [][32]byte `size:"HashCount"` // block locator object; newest back to genesis block (dense to start, but then sparse)
	HashStop           [32]byte   // hash of the last desired block; set to zero to get as many blocks as possible (500)
}

// String implements fmt.Stringer
func (gh GetHeaders) String() string {
	return fmt.Sprintf("{Version: %d, BlockLocatorHashes: %s, HashStop: %s}",
		gh.Version,
		FmtSlice(gh.BlockLocatorHashes, func(t [32]byte) string {
			return hex.EncodeToString(t[:])
		}),
		hex.EncodeToString(gh.HashStop[:]),
	)
}

type SendCmpct struct {
	// If announce is set to false the receive node must announce new blocks
	// via the standard inv relay. If announce is true, a new Compact Block
	// can be pushed directly to the peer.
	Announce bool

	// The version of this protocol is currently 1.
	Version uint64
}

type HeaderAndShortIDs struct {
	Header       BlockHeader
	Nonce        uint64
	shortIDCount VarInt
	ShortIDs     [][6]byte
}
