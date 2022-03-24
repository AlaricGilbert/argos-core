package sniffer

import (
	"net"
	"time"
)

// TransactionNotify represents an abstract transaction which has been
type TransactionNotify struct {
	// SourceIP is the source where the daemon get notified
	SourceIP net.IP
	// Timestamp is the time when the daemon get notified
	Timestamp time.Time
	// TxID is the re-hashed abstract representation of an abstract transaction, which can be computed by real
	// implementation-related cryptocurrency transaction ids
	TxID []byte
}

type TransactionHandler func(tx TransactionNotify) error

// Daemon is an interface that describes the behaviour of an abstract cryptocurrency daemon in argos system
type Daemon interface {
	// Spin tries to connect the specified server and start spinning up the daemon packet handler.
	// It will never return until an error occurred.
	Spin() error
	// Halt will immediately stop the daemon handle procedure
	Halt() error
	// OnTransactionReceived make the daemon calls the given handler every time the daemon received a tx packet
	OnTransactionReceived(handler TransactionHandler)
}
