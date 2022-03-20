package argos

import (
	"github.com/AlaricGilbert/argos-core/protocol"
	"net"
)

type ClientConstructor func(ctx *Context, ip net.IP, port int) protocol.Client
type SeedProvider func() ([]net.IP, error)

// Context contains
type Context struct {
	constructors  map[string]ClientConstructor
	seedProviders map[string]SeedProvider
}

func (c *Context) RegisterClientConstructor(name string, constructor ClientConstructor) {
	c.constructors[name] = constructor
}

func (c *Context) RegisterSeedProvider(name string, provider SeedProvider) {
	c.seedProviders[name] = provider
}

func (c *Context) NewClient(network string, ip net.IP, port int) (protocol.Client, error) {
	if ctor, ok := c.constructors[network]; ok {
		return ctor(c, ip, port), nil
	}
	return nil, ProtocolNotImplemented
}

func (c *Context) GetSeedNodes(network string) ([]net.IP, error) {
	if provider, ok := c.seedProviders[network]; ok {
		return provider()
	}
	return nil, ProtocolNotImplemented
}

func NewContext() *Context {
	return &Context{
		constructors:  make(map[string]ClientConstructor),
		seedProviders: make(map[string]SeedProvider),
	}
}
