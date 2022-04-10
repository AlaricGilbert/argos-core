package argos

import "errors"

var (
	// ErrProtocolNotImplemented
	ErrProtocolNotImplemented = errors.New("protocol not implemented")
	// ErrConnectFailed means the trail of connecting into the server failed
	ErrConnectFailed = errors.New("connect failed")
	// ErrDisconnected means the remote server has been disconnected
	ErrDisconnected = errors.New("remote socket disconnected")
	// ErrPeerHalted means the peer has been halted
	ErrPeerHalted = errors.New("peer spinning halted")
	// ErrPeerNotRunning means the peer not running
	ErrPeerNotRunning = errors.New("peer not running")
)
