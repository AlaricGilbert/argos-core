package argos

import (
	"net"
	"time"
)

// TransactionNotify represents an abstract transaction which has been
type TransactionNotify struct {
	// Source is the source where the current node get notified
	Source net.TCPAddr
	// Timestamp is the time when the current node get notified
	Timestamp time.Time
	// TxID is the re-hashed abstract representation of an abstract transaction, which can be computed by real
	// implementation-related cryptocurrency transaction ids
	TxID [32]byte
}

// Peer is an interface that describes the behaviour of an abstract cryptocurrency peer in argos system
type Peer interface {
	// Spin tries to connect the specified server and start spinning up the peer packet handler.
	// It will never return until an error occurred.
	Spin() error
	// Halt will immediately stop the peer handle procedure
	Halt() error
}
