package serverless_handler

import "github.com/gin-gonic/gin"

func FunctionHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "FunctionHandler",
	})
}
