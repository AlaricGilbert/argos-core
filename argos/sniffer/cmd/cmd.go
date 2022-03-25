package main

import (
	"net"

	"github.com/AlaricGilbert/argos-core/argos/sniffer"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// currently works as a bitcoin daemon
func main() {
	s := sniffer.NewSniffer()
	if bitcoin.Init(s) != nil {
		panic("bitcoin init failed")
	}

	s.Logger().AddHook(lfshook.NewHook(
		lfshook.PathMap{
			logrus.InfoLevel:  "logs.log",
			logrus.ErrorLevel: "logs.log",
			logrus.WarnLevel:  "logs.log",
		},
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

	d, _ := s.NewDaemon("bitcoin", &net.TCPAddr{
		IP:   ip,
		Port: 8333,
	})

	go func() {
		for true {
			n := <-s.Transactions

			logrus.Infof("btc tx notify: %v", n)
		}
	}()

	d.Spin()
}
