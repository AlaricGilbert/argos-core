package handlers

import (
	"github.com/AlaricGilbert/argos-core/master/metrics"
	"github.com/gin-gonic/gin"
)

func GetReportStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": metrics.ReportMetrics.GetMetrics(),
	})
}
