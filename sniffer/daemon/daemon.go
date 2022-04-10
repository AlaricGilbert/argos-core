package daemon

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/master"
	am "github.com/AlaricGilbert/argos-core/master/kitex_gen/master/argosmaster"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
	"github.com/cloudwego/kitex/client"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type SnifferDaemon struct {
	logger  *logrus.Logger
	sniffer argos.Sniffer
	master  am.Client
	config  *Config
}

var instance *SnifferDaemon

// ping is a loop that pings the master every 10 seconds, to keep the connection to the master alive
// and to get the task provided by the master
func (d *SnifferDaemon) ping() {
	var errTimes = 0
	for {
		var req *master.PingRequest
		var resp *master.PingResponse
		var err error
		time.Sleep(time.Second * 10)

		// when errTimes > 10, we think the connection to the master is broken, so we will exit the loop
		if errTimes > 10 {
			d.logger.Fatal("ping failed 10 times, argos master is not available")
		}

		req = &master.PingRequest{
			Protocols:  argos.GetSupportedProtocols(),
			Identifier: d.config.Identifier,
		}

		if resp, err = d.master.Ping(context.Background(), req); err != nil {
			d.logger.WithError(err).Error("argos sniffer ping failed")
			errTimes++
			continue
		}

		if resp.Status.Code != 0 {
			d.logger.WithField("status", resp.Status).Fatal("argos sniffer ping failed")
			errTimes++
			continue
		}

	}
}

// GetLogger returns the logger of the daemon
func (d *SnifferDaemon) GetLogger() *logrus.Logger {
	return d.logger
}

// GetSniffer returns the sniffer of the daemon
func (d *SnifferDaemon) GetSniffer() argos.Sniffer {
	return d.sniffer
}

// GetMasterClient returns the master client of the daemon
func (d *SnifferDaemon) GetMasterClient() am.Client {
	return d.master
}

// GetConfig returns the config of the daemon
func (d *SnifferDaemon) GetConfig() *Config {
	return d.config
}

func (d *SnifferDaemon) GetNodes() []net.TCPAddr {
	return nil
}

// Spin runs the ping and sniffer loop
func (d *SnifferDaemon) Spin() {
	go d.ping()
}

func Init() {
	if instance != nil {
		instance.logger.Fatal("argos sniffer daemon already initialized")
	}

	var err error

	// it should be noticed that logger should be initialized BEFORE any other global consts
	instance.logger = logrus.New()
	instance.logger.AddHook(lfshook.NewHook(
		fmt.Sprintf("logs/%s.log", time.Now().Format(time.RFC3339)),
		&logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		},
	))

	argos.SetLogger(instance.logger)

	// read config
	if err = instance.ReadConfig(); err != nil {
		instance.logger.WithError(err).Fatal("read config failed")
	}

	instance.sniffer = NewSniffer(instance.logger)

	if err = bitcoin.Init(instance.sniffer); err != nil {
		instance.logger.WithError(err).Fatal("bitcoin init failed")
	}

	if instance.master, err = am.NewClient("argos.master", client.WithHostPorts(instance.config.MasterAddress)); err != nil {
		instance.logger.WithError(err).Fatal("argos master client init failed")
	}

	resp, err := instance.master.Ping(context.Background(), &master.PingRequest{
		Identifier: instance.config.Identifier,
		Protocols:  argos.GetSupportedProtocols(),
	})

	if err != nil {
		instance.logger.WithError(err).Fatal("argos sniffer register failed")
	}

	if resp.Status.Code != 0 {
		instance.logger.WithField("status", resp.Status).Fatal("argos sniffer register failed")
	}
}

func Instance() *SnifferDaemon {
	if instance == nil {
		panic("argos sniffer daemon not initialized")
	}
	return instance
}
