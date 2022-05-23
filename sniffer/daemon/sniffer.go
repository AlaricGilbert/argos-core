package daemon

import (
	"net"
	"sync"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/graph"
	"github.com/sirupsen/logrus"
)

const ReportCenterThreshold = 24

type addr struct {
	IP   [16]byte
	Port int16
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
	network      *graph.Graph[addr, struct{}]
	notifies     map[[32]byte]map[addr]time.Time
	newAddrs     chan net.TCPAddr
	peers        map[addr]argos.Peer
	mu           sync.Mutex
}

func (s *Sniffer) Logger() *logrus.Logger {
	return s.logger
}

func (s *Sniffer) NotifyTransaction(notify argos.TransactionNotify) {
	s.mu.Lock()
	defer s.mu.Unlock()

	address := newAddr(notify.Source)

	// when get a transaction, check if it has been ignored
	// we use [nil] to mark the transaction has been ignored
	var notifies map[addr]time.Time
	var ok bool
	// first check the s.notifies map if it has been notified
	if notifies, ok = s.notifies[notify.TxID]; !ok {
		// never notified, create a new map, so the notifies variable is not nil
		// until we set it to [nil] after the transaction is marked ignored in ReportCenterEstimator.
		notifies = make(map[addr]time.Time)
		s.notifies[notify.TxID] = notifies
	}
	// then check the notifies map if the transaction has been ignored
	if notifies == nil {
		return
	}

	func() {
		if tt, ok := notifies[address]; ok && tt.Before(notify.Timestamp) {
			return
		}

		notifies[address] = notify.Timestamp
	}()

	// FirstTimestampEstimate is the first time a transaction is seen
	// by a node.

	if len(notifies) == 1 {
		go Report(notify.TxID[:], notify.Source.IP[:], notify.Source.Port, notify.Timestamp, "FTE")
	}

	if len(notifies) == ReportCenterThreshold {
		// generete a subgraph contains all the nodes that have been notified
		// get all the nodes in the subgraph
		var nodes = make([]addr, 0)
		var candicates = make([]addr, 0)
		var subgraph = graph.NewGraph[addr, struct{}]()
		for k := range notifies {
			// check the node is alive
			if _, ok := s.peers[k]; ok {
				nodes = append(nodes, k)
				subgraph.AddVertex(k, struct{}{})
			}
		}

		// get all the edges in the subgraph and add them to the subgraph
		for _, node := range nodes {
			for _, nn := range s.network.GetVertex(node).GetNeighbors() {
				if subgraph.ContainsVertex(nn) {
					subgraph.AddEdge(node, nn)
				}
			}
		}

		Yt := len(nodes)

		// run the ReportCenterEstimator in the subgraph
		// we thinks the subgraph is a k-degree regular tree use dfs to get the height
		var height = func(g *graph.Graph[addr, struct{}], center, m addr) int {
			h := 0
			g.DFS(m, func(key addr, value struct{}) bool {
				if key == center {
					return true
				}
				h += 1
				return true
			})
			return h
		}

		var max = func(a, b int) int {
			if a > b {
				return a
			}
			return b
		}

		for _, v := range subgraph.GetVertices() {
			if len(v.GetNeighbors()) == 0 {
				// if the vertex is a single node, then it is a candidate node
				candicates = append(candicates, v.GetKey())
			} else {
				// if the vertex is not a single node, then it is a center node
				// we need to use ReportCenterEstimator to check if it is a candidate node

				maxYt := 0
				for _, n := range v.GetNeighbors() {
					maxYt = max(maxYt, height(subgraph, v.GetKey(), n))
				}

				if maxYt < Yt/2 {
					candicates = append(candicates, v.GetKey())
				}
			}
		}

		// finally use first timestamp estimate to get a candidate node
		var candidate *addr
		var ts = time.Now()
		for _, c := range candicates {
			if tt, ok := notifies[c]; ok && tt.Before(ts) {
				candidate = &c
				ts = tt
			}
		}

		go Report(notify.TxID[:], candidate.IP[:], int(candidate.Port), ts, "RCE")

		// set the notifies map to [nil]
		s.notifies[notify.TxID] = nil
	}
}

func (s *Sniffer) NodeConn(src net.TCPAddr, conn []net.TCPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()

	srcAddr := newAddr(src)
	s.network.AddVertex(srcAddr, struct{}{})

	connAddrs := make([]addr, len(conn))
	for i, addr := range conn {
		connAddrs[i] = newAddr(addr)
		if srcAddr == connAddrs[i] {
			continue
		}
		s.network.AddVertex(connAddrs[i], struct{}{})
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

func (s *Sniffer) Spin(node net.TCPAddr) {
	// query protocol
	var protocol = "bitcoin"

	nodes, err := argos.GetSeedNodes(protocol)
	if err != nil {
		s.logger.WithError(err).Fatal("get seed nodes failed")
		return
	}

	// only print tx into log currently
	s.running = true
	go func() {
		for _, node := range nodes {
			s.Connect(node)
		}
	}()
	s.hostConnection()
}

func (s *Sniffer) hostConnection() {
	for address := range s.newAddrs {
		s.Connect(address)
	}
}

func (s *Sniffer) Connect(address net.TCPAddr) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var err error
	var peer argos.Peer
	var addr = newAddr(address)
	if _, ok := s.peers[addr]; ok {
		s.logger.WithField("address", address).Error("sniffer already connected to peer")
		return
	}

	if peer, err = argos.NewPeer(Instance().protocol, &address, s); err != nil {
		delete(s.peers, addr)
		s.logger.WithField("address", address).WithError(err).Error("failed to connect to peer")
	} else {
		s.network.AddVertex(addr, struct{}{})
		s.peers[addr] = peer
		go func() {
			peer.Spin()
			// delete peer
			s.mu.Lock()
			defer s.mu.Unlock()
			delete(s.peers, addr)
			s.network.RemoveVertex(addr)
		}()
	}
}

func (s *Sniffer) Halt() {
	s.running = false
}

func NewSniffer(logger *logrus.Logger) *Sniffer {
	return &Sniffer{
		transactions: make(chan argos.TransactionNotify),
		newAddrs:     make(chan net.TCPAddr, 1000),
		notifies:     make(map[[32]byte]map[addr]time.Time),
		network:      graph.NewGraph[addr, struct{}](),
		peers:        make(map[addr]argos.Peer),
		running:      false,
		logger:       logger,
	}
}
