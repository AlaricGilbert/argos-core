package main

import (
	"log"

	"github.com/AlaricGilbert/argos-core/sniffer/daemon"
	sniffer "github.com/AlaricGilbert/argos-core/sniffer/kitex_gen/sniffer/argossniffer"
)

func main() {
	daemon.Init()
	go daemon.Instance().Spin()

	svr := sniffer.NewServer(new(ArgosSnifferImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
