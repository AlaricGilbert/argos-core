package protocol

import "errors"

var (
	// ConnectFailedError means the trail of connecting into the server failed
	ConnectFailedError = errors.New("connect failed")
	// DisconnectedError means the remote server has been disconnected
	DisconnectedError = errors.New("remote socket disconnected")
	// ClientHaltedError means the client has been halted
	ClientHaltedError = errors.New("client spinning halted")
	// ClientNotRunningError means the client not running
	ClientNotRunningError = errors.New("not running")
)
