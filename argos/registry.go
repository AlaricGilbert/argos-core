package argos

import "net"

type PeerConstructor func(s Sniffer, addr *net.TCPAddr) Peer
type SeedProvider func() ([]net.TCPAddr, error)
type RandomRemoteAddressProvider func() (*net.TCPAddr, error)

var (
	constructors                 = make(map[string]PeerConstructor)
	seedProviders                = make(map[string]SeedProvider)
	randomRemoteAddressProviders = make(map[string]RandomRemoteAddressProvider)
)

func RegisterPeerConstructor(name string, constructor PeerConstructor) {
	constructors[name] = constructor
}

func RegisterSeedProvider(name string, provider SeedProvider) {
	seedProviders[name] = provider
}

func RegisterRandomRemoteAddressProvider(name string, provider RandomRemoteAddressProvider) {
	randomRemoteAddressProviders[name] = provider
}

func NewPeer(protocol string, addr *net.TCPAddr, s Sniffer) (Peer, error) {
	if ctor, ok := constructors[protocol]; ok {
		return ctor(s, addr), nil
	}
	return nil, ErrProtocolNotImplemented
}

func GetSeedNodes(protocol string) ([]net.TCPAddr, error) {
	if provider, ok := seedProviders[protocol]; ok {
		return provider()
	}
	return nil, ErrProtocolNotImplemented
}

func GetRandomRemoteAddress(protocol string) (*net.TCPAddr, error) {
	if provider, ok := randomRemoteAddressProviders[protocol]; ok {
		return provider()
	}
	return nil, ErrProtocolNotImplemented
}

func GetSupportedProtocols() []string {
	protocols := make([]string, 0, len(constructors))
	for k := range constructors {
		protocols = append(protocols, k)
	}
	return protocols
}
