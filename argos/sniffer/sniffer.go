package sniffer

import (
	"net"

	"github.com/sirupsen/logrus"
)

type PeerConstructor func(ctx *Sniffer, addr *net.TCPAddr) Peer
type SeedProvider func() ([]net.IP, error)

// Sniffer repesents an abstract super node that could connect to multiple nodes simultaneously
type Sniffer struct {
	constructors  map[string]PeerConstructor
	seedProviders map[string]SeedProvider
	logger        *logrus.Logger
	transactions  chan TransactionNotify
	running       bool
}

func (c *Sniffer) RegisterPeerConstructor(name string, constructor PeerConstructor) {
	c.constructors[name] = constructor
}

func (c *Sniffer) RegisterSeedProvider(name string, provider SeedProvider) {
	c.seedProviders[name] = provider
}

func (c *Sniffer) NewPeer(protocol string, addr *net.TCPAddr) (Peer, error) {
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

func NewSniffer() *Sniffer {
	return &Sniffer{
		constructors:  make(map[string]PeerConstructor),
		seedProviders: make(map[string]SeedProvider),
		logger:        logrus.New(),
		transactions:  make(chan TransactionNotify),
		running:       false,
	}
}

func (s *Sniffer) Logger() *logrus.Logger {
	return s.logger
}

func (s *Sniffer) NotifyTransaction(notify TransactionNotify) {
	s.transactions <- notify
}

func (s *Sniffer) NodeConn(src net.TCPAddr, conn []net.TCPAddr) {

}

func (s *Sniffer) NodeExit(node net.Addr) {

}

func (s *Sniffer) Spin() {
	// only print tx into log currently
	s.running = true
	for s.running {
		if tx, ok := <-s.transactions; ok {
			s.Logger().WithField("tx", tx).Info("sniffer received tx")
		}
	}
}

func (s *Sniffer) Halt() {
	s.running = false
}
