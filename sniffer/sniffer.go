package main

import (
	"net"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/sirupsen/logrus"
)

type Sniffer struct {
	constructors  map[string]argos.PeerConstructor
	seedProviders map[string]argos.SeedProvider
	logger        *logrus.Logger
	transactions  chan argos.TransactionNotify
	running       bool
}

func (c *Sniffer) RegisterPeerConstructor(name string, constructor argos.PeerConstructor) {
	c.constructors[name] = constructor
}

func (c *Sniffer) RegisterSeedProvider(name string, provider argos.SeedProvider) {
	c.seedProviders[name] = provider
}

func (c *Sniffer) NewPeer(protocol string, addr *net.TCPAddr) (argos.Peer, error) {
	if ctor, ok := c.constructors[protocol]; ok {
		return ctor(c, addr), nil
	}
	return nil, argos.ProtocolNotImplementedError
}

func (c *Sniffer) GetSeedNodes(protocol string) ([]net.IP, error) {
	if provider, ok := c.seedProviders[protocol]; ok {
		return provider()
	}
	return nil, argos.ProtocolNotImplementedError
}

func NewSniffer() *Sniffer {
	return &Sniffer{
		constructors:  make(map[string]argos.PeerConstructor),
		seedProviders: make(map[string]argos.SeedProvider),
		logger:        logrus.New(),
		transactions:  make(chan argos.TransactionNotify),
		running:       false,
	}
}

func (s *Sniffer) Logger() *logrus.Logger {
	return s.logger
}

func (s *Sniffer) NotifyTransaction(notify argos.TransactionNotify) {
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
