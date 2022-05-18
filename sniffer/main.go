package main

import (
	"github.com/AlaricGilbert/argos-core/sniffer/daemon"
)

func main() {
	daemon.Init()
	daemon.Instance().Spin()
}
