package daemon

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/base"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/master"
	am "github.com/AlaricGilbert/argos-core/master/kitex_gen/master/argosmaster"
	"github.com/AlaricGilbert/argos-core/protocol/bitcoin"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/cloudwego/kitex/client"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

type SnifferDaemon struct {
	logger    *logrus.Logger
	sniffer   argos.Sniffer
	master    am.Client
	protocol  string
	timeDelta int64
	config    *Config
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
			Identifier: d.config.Identifier,
			Timestamp:  time.Now().UnixNano(),
			DeltaTime:  thrift.Int64Ptr(d.timeDelta),
		}

		if resp, err = d.master.Ping(context.Background(), req); err != nil {
			d.logger.WithError(err).Error("argos sniffer ping failed")
			errTimes++
			continue
		}

		if resp == nil || resp.Status.Code != 0 {
			d.logger.WithField("status", resp.Status).Fatal("argos sniffer ping failed")
			errTimes++
			continue
		}

		// if the master is available, we will sync the time with master
		d.timeDelta = (resp.GetTimeSync().RecvTimestamp - resp.GetTimeSync().SendTimestamp) / 2

		// check the if the protocol has changed, if so, we will exit the sniffer and restart it
		// it should be noticed that restart should be done in some shell scripts, not in the sniffer codes
		if d.protocol != resp.GetProtocol() {
			d.logger.WithField("protocol", resp.GetProtocol()).Info("protocol changed, restarting sniffer")
			os.Exit(0)
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
	// start the ping loop
	go d.ping()

	// get the init node

	if addr, err := argos.GetRandomRemoteAddress(d.protocol); err != nil {
		d.logger.WithError(err).Fatal("get seed nodes failed")
		os.Exit(-1)
	} else {
		// start the sniffer loop
		d.sniffer.Spin(*addr)
	}
}

func Init() {
	if instance != nil {
		instance.logger.Fatal("argos sniffer daemon already initialized")
	}

	instance = &SnifferDaemon{}

	var err error

	// it should be noticed that logger should be initialized BEFORE any other global consts except the daemon instance
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

	if err = bitcoin.Init(); err != nil {
		instance.logger.WithError(err).Fatal("bitcoin init failed")
	}

	if instance.master, err = am.NewClient("argos.master", client.WithHostPorts(instance.config.MasterAddress)); err != nil {
		instance.logger.WithError(err).Fatal("argos master client init failed")
	}

	resp, err := instance.master.Ping(context.Background(), &master.PingRequest{
		Identifier: instance.config.Identifier,
		Timestamp:  time.Now().UnixNano(),
	})

	if err != nil {
		instance.logger.WithError(err).Fatal("argos sniffer register failed")
	}

	if resp == nil || resp.Status.Code != 0 {
		instance.logger.WithField("status", resp.Status).Fatal("argos sniffer register failed")
	}

	// save time delta and protocol
	instance.timeDelta = (resp.GetTimeSync().RecvTimestamp - resp.GetTimeSync().SendTimestamp) / 2
	instance.protocol = resp.GetProtocol()
}

func Instance() *SnifferDaemon {
	if instance == nil {
		panic("argos sniffer daemon not initialized")
	}
	return instance
}

func Report(txid []byte, ip []byte, port int, timestamp time.Time, method string) {
	if instance == nil {
		panic("argos sniffer daemon not initialized")
	}

	instance.master.Report(context.Background(), &master.ReportRequest{
		Identifier: Instance().config.Identifier,
		Method:     method,
		Transaction: &base.Transaction{
			Txid:      txid,
			Timestamp: timestamp.UnixNano() + instance.timeDelta,
			From: &base.TcpAddress{
				Ip:   ip,
				Port: int32(port),
			},
		},
		Protocol: instance.protocol,
	})
}
