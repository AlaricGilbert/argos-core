package argos

import "errors"

var (
	// ProtocolNotImplementedError
	ProtocolNotImplementedError = errors.New("protocol not implemented")
	// ConnectFailedError means the trail of connecting into the server failed
	ConnectFailedError = errors.New("connect failed")
	// DisconnectedError means the remote server has been disconnected
	DisconnectedError = errors.New("remote socket disconnected")
	// PeerHaltedError means the peer has been halted
	PeerHaltedError = errors.New("peer spinning halted")
	// PeerNotRunningError means the peer not running
	PeerNotRunningError = errors.New("peer not running")
)
