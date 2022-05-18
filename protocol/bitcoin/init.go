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

func Init() error {
	once.Do(initOnce)
	argos.RegisterPeerConstructor("bitcoin", NewPeer)
	argos.RegisterSeedProvider("bitcoin", LookupBTCNetwork)
	argos.RegisterRandomRemoteAddressProvider("bitcoin", LookupRandomBTCHostAddress)
	return nil
}
