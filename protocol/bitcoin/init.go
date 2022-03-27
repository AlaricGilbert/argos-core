package bitcoin

import (
	"sync"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/argos/serialization"
)

var once sync.Once

func initOnce() {
	serialization.RegisterSerializer("bitcoin", &BitcoinSerializer{})
}

func Init(s argos.Sniffer) error {
	once.Do(initOnce)
	s.RegisterPeerConstructor("bitcoin", NewPeer)
	s.RegisterSeedProvider("bitcoin", LookupBTCNetwork)
	return nil
}
