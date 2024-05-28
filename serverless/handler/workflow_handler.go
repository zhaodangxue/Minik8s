package serverless_handler

import (
	"encoding/json"
	"io"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func WorkflowHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "WorkflowHandler",
	})
}
func TestHandler(c *gin.Context) {
	data, _ := io.ReadAll(c.Request.Body)
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		c.String(200, err.Error())
		return
	}
	jsonData["WorkflowId"] = utils.NewUUID()
	c.JSON(200, jsonData)
}
