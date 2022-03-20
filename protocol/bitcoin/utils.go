package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/cloudwego/netpoll"
)

// hash returns the hash result of sha256(sha256(data))
func hash(data []byte) []byte {
	h := sha256.New()
	hh := sha256.New()
	return hh.Sum(h.Sum(data))
}

type indexType interface{ int | byte }

// Index finds the index of given element e in array s returns index and succeed result of the process
// TODO: change decl [T indexType] -> [T comparable], which is a temporary solution for GoLand false positive.
func Index[T indexType](s []T, e T) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return -1, false
}

// deserialize provides extension of bitcoin protocol type for serialization.Deserialize
func deserialize(r netpoll.Reader, data any, order binary.ByteOrder) (int, error) {
	// current byte
	var b byte
	// bytes read
	var bs []byte
	// error
	var err error
	// count of total bytes have been read
	var bytes = 0
	// tmp count of bytes have been read
	var n = 0

	// handle basic types defined in types.go for bitcoin protocol
	switch data := data.(type) {
	case *VarInt:
		if b, err = r.ReadByte(); err != nil {
			return 0, err
		}
		if b < 0xFD {
			*data = VarInt(b)
			return 1, nil
		}
		switch b {
		case 0xFD:
			if bs, err = r.ReadBinary(2); err != nil {
				return 1, err
			}
			*data = VarInt(binary.LittleEndian.Uint16(bs))
			return 3, nil
		case 0xFE:
			if bs, err = r.ReadBinary(4); err != nil {
				return 1, err
			}
			*data = VarInt(binary.LittleEndian.Uint32(bs))
			return 5, nil
		case 0xFF:
			if bs, err = r.ReadBinary(8); err != nil {
				return 1, err
			}
			*data = VarInt(binary.LittleEndian.Uint64(bs))
			return 9, nil
		}
	case *InventoryType:
		var invType uint32
		bytes, err = serialization.Deserialize(r, &invType)
		*data = InventoryType(invType)
		return bytes, err
	case *NetworkMagic:
		var magic uint32
		bytes, err = serialization.Deserialize(r, &magic)
		*data = NetworkMagic(magic)
		return bytes, err
	case *ServiceType:
		var svc uint64
		bytes, err = serialization.Deserialize(r, &svc)
		*data = ServiceType(svc)
		return bytes, err
	case *VarString:
		var size VarInt
		var str string
		if bytes, err = serialization.Deserialize(r, &size); err != nil {
			return bytes, err
		}
		str, err = r.ReadString(int(size))
		*data = VarString(str)
		return bytes + len(*data), err
	case *Transaction:
		if bytes, err = serialization.Deserialize(r, &data.Version); err != nil {
			return bytes, err
		}
		if n, err = serialization.Deserialize(r, &data.TxInCount); err != nil {
			return bytes + n, err
		}
		bytes += n

		data.Flag[0] = 0
		data.Flag[1] = 0
		if data.TxInCount == 0 {
			if n, err = serialization.Deserialize(r, &data.Flag[1]); err != nil {
				return bytes + n, err
			}
			bytes += n
			if n, err = serialization.Deserialize(r, &data.TxInCount); err != nil {
				return bytes + n, err
			}
			bytes += n
		}

		data.TxIn = make([]TransactionIn, data.TxInCount)
		for i := 0; i < int(data.TxInCount); i++ {
			if n, err = serialization.Deserialize(r, &data.TxIn[i]); err != nil {
				return bytes + n, err
			}
			bytes += n
		}

		if n, err = serialization.Deserialize(r, &data.TxOutCount); err != nil {
			return bytes + n, err
		}
		bytes += n

		data.TxOut = make([]TransactionOut, data.TxOutCount)
		for i := 0; i < int(data.TxOutCount); i++ {
			if n, err = serialization.Deserialize(r, &data.TxOut[i]); err != nil {
				return bytes + n, err
			}
			bytes += n
		}

		// contains witness data
		if data.Flag[1] == 1 {
			if n, err = serialization.Deserialize(r, &data.TxWitnessCount); err != nil {
				return bytes + n, err
			}
			bytes += n

			for i := 0; i < int(data.TxWitnessCount); i++ {
				if n, err = serialization.Deserialize(r, &data.TxWitness[i]); err != nil {
					return bytes + n, err
				}
				bytes += n
			}
		}

		n, err = serialization.Deserialize(r, &data.LockTime)
		return bytes + n, err
	}
	return bytes, serialization.DeserializeTypeDismatchError
}

// SliceToString converts a byte slice into a string ends with '\0'
func SliceToString(data []byte) string {
	n := bytes.IndexByte(data, 0)
	if n != -1 {
		data = data[:n]
	}
	return string(data)
}
