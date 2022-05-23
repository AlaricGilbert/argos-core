package handlers

import (
	"encoding/hex"
	"net"
	"strconv"

	"github.com/AlaricGilbert/argos-core/master/dal"
	"github.com/gin-gonic/gin"
)

func QueryByTime(c *gin.Context) {
	method := c.Query("method")

	prev := c.Query("prev")
	psStr := c.Query("ps")

	fromStr := c.Query("from")
	toStr := c.Query("to")

	var query dal.ConclusionQuery
	var err error
	if method == "" {
		retErrMsg(c, "method should not be empty")
		return
	}
	query.Method = method
	query.Protocol = c.Query("protocol")

	if prev != "" {
		if prevTx, err := dal.GetSingleConclusion(prev, method); err != nil {
			retErrMsg(c, "prev txid is not valid")
			return
		} else {
			query.Offset = prevTx.ID
		}
	}

	if psStr != "" {
		if query.Limits, err = strconv.Atoi(psStr); err != nil {
			query.Limits = 0
		}
	}

	if fromStr != "" {
		if query.Offset, err = strconv.ParseInt(fromStr, 10, 64); err != nil {
			query.TimeFrom = 0
		}
	}

	if toStr != "" {
		if query.Offset, err = strconv.ParseInt(toStr, 10, 64); err != nil {
			query.TimeTo = 0
		}
	}

	if result, err := dal.GetConclusions(query); err != nil {
		retErr(c, err)
	} else {
		retData(c, result)
	}
}

func QueryByIP(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		retErrMsg(c, "ip should not be empty")
		return
	}

	// try parse ip
	if net.ParseIP(ip) == nil {
		retErrMsg(c, "ip is not valid")
		return
	}

	if result, err := dal.GetRecordsWithIP(ip); err != nil {
		retErr(c, err)
	} else {
		retData(c, result)
	}
}

func QueryByTx(c *gin.Context) {
	txid := c.Query("txid")
	if txid == "" {
		retErrMsg(c, "txid should not be empty")
		return
	}

	if _, err := hex.DecodeString(txid); err != nil {
		retErrMsg(c, "txid is not valid")
		return
	}

	if result, err := dal.GetRecordsWithTxid(txid); err != nil {
		retErr(c, err)
	} else {
		retData(c, result)
	}
}
