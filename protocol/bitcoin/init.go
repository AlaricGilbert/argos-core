package bitcoin

import (
	"sync"

	"github.com/AlaricGilbert/argos-core/argos/serialization"
	"github.com/AlaricGilbert/argos-core/argos/sniffer"
)

var once sync.Once

func initOnce() {
	serialization.RegisterSerializer("bitcoin", &BitcoinSerializer{})
}

func Init(s *sniffer.Sniffer) error {
	once.Do(initOnce)
	s.RegisterDaemonConstructor("bitcoin", NewDaemon)
	s.RegisterSeedProvider("bitcoin", LookupBTCNetwork)
	return nil
}
