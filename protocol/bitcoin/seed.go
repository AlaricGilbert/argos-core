package bitcoin

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
)

// The DNS host from https://github.com/bitcoin/bitcoin core repository
// When started for the first time, programs don’t know the IP addresses of any active full nodes. In order to discover
// some IP addresses, we query one or more DNS names hardcoded here. The response to the lookup should include one or
// more DNS A records with the IP addresses of full nodes that may accept new incoming connections.
var btcSeedHosts = []string{
	"seed.bitcoin.sipa.be.",          // Pieter Wuille, only supports x1, x5, x9, and xd
	"dnsseed.bluematt.me.",           // Matt Corallo, only supports x9
	"dnsseed.bitcoin.dashjr.org.",    // Luke Dashjr
	"seed.bitcoinstats.com.",         // Christian Decker, supports x1 - xf
	"seed.bitcoin.jonasschnelli.ch.", // Jonas Schnelli, only supports x1, x5, x9, and xd
	"seed.btc.petertodd.org.",        // Peter Todd, only supports x1, x5, x9, and xd
	"seed.bitcoin.sprovoost.nl.",     // Sjors Provoost
	"dnsseed.emzy.de.",               // Stephan Oeste
	"seed.bitcoin.wiz.biz.",          // Jason Maurice
}

var seedRng = rand.New((rand.NewSource(time.Now().Unix())))

// LookupBTCNetwork queries all the possible connected DNS servers and returns the core BTC network seed.
func LookupBTCNetwork() ([]net.TCPAddr, error) {
	var result = make([]net.IP, 0)
	var nodes = make([]net.TCPAddr, 0)
	for _, host := range btcSeedHosts {
		ips, err := net.LookupIP(host)
		if err != nil {
			argos.StandardLogger().Errorf("[LookupBTCNetwork] Get BTC DNS seed from host `%s` failed: %v", host, err)
			return nil, errors.New(fmt.Sprintf("Get BTC DNS seed from host `%s` failed: %v", host, err))
		}
		result = append(result, ips...)
	}
	for _, ip := range result {
		nodes = append(nodes, net.TCPAddr{IP: ip, Port: 8333})
	}
	return nodes, nil
}

func LookupRandomBTCNetwork() (net.IP, error) {
	if ips, err := net.LookupIP(btcSeedHosts[seedRng.Intn(len(btcSeedHosts))]); err != nil {
		return nil, err
	} else {
		return ips[seedRng.Intn(len(ips))], nil
	}
}

func LookupRandomBTCHostAddress() (*net.TCPAddr, error) {
	if ip, err := LookupRandomBTCNetwork(); err != nil {
		return nil, err
	} else {
		return &net.TCPAddr{IP: ip, Port: 8333}, nil
	}
}
