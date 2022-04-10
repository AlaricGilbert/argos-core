package handler

import (
	"context"
	"time"

	"github.com/AlaricGilbert/argos-core/sniffer/daemon"
	"github.com/AlaricGilbert/argos-core/sniffer/kitex_gen/base"
	"github.com/AlaricGilbert/argos-core/sniffer/kitex_gen/sniffer"
)

func Time(ctx context.Context, req *sniffer.TimeRequest) (resp *sniffer.TimeResponse, err error) {
	ts := time.Now().UnixNano()
	resp = &sniffer.TimeResponse{
		SendTimestamp: req.Timestamp,
		RecvTimestamp: ts,
		RespTimestamp: ts,
	}
	return
}

func Transaction(ctx context.Context, req *sniffer.TimeRequest) (resp *sniffer.TransactionResponse, err error) {
	// TODO: Your code here...
	return
}

func Address(ctx context.Context, req *sniffer.AddressRequest) (resp *sniffer.AddressResponse, err error) {
	resp = &sniffer.AddressResponse{}
	addresses := daemon.Instance().GetNodes()
	resp.Addresses = make([]*base.TcpAddress, len(addresses))
	for i, addr := range addresses {
		resp.Addresses[i] = &base.TcpAddress{
			Ip:   []byte(addr.IP),
			Port: int32(addr.Port),
		}
	}
	return
}
