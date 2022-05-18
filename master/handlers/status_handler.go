package handlers

import "github.com/gin-gonic/gin"

func GetStatus(c *gin.Context) {
	// TODO: add status monitoring
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
