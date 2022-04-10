package argos

import (
	"net"

	"github.com/sirupsen/logrus"
)

// Sniffer repesents an abstract super node that could connect to multiple nodes simultaneously
type Sniffer interface {
	Logger() *logrus.Logger
	NotifyTransaction(notify TransactionNotify)
	Connect(address net.TCPAddr)
	NodeConn(src net.TCPAddr, conn []net.TCPAddr)
	NodeExit(address net.TCPAddr)
	Spin()
	Halt()
}
