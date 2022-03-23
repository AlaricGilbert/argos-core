package bitcoin

import (
	"context"
	"math/rand"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/cloudwego/netpoll"
)

type Daemon struct {
	ctx        *argos.Context
	addr       *netpoll.TCPAddr
	conn       *netpoll.TCPConnection
	txHandler  argos.TransactionHandler
	mock       bool
	mockReader netpoll.Reader
	mockWriter netpoll.Writer
	nonce      uint64
}

func (d *Daemon) send(command string, data any) error {
	buf := netpoll.NewLinkBuffer()

	if _, err := serialization.Serialize(buf, data); err != nil {
		return err
	}

	// link buffer never returns an error
	_ = buf.Flush()

	msgLen := buf.Len()

	msg, _ := buf.ReadBinary(msgLen)

	var cmd [12]byte
	copy(cmd[:], []byte(command))

	header := &MessageHeader{
		Magic:    MagicMain,
		Command:  cmd,
		Length:   uint32(msgLen),
		Checksum: checksum(msg),
	}

	w := d.writer()

	if _, err := serialization.Serialize(w, header); err != nil {
		return err
	}

	if _, err := w.WriteBinary(msg); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	_ = buf.Close()

	return nil
}

func (d *Daemon) sendReject(msg string, code byte, reason string, data [32]byte) error {
	return d.send(CommandReject, &Reject{
		Message: VarString(msg),
		CCode:   code,
		Reason:  VarString(reason),
		Data:    data,
	})
}

func (d *Daemon) sendVersion() error {
	return d.send(CommandVersion, &Version{
		Version:      70015,
		Services:     NODE_NETWORK,
		Timestamp:    time.Now().Unix(),
		AddrReceived: *newNetworkAddress(NODE_NETWORK, &d.addr.TCPAddr),
		AddrFrom:     *newNetworkAddress(NODE_NETWORK, &d.addr.TCPAddr),
		Nonce:        d.nonce,
		UserAgent:    UserAgent,
		StartHeight:  0,
		Relay:        false,
	})
}

func (d *Daemon) reader() netpoll.Reader {
	if d.mock && d.mockReader != nil {
		return d.mockReader
	}
	return d.conn.Reader()
}

func (d *Daemon) writer() netpoll.Writer {
	if d.mock && d.mockWriter != nil {
		return d.mockWriter
	}
	return d.conn.Writer()
}

func (d *Daemon) header() (*MessageHeader, error) {
	defer d.reader().Release()

	if d.mock || d.conn.IsActive() {
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
			if b, err = d.reader().ReadByte(); err != nil {
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
		if _, err = serialization.Deserialize(d.reader(), &header.Command); err != nil {
			return nil, err
		}
		if _, err = serialization.Deserialize(d.reader(), &header.Length); err != nil {
			return nil, err
		}
		if _, err = serialization.Deserialize(d.reader(), &header.Checksum); err != nil {
			return nil, err
		}

		return &header, nil
	}
	return nil, argos.DaemonNotRunningError
}

func (d *Daemon) handle() error {
	defer d.reader().Release()
	defer d.writer().Flush()
	var rejectData [32]byte

	h, err := d.header()
	if err != nil {
		return err
	}

	command := SliceToString(h.Command[:])
	// received unexpected message that more than buffer limit
	// consider it's a message that with wrong transmission status
	// send a reject message
	if h.Length > BitcoinMessageMaxLength {
		return d.sendReject(command, REJECT_INVALID, "message too long", rejectData)
	}

	slice, err := d.reader().Slice(int(h.Length))
	if err != nil {
		return err
	}

	payload, _ := slice.ReadBinary(int(h.Length))
	if checksum(payload) != h.Checksum {
		return d.sendReject(command, REJECT_INVALID, "message checksum invalid", rejectData)
	}

	switch command {
	case CommandReject, CommandVerack:
		// DO nothing
	case CommandVersion:
		return d.handleVersion(d.reader())
	default:
		d.sendReject(command, REJECT_INVALID, "unsupported", rejectData)
	}

	return nil
}

func (d *Daemon) handleVersion(reader netpoll.Reader) error {
	var version Version
	serialization.Deserialize(reader, &version)
	return d.sendVersion()
}

func (d *Daemon) Spin() error {
	var err error

	defer d.conn.Close()

	if !d.mock {
		d.nonce = rand.Uint64()

		if d.conn, err = netpoll.DialTCP(context.Background(), "tcp", nil, d.addr); err != nil {
			return err
		}
	}

	if err = d.sendVersion(); err != nil {
		return err
	}

	for d.mock || d.conn.IsActive() {
		if err = d.handle(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Daemon) Halt() error {
	if d.conn == nil || !d.conn.IsActive() {
		return argos.DaemonNotRunningError
	}
	return d.conn.Close()
}

func (d *Daemon) OnTransactionReceived(handler argos.TransactionHandler) {
	d.txHandler = handler
}

func NewDaemon(ctx *argos.Context, addr *net.TCPAddr) argos.Daemon {
	return &Daemon{
		ctx: ctx,
		addr: &netpoll.TCPAddr{
			TCPAddr: *addr,
		},
	}
}

func (d *Daemon) Mock(reader netpoll.Reader, writer netpoll.Writer) {
	d.mock = true
	d.mockReader = reader
	d.mockWriter = writer
}
