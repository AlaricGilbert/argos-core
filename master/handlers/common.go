package handlers

import "github.com/gin-gonic/gin"

func retOK(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 200,
	})
}

func retData(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code": 200,
		"data": data,
	})
}

func retErr(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"code": -1,
		"err":  err.Error(),
	})
}

func retErrMsg(c *gin.Context, err string) {
	c.JSON(200, gin.H{
		"code": -1,
		"err":  err,
	})
}

// retUnwarpErr performs as retErr if err is not nil and performs as retOK in other situations.
func retUnwarpErr(c *gin.Context, err error) {
	if err == nil {
		retOK(c)
	} else {
		retErr(c, err)
	}
}
