package bitcoin

import (
	"sync"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/argos/serialization"
)

var once sync.Once

func initOnce() {
	serialization.RegisterDeserializer(deserialize)
}

func Init(ctx *argos.Context) error {
	once.Do(initOnce)
	ctx.RegisterClientConstructor("bitcoin", NewClient)
	ctx.RegisterSeedProvider("bitcoin", LookupBTCNetwork)
	return nil
}
