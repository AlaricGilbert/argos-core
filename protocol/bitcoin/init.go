package bitcoin

import (
	"sync"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/argos/serialization"
)

var once sync.Once

func initOnce() {
	_ = serialization.RegisterSerializer("bitcoin", &BitcoinSerializer{})
}

func Init(s argos.Sniffer) error {
	once.Do(initOnce)
	argos.RegisterPeerConstructor("bitcoin", NewPeer)
	argos.RegisterSeedProvider("bitcoin", LookupBTCNetwork)
	return nil
}
