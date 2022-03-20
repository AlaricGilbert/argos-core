package protocol

import "errors"

var (
	// ConnectFailed means the trail of connecting into the server failed
	ConnectFailed = errors.New("connect failed")
	// Disconnected means the remote server has been disconnected
	Disconnected = errors.New("remote socket disconnected")
	// Halted means the client has been halted
	Halted = errors.New("client spinning halted")
	// NotRunning means the client not running
	NotRunning = errors.New("not running")
)
