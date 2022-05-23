package main

import (
	"context"
	"encoding/hex"
	"net"
	"strings"
	"time"

	"github.com/AlaricGilbert/argos-core/argos"
	"github.com/AlaricGilbert/argos-core/master/dal"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/base"
	"github.com/AlaricGilbert/argos-core/master/kitex_gen/master"
	"github.com/AlaricGilbert/argos-core/master/metrics"
	"github.com/AlaricGilbert/argos-core/master/model"
)

// ArgosMasterImpl implements the last service interface defined in the IDL.
type ArgosMasterImpl struct{}

// Ping implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Ping(ctx context.Context, req *master.PingRequest) (resp *master.PingResponse, err error) {
	tt := time.Now().UnixNano()
	logger := argos.StandardLogger()

	logger.WithField("request", req).Info("ping called")

	badResp := &master.PingResponse{
		Status: &base.ResponseStatus{
			Code:    base.StatusInternalError,
			Message: "request nil",
		},
	}

	// query task

	if req == nil {
		logger.Warn("ping exited since request is nil")
		return badResp, nil
	}

	id := req.GetIdentifier()
	protocol := ""
	if pref, _, ok := strings.Cut(id, "-"); ok {
		if task, err := dal.GetTask(pref); err == nil {
			protocol = task.Protocol
		} else {
			logger.WithError(err).Info("query task failed")
		}
	}

	resp = &master.PingResponse{
		Status: &base.ResponseStatus{
			Code:    base.StatusOK,
			Message: "",
		},
		Protocol: protocol,
		TimeSync: &master.TimeSync{
			SendTimestamp: req.GetTimestamp(),
			RecvTimestamp: tt,
			RespTimestamp: time.Now().Unix(),
		},
	}
	logger.WithField("resp", resp).Info("ponged")
	return
}

// Report implements the ArgosMasterImpl interface.
func (s *ArgosMasterImpl) Report(ctx context.Context, req *master.ReportRequest) (resp *master.ReportResponse, err error) {
	logger := argos.StandardLogger()
	logger.WithField("report", req).Info("received report")

	metrics.ReportMetrics.Mark(1)

	r := model.Record{
		Txid:      hex.EncodeToString(req.Transaction.Txid),
		Timestamp: req.Transaction.Timestamp,
		SourceIp:  net.IP(req.Transaction.From.Ip).String(),
		Sniffer:   req.Identifier,
		Protocol:  req.Protocol,
		Method:    req.Method,
	}

	if err := dal.CreateRecord(&r); err != nil {
		logger.WithField("record", r).Info("record create failed")
	}

	if err := dal.CreateOrUpdateConclustion(&r); err != nil {
		logger.WithField("conclustion", r).Info("conclustion create or update failed")
	}

	return &master.ReportResponse{
		Status: &base.ResponseStatus{
			Code:    base.StatusOK,
			Message: "",
		},
	}, nil
}
