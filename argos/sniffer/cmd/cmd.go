package main

import (
	"fmt"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos/sniffer"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// currently works as a bitcoin peer
func main() {
	s := sniffer.NewSniffer()
	if bitcoin.Init(s) != nil {
		panic("bitcoin init failed")
	}

	s.Logger().AddHook(lfshook.NewHook(
		fmt.Sprintf("logs/%s.log", time.Now().Format(time.RFC3339)),
		&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		},
	))
	bitcoin.Init(s)

	ip, err := bitcoin.LookupRandomBTCNetwork()
	if err != nil {
		panic(err)
	}

	d, _ := s.NewPeer("bitcoin", &net.TCPAddr{
		IP:   ip,
		Port: 8333,
	})
	go s.Spin()
	d.Spin()
	s.Halt()
}
