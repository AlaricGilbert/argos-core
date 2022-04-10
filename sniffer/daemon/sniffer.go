package daemon

import (
	"net"
	"sync"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/graph"
	"github.com/sirupsen/logrus"
)

type addr struct {
	IP   [16]byte
	Port int16
}

type node struct {
	notified map[[32]byte]time.Time
	mu       sync.Mutex
}

func newNode() *node {
	return &node{
		notified: make(map[[32]byte]time.Time),
		mu:       sync.Mutex{},
	}
}

func (n *node) notify(tx [32]byte, t time.Time) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if tt, ok := n.notified[tx]; ok {
		if t.After(tt) {
			return
		}
	}
	n.notified[tx] = t
}

func newAddr(address net.TCPAddr) addr {
	var ip [16]byte
	copy(ip[:], address.IP)
	return addr{
		IP:   ip,
		Port: int16(address.Port),
	}
}

type Sniffer struct {
	transactions chan argos.TransactionNotify
	running      bool
	logger       *logrus.Logger
	network      *graph.Graph[addr, *node]
	newAddrs     chan net.TCPAddr
	peers        map[addr]argos.Peer
	mu           sync.Mutex
	protocol     string
}

func (s *Sniffer) Logger() *logrus.Logger {
	return s.logger
}

func (s *Sniffer) NotifyTransaction(notify argos.TransactionNotify) {
	s.mu.Lock()
	defer s.mu.Unlock()

	address := newAddr(notify.Source)
	v := s.network.GetVertex(address)
	if v == nil {
		s.logger.WithField("notify", notify).Error("sniffer received tx from unknown node")
		return
	}

	v.GetValue().notify(notify.TxID, time.Now())
}

func (s *Sniffer) NodeConn(src net.TCPAddr, conn []net.TCPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	srcAddr := newAddr(src)
	s.network.AddVertexWithFactory(srcAddr, newNode)

	connAddrs := make([]addr, len(conn))
	for i, addr := range conn {
		connAddrs[i] = newAddr(addr)
		s.network.AddVertexWithFactory(connAddrs[i], newNode)
		s.network.AddEdge(srcAddr, connAddrs[i])
		s.newAddrs <- addr
	}
}

func (s *Sniffer) NodeExit(address net.TCPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	addr := newAddr(address)
	s.network.RemoveVertex(addr)
}

func (s *Sniffer) Spin() {
	// only print tx into log currently
	s.running = true
	go s.hostConnection()
}

func (s *Sniffer) hostConnection() {
	for {
		select {
		case address := <-s.newAddrs:
			s.Connect(address)
		}
	}
}

func (s *Sniffer) Connect(address net.TCPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var err error
	if _, ok := s.peers[newAddr(address)]; ok {
		s.logger.WithField("address", address).Error("sniffer already connected to peer")
		return
	}

	if s.peers[newAddr(address)], err = argos.NewPeer(s.protocol, &address, s); err != nil {
		delete(s.peers, newAddr(address))
		s.logger.WithField("address", address).WithError(err).Error("failed to connect to peer")
	}
}

func (s *Sniffer) Halt() {
	s.running = false
}

func NewSniffer(logger *logrus.Logger) *Sniffer {
	return &Sniffer{
		transactions: make(chan argos.TransactionNotify),
		running:      false,
		logger:       logger,
	}
}
