package main

import (
	"context"
	"github.com/AlaricGilbert/argos-core/sniffer/kitex_gen/api"
)

// ArgosSnifferImpl implements the last service interface defined in the IDL.
type ArgosSnifferImpl struct{}

// Time implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) Time(ctx context.Context, req *api.TimeRequest) (resp *api.TimeResponse, err error) {
	// TODO: Your code here...
	return
}

// Transaction implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) Transaction(ctx context.Context, req *api.TimeRequest) (resp *api.TransactionResponse, err error) {
	// TODO: Your code here...
	return
}

// Address implements the ArgosSnifferImpl interface.
func (s *ArgosSnifferImpl) Address(ctx context.Context, req *api.AddressRequest) (resp *api.AddressResponse, err error) {
	// TODO: Your code here...
	return
}
