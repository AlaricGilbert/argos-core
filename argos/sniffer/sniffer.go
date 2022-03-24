package sniffer

import (
	"net"
)

type DaemonConstructor func(ctx *Sniffer, addr *net.TCPAddr) Daemon
type SeedProvider func() ([]net.IP, error)

// Sniffer repesents an abstract super node that could connect to multiple nodes simultaneously
type Sniffer struct {
	constructors  map[string]DaemonConstructor
	seedProviders map[string]SeedProvider
	// Transactions are used for daemons to report transaction notifies
	Transactions chan TransactionNotify
}

func (c *Sniffer) RegisterDaemonConstructor(name string, constructor DaemonConstructor) {
	c.constructors[name] = constructor
}

func (c *Sniffer) RegisterSeedProvider(name string, provider SeedProvider) {
	c.seedProviders[name] = provider
}

func (c *Sniffer) NewDaemon(protocol string, addr *net.TCPAddr) (Daemon, error) {
	if ctor, ok := c.constructors[protocol]; ok {
		return ctor(c, addr), nil
	}
	return nil, ProtocolNotImplementedError
}

func (c *Sniffer) GetSeedNodes(protocol string) ([]net.IP, error) {
	if provider, ok := c.seedProviders[protocol]; ok {
		return provider()
	}
	return nil, ProtocolNotImplementedError
}

func NewContext() *Sniffer {
	return &Sniffer{
		constructors:  make(map[string]DaemonConstructor),
		seedProviders: make(map[string]SeedProvider),
		Transactions:  make(chan TransactionNotify, 32),
	}
}
