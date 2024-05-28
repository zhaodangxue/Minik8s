package serverless_handler

import "github.com/gin-gonic/gin"

func WorkflowHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "WorkflowHandler",
	})
}
