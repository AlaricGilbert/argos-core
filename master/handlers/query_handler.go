package handlers

import (
	"strconv"

	"github.com/AlaricGilbert/argos-core/master/dal"
	"github.com/gin-gonic/gin"
)

func QueryByTime(c *gin.Context) {
	method := c.Query("method")

	offsetStr := c.Query("offset")
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

	if offsetStr != "" {
		if query.Offset, err = strconv.ParseInt(offsetStr, 10, 64); err != nil {
			query.Offset = 0
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
	if result, err := dal.GetRecordsWithIP(txid); err != nil {
		retErr(c, err)
	} else {
		retData(c, result)
	}
}
