package argos

import (
	"net"

	"github.com/sirupsen/logrus"
)

type PeerConstructor func(ctx Sniffer, addr *net.TCPAddr) Peer
type SeedProvider func() ([]net.IP, error)

// Sniffer repesents an abstract super node that could connect to multiple nodes simultaneously
type Sniffer interface {
	RegisterPeerConstructor(name string, constructor PeerConstructor)
	RegisterSeedProvider(name string, provider SeedProvider)
	NewPeer(protocol string, addr *net.TCPAddr) (Peer, error)
	GetSeedNodes(protocol string) ([]net.IP, error)
	Logger() *logrus.Logger
	NotifyTransaction(notify TransactionNotify)
	NodeConn(src net.TCPAddr, conn []net.TCPAddr)
	NodeExit(node net.Addr)
	Spin()
	Halt()
}
