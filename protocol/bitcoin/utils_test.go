package bitcoin

import (
	"testing"

	"github.com/cloudwego/netpoll"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	var arr = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, b := range arr {
		if id, ok := Index(arr, b); !ok || id != i {
			t.Errorf("Index() test failed on i=%d, got Index=%d and status=%t ", i, id, ok)
		}
	}
}

func TestDeserializeVersion(t *testing.T) {
	initOnce()
	var data = []byte{
		0xf9, 0xbe, 0xb4, 0xd9, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x65, 0x00, 0x00, 0x00, 0x35, 0x8d, 0x49, 0x32, 0x62, 0xea, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x11, 0xb2, 0xd0, 0x50, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x3b, 0x2e, 0xb3, 0x5d, 0x8c, 0xe6, 0x17, 0x65, 0x0f, 0x2f, 0x53, 0x61, 0x74, 0x6f, 0x73, 0x68,
		0x69, 0x3a, 0x30, 0x2e, 0x37, 0x2e, 0x32, 0x2f, 0xc0, 0x3e, 0x03, 0x00, 0x01,
	}

	// 0000   f9 be b4 d9 76 65 72 73 69 6f 6e 00 00 00 00 00  ....version.....
	// 0010   64 00 00 00 35 8d 49 32 62 ea 00 00 01 00 00 00  d...5.I2b.......
	// 0020   00 00 00 00 11 b2 d0 50 00 00 00 00 01 00 00 00  .......P........
	// 0030   00 00 00 00 00 00 00 00 00 00 00 00 00 00 ff ff  ................
	// 0040   00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	// 0050   00 00 00 00 00 00 00 00 ff ff 00 00 00 00 00 00  ................
	// 0060   3b 2e b3 5d 8c e6 17 65 0f 2f 53 61 74 6f 73 68  ;..]...e./Satosh
	// 0070   69 3a 30 2e 37 2e 32 2f c0 3e 03 00 01           i:0.7.2/.>...

	// Message Header:
	//  F9 BE B4 D9                                                                   - Main network magic bytes
	//  76 65 72 73 69 6F 6E 00 00 00 00 00                                           - "version" command
	//  64 00 00 00                                                                   - Payload is 100 bytes long
	//  35 8d 49 32                                                                   - payload checksum (internal byte order)

	// Version message:
	//  62 EA 00 00                                                                   - 60002 (protocol version 60002)
	//  01 00 00 00 00 00 00 00                                                       - 1 (NODE_NETWORK services)
	//  11 B2 D0 50 00 00 00 00                                                       - Tue Dec 18 10:12:33 PST 2012
	//  01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF FF 00 00 00 00 00 00 - Recipient address info - see Network Address
	//  01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF FF 00 00 00 00 00 00 - Sender address info - see Network Address
	//  3B 2E B3 5D 8C E6 17 65                                                       - Node ID
	//  0F 2F 53 61 74 6F 73 68 69 3A 30 2E 37 2E 32 2F                               - "/Satoshi:0.7.2/" sub-version string (string is 15 bytes long)
	//  C0 3E 03 00                                                                   - Last block sending node has is block #212672
	//  01																			  - Relay enabled
	var ok bool
	var ver *Version
	var header *MessageHeader
	//var addr Addr
	var payload interface{}
	var err error
	//var n int

	buf := netpoll.NewLinkBuffer()
	_, _ = buf.WriteBinary(data)
	_ = buf.Flush()

	c := NewClient(nil, nil).(*Client)
	c.Mock(buf)

	header, err = c.header()
	assert.Nil(t, err, "while reading header")
	assert.NotNil(t, header)

	assert.Equal(t, uint32(101), header.length)
	assert.Equal(t, "version", SliceToString(header.command[:]))

	payload, err = c.payload(header)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	ver, ok = payload.(*Version)
	assert.True(t, ok)
	assert.NotNil(t, ver)

	assert.Equal(t, int32(60002), ver.Version)
	assert.True(t, ver.Services.Serves(NODE_NETWORK))
	assert.Equal(t, int64(1355854353), ver.Timestamp)

	addrFrom := ver.AddrFrom.TCPAddr()
	addrReceived := ver.AddrReceived.TCPAddr()

	assert.Equal(t, "0.0.0.0:0", addrFrom.String())
	assert.Equal(t, "0.0.0.0:0", addrReceived.String())
	assert.Equal(t, uint64(0x6517e68c5db32e3b), ver.Nonce)
	assert.Equal(t, VarString("/Satoshi:0.7.2/"), ver.UserAgent)
	assert.Equal(t, int32(212672), ver.StartHeight)
	assert.True(t, ver.Relay)
}

func TestDeserializeAddr(t *testing.T) {
	initOnce()

	var data = []byte{
		0xF9, 0xBE, 0xB4, 0xD9, 0x61, 0x64, 0x64, 0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x1F, 0x00, 0x00, 0x00, 0xED, 0x52, 0x39, 0x9B, 0x01, 0xE2, 0x15, 0x10, 0x4D, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF,
		0xFF, 0x0A, 0x00, 0x00, 0x01, 0x20, 0x8D,
	}

	// 0000   F9 BE B4 D9 61 64 64 72  00 00 00 00 00 00 00 00   ....addr........
	// 0010   1F 00 00 00 ED 52 39 9B  01 E2 15 10 4D 01 00 00   .....R9.....M...
	// 0020   00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 FF   ................
	// 0030   FF 0A 00 00 01 20 8D                               ..... .

	// Message Header:
	//  F9 BE B4 D9                                     - Main network magic bytes
	//  61 64 64 72  00 00 00 00 00 00 00 00            - "addr"
	//  1F 00 00 00                                     - payload is 31 bytes long
	//  ED 52 39 9B                                     - payload checksum (internal byte order)

	// Payload:
	//  01                                              - 1 address in this message

	// Address:
	//  E2 15 10 4D                                     - Mon Dec 20 21:50:10 EST 2010 (only when version is >= 31402)
	//  01 00 00 00 00 00 00 00                         - 1 (NODE_NETWORK service - see version message)
	//  00 00 00 00 00 00 00 00 00 00 FF FF 0A 00 00 01 - IPv4: 10.0.0.1, IPv6: ::ffff:10.0.0.1 (IPv4-mapped IPv6 address)
	//  20 8D                                           - port 8333
	var header *MessageHeader
	var addr *Addr
	var ok bool
	var payload interface{}
	var err error

	buf := netpoll.NewLinkBuffer()
	_, _ = buf.WriteBinary(data)
	_ = buf.Flush()

	c := NewClient(nil, nil).(*Client)
	c.Mock(buf)

	header, err = c.header()
	assert.Nil(t, err, "while reading header")
	assert.NotNil(t, header)

	assert.Equal(t, uint32(31), header.length)
	assert.Equal(t, "addr", SliceToString(header.command[:]))

	payload, err = c.payload(header)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	addr, ok = payload.(*Addr)
	assert.True(t, ok)
	assert.NotNil(t, addr)

	assert.Equal(t, VarInt(1), addr.Count)
	assert.Equal(t, uint32(0x4d1015e2), addr.AddrList[0].Time)
	assert.Equal(t, "10.0.0.1:8333", addr.AddrList[0].TCPAddr().String())
}

func TestDeserializeTx(t *testing.T) {
	initOnce()

	var data = []byte{
		0xF9, 0xBE, 0xB4, 0xD9, 0x74, 0x78, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x01, 0x00, 0x00, 0xE2, 0x93, 0xCD, 0xBE, 0x01, 0x00, 0x00, 0x00, 0x01, 0x6D, 0xBD, 0xDB,
		0x08, 0x5B, 0x1D, 0x8A, 0xF7, 0x51, 0x84, 0xF0, 0xBC, 0x01, 0xFA, 0xD5, 0x8D, 0x12, 0x66, 0xE9,
		0xB6, 0x3B, 0x50, 0x88, 0x19, 0x90, 0xE4, 0xB4, 0x0D, 0x6A, 0xEE, 0x36, 0x29, 0x00, 0x00, 0x00,
		0x00, 0x8B, 0x48, 0x30, 0x45, 0x02, 0x21, 0x00, 0xF3, 0x58, 0x1E, 0x19, 0x72, 0xAE, 0x8A, 0xC7,
		0xC7, 0x36, 0x7A, 0x7A, 0x25, 0x3B, 0xC1, 0x13, 0x52, 0x23, 0xAD, 0xB9, 0xA4, 0x68, 0xBB, 0x3A,
		0x59, 0x23, 0x3F, 0x45, 0xBC, 0x57, 0x83, 0x80, 0x02, 0x20, 0x59, 0xAF, 0x01, 0xCA, 0x17, 0xD0,
		0x0E, 0x41, 0x83, 0x7A, 0x1D, 0x58, 0xE9, 0x7A, 0xA3, 0x1B, 0xAE, 0x58, 0x4E, 0xDE, 0xC2, 0x8D,
		0x35, 0xBD, 0x96, 0x92, 0x36, 0x90, 0x91, 0x3B, 0xAE, 0x9A, 0x01, 0x41, 0x04, 0x9C, 0x02, 0xBF,
		0xC9, 0x7E, 0xF2, 0x36, 0xCE, 0x6D, 0x8F, 0xE5, 0xD9, 0x40, 0x13, 0xC7, 0x21, 0xE9, 0x15, 0x98,
		0x2A, 0xCD, 0x2B, 0x12, 0xB6, 0x5D, 0x9B, 0x7D, 0x59, 0xE2, 0x0A, 0x84, 0x20, 0x05, 0xF8, 0xFC,
		0x4E, 0x02, 0x53, 0x2E, 0x87, 0x3D, 0x37, 0xB9, 0x6F, 0x09, 0xD6, 0xD4, 0x51, 0x1A, 0xDA, 0x8F,
		0x14, 0x04, 0x2F, 0x46, 0x61, 0x4A, 0x4C, 0x70, 0xC0, 0xF1, 0x4B, 0xEF, 0xF5, 0xFF, 0xFF, 0xFF,
		0xFF, 0x02, 0x40, 0x4B, 0x4C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0x76, 0xA9, 0x14, 0x1A, 0xA0,
		0xCD, 0x1C, 0xBE, 0xA6, 0xE7, 0x45, 0x8A, 0x7A, 0xBA, 0xD5, 0x12, 0xA9, 0xD9, 0xEA, 0x1A, 0xFB,
		0x22, 0x5E, 0x88, 0xAC, 0x80, 0xFA, 0xE9, 0xC7, 0x00, 0x00, 0x00, 0x00, 0x19, 0x76, 0xA9, 0x14,
		0x0E, 0xAB, 0x5B, 0xEA, 0x43, 0x6A, 0x04, 0x84, 0xCF, 0xAB, 0x12, 0x48, 0x5E, 0xFD, 0xA0, 0xB7,
		0x8B, 0x4E, 0xCC, 0x52, 0x88, 0xAC, 0x00, 0x00, 0x00, 0x00,
	}

	// 000000	F9 BE B4 D9 74 78 00 00  00 00 00 00 00 00 00 00   ....tx..........
	// 000010	02 01 00 00 E2 93 CD BE  01 00 00 00 01 6D BD DB   .............m..
	// 000020	08 5B 1D 8A F7 51 84 F0  BC 01 FA D5 8D 12 66 E9   .[...Q........f.
	// 000030	B6 3B 50 88 19 90 E4 B4  0D 6A EE 36 29 00 00 00   .;P......j.6)...
	// 000040	00 8B 48 30 45 02 21 00  F3 58 1E 19 72 AE 8A C7   ..H0E.!..X..r...
	// 000050	C7 36 7A 7A 25 3B C1 13  52 23 AD B9 A4 68 BB 3A   .6zz%;..R#...h.:
	// 000060	59 23 3F 45 BC 57 83 80  02 20 59 AF 01 CA 17 D0   Y#?E.W... Y.....
	// 000070	0E 41 83 7A 1D 58 E9 7A  A3 1B AE 58 4E DE C2 8D   .A.z.X.z...XN...
	// 000080	35 BD 96 92 36 90 91 3B  AE 9A 01 41 04 9C 02 BF   5...6..;...A....
	// 000090	C9 7E F2 36 CE 6D 8F E5  D9 40 13 C7 21 E9 15 98   .~.6.m...@..!...
	// 0000A0	2A CD 2B 12 B6 5D 9B 7D  59 E2 0A 84 20 05 F8 FC   *.+..].}Y... ...
	// 0000B0	4E 02 53 2E 87 3D 37 B9  6F 09 D6 D4 51 1A DA 8F   N.S..=7.o...Q...
	// 0000C0	14 04 2F 46 61 4A 4C 70  C0 F1 4B EF F5 FF FF FF   ../FaJLp..K.....
	// 0000D0	FF 02 40 4B 4C 00 00 00  00 00 19 76 A9 14 1A A0   ..@KL......v....
	// 0000E0	CD 1C BE A6 E7 45 8A 7A  BA D5 12 A9 D9 EA 1A FB   .....E.z........
	// 0000F0	22 5E 88 AC 80 FA E9 C7  00 00 00 00 19 76 A9 14   "^...........v..
	// 000100	0E AB 5B EA 43 6A 04 84  CF AB 12 48 5E FD A0 B7   ..[.Cj.....H^...
	// 000110	8B 4E CC 52 88 AC 00 00  00 00                     .N.R......

	// Message header:
	//  F9 BE B4 D9                                       - main network magic bytes
	//  74 78 00 00 00 00 00 00 00 00 00 00               - "tx" command
	//  02 01 00 00                                       - payload is 258 bytes long
	//  E2 93 CD BE                                       - payload checksum (internal byte order)

	// Transaction:
	//  01 00 00 00                                       - version

	// Inputs:
	//  01                                                - number of transaction inputs

	// Input 1:
	//  6D BD DB 08 5B 1D 8A F7  51 84 F0 BC 01 FA D5 8D  - previous output (outpoint)
	//  12 66 E9 B6 3B 50 88 19  90 E4 B4 0D 6A EE 36 29
	//  00 00 00 00

	//  8B                                                - script is 139 bytes long

	//  48 30 45 02 21 00 F3 58  1E 19 72 AE 8A C7 C7 36  - signature script (scriptSig)
	//  7A 7A 25 3B C1 13 52 23  AD B9 A4 68 BB 3A 59 23
	//  3F 45 BC 57 83 80 02 20  59 AF 01 CA 17 D0 0E 41
	//  83 7A 1D 58 E9 7A A3 1B  AE 58 4E DE C2 8D 35 BD
	//  96 92 36 90 91 3B AE 9A  01 41 04 9C 02 BF C9 7E
	//  F2 36 CE 6D 8F E5 D9 40  13 C7 21 E9 15 98 2A CD
	//  2B 12 B6 5D 9B 7D 59 E2  0A 84 20 05 F8 FC 4E 02
	//  53 2E 87 3D 37 B9 6F 09  D6 D4 51 1A DA 8F 14 04
	//  2F 46 61 4A 4C 70 C0 F1  4B EF F5

	//  FF FF FF FF                                       - sequence

	// Outputs:
	//  02                                                - 2 Output Transactions

	// Output 1:
	//  40 4B 4C 00 00 00 00 00                           - 0.05 BTC (5000000)
	//  19                                                - pk_script is 25 bytes long

	//  76 A9 14 1A A0 CD 1C BE  A6 E7 45 8A 7A BA D5 12  - pk_script
	//  A9 D9 EA 1A FB 22 5E 88  AC

	// Output 2:
	//  80 FA E9 C7 00 00 00 00                           - 33.54 BTC (3354000000)
	//  19                                                - pk_script is 25 bytes long

	//  76 A9 14 0E AB 5B EA 43  6A 04 84 CF AB 12 48 5E  - pk_script
	//  FD A0 B7 8B 4E CC 52 88  AC

	// Locktime:
	//  00 00 00 00                                       - lock time
	var header *MessageHeader
	var tx *Transaction
	var ok bool
	var payload interface{}
	var err error

	buf := netpoll.NewLinkBuffer()
	_, _ = buf.WriteBinary(data)
	_ = buf.Flush()

	c := NewClient(nil, nil).(*Client)
	c.Mock(buf)

	header, err = c.header()
	assert.Nil(t, err, "while reading header")
	assert.NotNil(t, header)

	assert.Equal(t, uint32(258), header.length)
	assert.Equal(t, "tx", SliceToString(header.command[:]))

	payload, err = c.payload(header)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	tx, ok = payload.(*Transaction)
	assert.True(t, ok)
	assert.NotNil(t, tx)

	assert.Equal(t, uint32(1), tx.Version)
	assert.Equal(t, [2]uint8{0, 0}, tx.Flag)

	assert.Equal(t, VarInt(1), tx.TxInCount)
	assert.Equal(t, uint32(0xffffffff), tx.TxIn[0].Sequence)

	assert.Equal(t, VarInt(2), tx.TxOutCount)
	assert.Equal(t, int64(5000000), tx.TxOut[0].Value)
	assert.Equal(t, int64(3354000000), tx.TxOut[1].Value)

	assert.Equal(t, uint32(0), tx.LockTime)
}
