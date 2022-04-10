package main

import (
	"context"

	"github.com/AlaricGilbert/argos-core/sniffer/handler"
	"github.com/AlaricGilbert/argos-core/sniffer/kitex_gen/sniffer"
)

// ArgosSnifferImpl implements the last service interface defined in the IDL.
type ArgosSnifferImpl struct{}

// Time implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) Time(ctx context.Context, req *sniffer.TimeRequest) (resp *sniffer.TimeResponse, err error) {
	return handler.Time(ctx, req)
}

// Transaction implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) Transaction(ctx context.Context, req *sniffer.TimeRequest) (resp *sniffer.TransactionResponse, err error) {
	return handler.Transaction(ctx, req)
}

// Address implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) Address(ctx context.Context, req *sniffer.AddressRequest) (resp *sniffer.AddressResponse, err error) {
	return handler.Address(ctx, req)
}

// TransactionNotify implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) TransactionNotify(ctx context.Context, req *sniffer.TransactionNotify) (err error) {
	// TODO: Your code here...
	return
}
