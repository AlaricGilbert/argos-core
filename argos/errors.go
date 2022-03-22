package argos

import "errors"

var (
	// ProtocolNotImplementedError
	ProtocolNotImplementedError = errors.New("protocol not implemented")
	// ConnectFailedError means the trail of connecting into the server failed
	ConnectFailedError = errors.New("connect failed")
	// DisconnectedError means the remote server has been disconnected
	DisconnectedError = errors.New("remote socket disconnected")
	// DaemonHaltedError means the daemon has been halted
	DaemonHaltedError = errors.New("daemon spinning halted")
	// DaemonNotRunningError means the daemon not running
	DaemonNotRunningError = errors.New("daemon not running")
)
