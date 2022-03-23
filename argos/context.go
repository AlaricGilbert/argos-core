package argos

import (
	"net"
)

type DaemonConstructor func(ctx *Context, addr *net.TCPAddr) Daemon
type SeedProvider func() ([]net.IP, error)

// Context contains
type Context struct {
	constructors  map[string]DaemonConstructor
	seedProviders map[string]SeedProvider
}

func (c *Context) RegisterDaemonConstructor(name string, constructor DaemonConstructor) {
	c.constructors[name] = constructor
}

func (c *Context) RegisterSeedProvider(name string, provider SeedProvider) {
	c.seedProviders[name] = provider
}

func (c *Context) NewDaemon(network string, addr *net.TCPAddr) (Daemon, error) {
	if ctor, ok := c.constructors[network]; ok {
		return ctor(c, addr), nil
	}
	return nil, ProtocolNotImplementedError
}

func (c *Context) GetSeedNodes(network string) ([]net.IP, error) {
	if provider, ok := c.seedProviders[network]; ok {
		return provider()
	}
	return nil, ProtocolNotImplementedError
}

func NewContext() *Context {
	return &Context{
		constructors:  make(map[string]DaemonConstructor),
		seedProviders: make(map[string]SeedProvider),
	}
}
