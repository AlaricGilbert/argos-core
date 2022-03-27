package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
	api "github.com/AlaricGilbert/argos-core/sniffer/kitex_gen/api/argossniffer"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func main() {
	s := NewSniffer()
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
	go func() {
		d.Spin()
		s.Halt()
	}()

	klog.SetLogger(argos.New(s.logger))

	svr := api.NewServer(new(ArgosSnifferImpl))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
