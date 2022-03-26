package bitcoin

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/AlaricGilbert/argos-core/argos/sniffer"
	"github.com/cloudwego/netpoll"
	"github.com/sirupsen/logrus"
)

type Peer struct {
	s           *sniffer.Sniffer
	addr        *netpoll.TCPAddr
	localAddr   *netpoll.TCPAddr
	conn        *netpoll.TCPConnection
	txs         map[[32]byte]Transaction
	announce    bool
	sendheaders bool
	filterLoad  *FilterLoad
	feeFilter   int64
	mock        bool
	mockReader  netpoll.Reader
	mockWriter  netpoll.Writer
	nonce       uint64
}

func (d *Peer) logger() *logrus.Entry {
	return d.s.Logger().WithFields(logrus.Fields{
		"address": d.addr,
		"nonce":   fmt.Sprintf("0x%x", d.nonce),
	})
}

func (d *Peer) send(command string, data any) error {
	var err error
	defer func() {
		if err != nil {
			d.logger().WithError(err).Warn("bitcoin peer sending message failed")
		}
	}()
	buf := netpoll.NewLinkBuffer()

	if _, err = serialization.Serialize(buf, data); err != nil {
		return err
	}

	// link buffer never returns an error
	_ = buf.Flush()

	msgLen := buf.Len()

	msg, _ := buf.ReadBinary(msgLen)

	var cmd [12]byte
	copy(cmd[:], []byte(command))

	_, sum := checksum(msg)
	header := &MessageHeader{
		Magic:    MagicMain,
		Command:  cmd,
		Length:   uint32(msgLen),
		Checksum: sum,
	}

	d.logger().WithFields(logrus.Fields{
		"command": command,
		"header":  header,
		"message": data,
	}).Info("bitcoin peer sending message")

	w := d.writer()

	if _, err = serialization.Serialize(w, header); err != nil {
		return err
	}

	if _, err = w.WriteBinary(msg); err != nil {
		return err
	}

	if err = w.Flush(); err != nil {
		return err
	}

	_ = buf.Close()

	return nil
}

func (d *Peer) sendReject(msg string, code byte, reason string, data [32]byte) error {
	return d.send(CommandReject, &Reject{
		Message: VarString(msg),
		CCode:   code,
		Reason:  VarString(reason),
		Data:    data,
	})
}

func (d *Peer) sendVersion() error {
	addr := d.addr
	addr.IP = addr.IP.To16()
	return d.send(CommandVersion, &Version{
		Version:      70015,
		Services:     0,
		Timestamp:    time.Now().Unix(),
		AddrReceived: *newNetworkAddress(0, &addr.TCPAddr),
		AddrFrom:     *newNetworkAddress(0, &addr.TCPAddr),
		Nonce:        d.nonce,
		UserAgent:    UserAgent,
		StartHeight:  0,
		Relay:        true,
	})
}

func (d *Peer) sendVerack() error {
	return d.send(CommandVerack, nil)
}

func (d *Peer) sendInv(invs ...Inventory) error {
	return d.send(CommandInv, &Inv{
		Count:     VarInt(len(invs)),
		Inventory: invs,
	})
}

func (d *Peer) sendGetData(invs ...Inventory) error {
	return d.send(CommandGetData, &Inv{
		Count:     VarInt(len(invs)),
		Inventory: invs,
	})
}

func (d *Peer) sendNotFound(invs ...Inventory) error {
	return d.send(CommandNotFound, &Inv{
		Count:     VarInt(len(invs)),
		Inventory: invs,
	})
}

func (d *Peer) sendPong(nonce uint64) error {
	return d.send(CommandPong, &Pong{
		Nonce: nonce,
	})
}

func (d *Peer) reader() netpoll.Reader {
	if d.mock && d.mockReader != nil {
		return d.mockReader
	}
	return d.conn.Reader()
}

func (d *Peer) writer() netpoll.Writer {
	if d.mock && d.mockWriter != nil {
		return d.mockWriter
	}
	return d.conn.Writer()
}

func (d *Peer) header(ctx *Ctx) {
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
		// if a magic code had been totally parsed
		var parsed bool = false
		// parsed magic code
		var magic NetworkMagic

		for !parsed {
			// tries to Read a byte from connection and returns error when it fails
			if b, ctx.err = d.reader().ReadByte(); ctx.err != nil {
				return
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

		ctx.header.Magic = magic
		if _, ctx.err = serialization.Deserialize(d.reader(), &ctx.header.Command); ctx.err != nil {
			return
		}
		if _, ctx.err = serialization.Deserialize(d.reader(), &ctx.header.Length); ctx.err != nil {
			return
		}
		if _, ctx.err = serialization.Deserialize(d.reader(), &ctx.header.Checksum); ctx.err != nil {
			return
		}
		return
	}
	ctx.err = sniffer.PeerNotRunningError
}

func (d *Peer) handle() error {
	var rejectData [32]byte
	var data []byte
	var ctx = &Ctx{
		peer: d,
	}
	defer func() {
		d.reader().Release()
		d.writer().Flush()
		if ctx.payload != nil {
			ctx.payload.Close()
		}
		if ctx.err != nil {
			d.logger().WithError(ctx.err).Warn("bitcoin peer handle func exited with error")
		}
	}()

	if d.header(ctx); ctx.err != nil {
		return ctx.err
	}

	ctx.command = SliceToString(ctx.header.Command[:])
	// received unexpected message that more than buffer limit
	// consider it's a message that with wrong transmission status
	// send a reject message
	if ctx.header.Length > BitcoinMessageMaxLength {
		ctx.err = d.sendReject(ctx.command, REJECT_INVALID, "message too long", rejectData)
		return ctx.err
	}

	if data, ctx.err = d.reader().ReadBinary(int(ctx.header.Length)); ctx.err != nil {
		panic(ctx.err)
		return ctx.err
	} else if ctx.payloadhash, ctx.checksum = checksum(data); ctx.header.Checksum != ctx.checksum {
		ctx.err = d.sendReject(ctx.command, REJECT_INVALID, "message checksum invalid", rejectData)
		return ctx.err
	}

	ctx.payload = netpoll.NewLinkBuffer()
	ctx.payload.WriteBinary(data)
	ctx.payload.Flush()

	if handler, ok := commandHandlers[ctx.command]; ok {
		handler(ctx)
	} else {
		ctx.err = d.sendReject(ctx.command, REJECT_INVALID, "unsupported", rejectData)
	}
	return ctx.err
}

func (d *Peer) Spin() error {
	var err error

	defer func() {
		if d.conn != nil {
			d.conn.Close()
		}
		d.logger().Info("bitcoin peer spin exited")
	}()

	d.logger().Info("bitcoin peer spinning")
	if !d.mock {
		d.nonce = rand.Uint64()

		if d.conn, err = netpoll.DialTCP(context.Background(), "tcp", nil, d.addr); err != nil {
			d.logger().WithError(err).Error("peer connect failed")
			return err
		}

		d.localAddr = &netpoll.TCPAddr{
			TCPAddr: *d.conn.LocalAddr().(*net.TCPAddr),
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

func (d *Peer) Halt() error {
	d.logger().Info("bitcoin peer spin halting")
	if d.conn == nil || !d.conn.IsActive() {
		return sniffer.PeerNotRunningError
	}
	return d.conn.Close()
}

func NewPeer(ctx *sniffer.Sniffer, addr *net.TCPAddr) sniffer.Peer {
	return &Peer{
		s: ctx,
		addr: &netpoll.TCPAddr{
			TCPAddr: *addr,
		},
		txs: make(map[[32]byte]Transaction),
	}
}

func (d *Peer) Mock(reader netpoll.Reader, writer netpoll.Writer) {
	d.mock = true
	d.mockReader = reader
	d.mockWriter = writer

	d.logger().Info("bitcoin peer mocked")
}
