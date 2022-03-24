package main

import (
	"github.com/AlaricGilbert/argos-core/argos/sniffer"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
)

func main() {
	ctx := sniffer.NewContext()
	if bitcoin.Init(ctx) != nil {
		panic("bitcoin init failed")
	}
}
