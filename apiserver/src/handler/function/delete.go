package function_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"

	"github.com/gin-gonic/gin"
)

func FunctionDeleteHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	url := route.FunctionPath + "/" + namespace + "/" + name
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "the function is not found")
		return
	}
	etcd.Delete(url)
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Delete
	topicMessage.Object = val
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Delete(url)
	listwatch.Publish(global.FunctionTopic(), string(topicMessageJson))
	c.String(200, "the function is deleted")
}
