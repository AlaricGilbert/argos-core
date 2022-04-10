package daemon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Config struct {
	MasterAddress string `json:"master_address"`
	LocalPort     int    `json:"local_port"`
	Identifier    string `json:"identifier"`
}

func randIdentifier() string {
	m := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	l := len(m)

	randRune := func() byte {
		return m[rng.Intn(l)]
	}

	var b strings.Builder
	for i := 0; i < 10; i++ {
		_ = b.WriteByte(randRune())
	}

	return b.String()
}

func newDefaultConfig() *Config {
	return &Config{
		MasterAddress: "127.0.0.1:4222",
		LocalPort:     8777,
		Identifier:    randIdentifier(),
	}
}

func (d *SnifferDaemon) ReadConfig() error {
	var err error
	defer d.logger.WithError(err).WithField("config", d.config).Warn("read config exited")
	var bytes []byte
	var cfgFile = "sniffer.json"
	var f *os.File

	d.config = newDefaultConfig()

	if _, err = os.Stat(cfgFile); errors.Is(err, os.ErrNotExist) {
		d.SaveConfig()
		return nil
	}

	if f, err = os.OpenFile(cfgFile, os.O_RDONLY, 0); err != nil {
		return err
	}

	if bytes, err = ioutil.ReadAll(f); err != nil {
		return err
	}

	if err = json.Unmarshal(bytes, d.config); err != nil {
		return err
	}

	// if config.Identifier empty, generate a new one
	if d.config.Identifier == "" {
		d.config.Identifier = randIdentifier()
		d.SaveConfig()
	}

	return nil
}

func (d *SnifferDaemon) SaveConfig() error {
	var bytes []byte
	var err error
	var cfgFile = "sniffer.json"
	var f *os.File
	if bytes, err = json.MarshalIndent(d.config, "", "    "); err != nil {
		d.logger.WithError(err).Warn("marshal config failed")
		return err
	}

	if f, err = os.OpenFile(cfgFile, os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		d.logger.WithError(err).Warn("create config and open failed")
		return err
	}

	if _, err = f.Write(bytes); err != nil {
		d.logger.WithError(err).Warn("write config failed")
	}
	return nil
}
