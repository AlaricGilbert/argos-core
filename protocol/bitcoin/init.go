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

func Init(ctx *argos.Context) error {
	once.Do(initOnce)
	ctx.RegisterDaemonConstructor("bitcoin", NewDaemon)
	ctx.RegisterSeedProvider("bitcoin", LookupBTCNetwork)
	return nil
}
