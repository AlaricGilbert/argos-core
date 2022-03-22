package bitcoin

import (
	"context"
	"net"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/AlaricGilbert/argos-core/protocol"
	"github.com/cloudwego/netpoll"
)

type Client struct {
	ctx        *argos.Context
	addr       net.Addr
	tcpAddr    *netpoll.TCPAddr
	conn       *netpoll.TCPConnection
	txHandler  argos.TransactionHandler
	mock       bool
	mockReader netpoll.Reader
}

func (c *Client) Spin() error {
	var err error

	defer c.conn.Close()

	if c.tcpAddr, err = netpoll.ResolveTCPAddr(c.addr.Network(), c.addr.String()); err != nil {
		return err
	}
	if c.conn, err = netpoll.DialTCP(context.Background(), "tcp", nil, c.tcpAddr); err != nil {
		return err
	}

	return nil
}

func (c *Client) reader() netpoll.Reader {
	if c.mockReader != nil {
		return c.mockReader
	}
	return c.conn.Reader()
}

func (c *Client) header() (*MessageHeader, error) {
	defer c.reader().Release()

	if c.mock || c.conn.IsActive() {
		// an easy handwritten state machine to detect magic code

		// there are five possible arrangements of magic code:
		//		{0xF9, 0xBE, 0xB4, 0xD9} (MagicMain     NetworkMagic = 0xD9B4BEF9)
		//		{0xF9, 0xBE, 0xB4, 0xFE} (MagicNamecoin NetworkMagic = 0xFEB4BEF9)
		//		{0xFA, 0xBF, 0xB5, 0xDA} (MagicTestnet  NetworkMagic = 0xDAB5BFFA)
		//		{0x0B, 0x11, 0x09, 0x07} (MagicTestnet3 NetworkMagic = 0x0709110B)
		//		{0x0A, 0x03, 0xCF, 0x40} (MagicSignet   NetworkMagic = 0x40CF030A)

		// every magic code consists of 4 bytes so there should be 5 internal states
		// and we use a temporary variable t to recognize prefixes (Main and Namecoin are same in pre 3 bytes)

		// state zero means init state or error state, states > 0 means code.to_bytes()[state - 1] has been detected
		var state = 0
		// t = 0 := magic starts with 0xf9
		// t = 1 := magic starts with 0xfa
		// t = 2 := magic starts with 0x0b
		// t = 3 := magic starts with 0x0a
		var t = 0

		// there are four possible magic detect starting
		var inits = []byte{0xf9, 0xfa, 0x0b, 0x0a}

		// middle two bytes state transfer array
		var trans = [][]byte{
			{0xbe, 0xbf, 0x11, 0x03},
			{0xb4, 0xb5, 0x09, 0xcf},
		}

		// last byte to magic code
		var ending = map[byte]NetworkMagic{
			0xd9: MagicMain,
			0xda: MagicTestnet,
			0x07: MagicTestnet3,
			0x40: MagicSignet,
			0xfe: MagicNamecoin,
		}

		// current byte
		var b byte
		// Read error
		var err error
		// if a magic code had been totally parsed
		var parsed bool = false
		// parsed magic code
		var magic NetworkMagic

		for !parsed {
			// tries to Read a byte from connection and returns error when it fails
			if b, err = c.reader().ReadByte(); err != nil {
				return nil, err
			}

			// despite any state now, when we met start codes, reset the state and t immediately
			if i, ok := Index(inits, b); ok {
				state = 1
				t = i
				continue
			}

			// previous 3 bytes are arranged valid, check last byte
			if state == 3 {
				if code, ok := ending[b]; ok {
					magic = code
					parsed = true
				} else {
					state = 0
				}
				continue
			}

			if state > 0 {
				if i, ok := Index(trans[state-1], b); ok && i == t {
					state++
				} else {
					state = 0
				}
			}
		}

		var header MessageHeader
		header.Magic = magic
		if _, err = serialization.Deserialize(c.reader(), &header.Command); err != nil {
			return nil, err
		}
		if _, err = serialization.Deserialize(c.reader(), &header.Length); err != nil {
			return nil, err
		}
		if _, err = serialization.Deserialize(c.reader(), &header.Checksum); err != nil {
			return nil, err
		}

		return &header, nil
	}
	return nil, protocol.ClientNotRunningError
}

func (c *Client) payload(header *MessageHeader) (any, error) {
	if header == nil || header.Length == 0 {
		return nil, nil
	}

	var payload interface{}
	var buf netpoll.Reader
	var err error
	var n int
	// prefetch buffer
	buf, err = c.reader().Slice(int(header.Length))
	if err != nil {
		return nil, err
	}

	switch string(SliceToString(header.Command[:])) {
	case "version":
		payload = &Version{}
	case "addr":
		payload = &Addr{}
	case "inv", "getdata", "notfound":
		payload = &Inventory{}
	case "tx":
		payload = &Transaction{}
	case "block":
		payload = &Block{}
	case "headers":
		payload = &Headers{}
	case "ping":
		payload = &Ping{}
	case "pong":
		payload = &Pong{}
	case "reject":
		payload = &Reject{}
	case "filteradd":
		payload = &FilterAdd{}
	case "filterload":
		payload = &FilterLoad{}
	case "merkleblock":
		payload = &MerkleBlock{}
	case "getaddr", "mempool", "verack", "filterclear":
		return nil, nil
	}
	n, err = serialization.Deserialize(buf, payload)
	if err != nil || n != int(header.Length) {
		return nil, err
	}

	return payload, nil
}

func (c *Client) Halt() error {
	if c.conn == nil || !c.conn.IsActive() {
		return protocol.ClientNotRunningError
	}
	return c.conn.Close()
}

func (c *Client) OnTransactionReceived(handler argos.TransactionHandler) {
	c.txHandler = handler
}

func NewClient(ctx *argos.Context, addr net.Addr) argos.Client {
	return &Client{
		ctx:  ctx,
		addr: addr,
	}
}

func (c *Client) Mock(reader netpoll.Reader) {
	c.mock = true
	c.mockReader = reader
}
