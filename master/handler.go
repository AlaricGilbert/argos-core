package main

import (
	"context"

	"github.com/AlaricGilbert/argos-core/master/kitex_gen/base"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/master"
)

// ArgosMasterImpl implements the last service interface defined in the IDL.
type ArgosMasterImpl struct{}

// Ping implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Ping(ctx context.Context, req *master.PingRequest) (resp *master.PingResponse, err error) {
	resp = &master.PingResponse{
		Status: &base.ResponseStatus{
			Code:  base.STATUSOK,
			Error: "",
		},
	}
	// TODO: Your code here...
	return
}

// Report implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Report(ctx context.Context, req *master.ReportRequest) (resp *master.ReportResponse, err error) {
	// TODO: Your code here...
	return
}
