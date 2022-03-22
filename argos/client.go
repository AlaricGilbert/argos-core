package argos

import (
	"net"
	"time"
)

// TransactionNotify represents an abstract transaction which has been
type TransactionNotify struct {
	// SourceIP is the source where the client get notified
	SourceIP net.IP
	// Timestamp is the time when the client get notified
	Timestamp time.Time
	// TxID is the re-hashed abstract representation of an abstract transaction, which can be computed by real
	// implementation-related cryptocurrency transaction ids
	TxID []byte
}

type TransactionHandler func(tx TransactionNotify) error

// Client is an interface that describes the behaviour of an abstract cryptocurrency client in argos system
type Client interface {
	// Spin tries to connect the specified server and start spinning up the client packet handler.
	// It will never return until an error occurred.
	Spin() error
	// Halt will immediately stop the client handle procedure
	Halt() error
	// OnTransactionReceived make the client calls the given handler every time the client received a tx packet
	OnTransactionReceived(handler TransactionHandler)
}
