package main

import (
	"context"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/api"
)

// ArgosMasterImpl implements the last service interface defined in the IDL.
type ArgosMasterImpl struct{}

// Register implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Register(ctx context.Context, req *api.RegisterRequest) (resp *api.RegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Ping implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Ping(ctx context.Context, req *api.PingRequest) (resp *api.PingResponse, err error) {
	// TODO: Your code here...
	return
}

// Report implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Report(ctx context.Context, req *api.ReportRequest) (resp *api.ReportResponse, err error) {
	// TODO: Your code here...
	return
}
