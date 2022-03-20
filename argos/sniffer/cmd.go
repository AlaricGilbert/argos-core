package main

import (
	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
)

func main() {
	ctx := argos.NewContext()
	if bitcoin.Init(ctx) != nil {
		panic("bitcoin init failed")
	}
}
