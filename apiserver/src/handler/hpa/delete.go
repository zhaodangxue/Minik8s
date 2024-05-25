package hpa_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"

	"github.com/gin-gonic/gin"
)

func HPADeleteHandler(c *gin.Context) {
	np := c.Param("namespace")
	name := c.Param("name")
	url := route.HorizontalPodAutoscalerPath + "/" + np + "/" + name
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "hpa %s/%s not found", np, name)
		return
	}
	hpa := apiobjects.HorizontalPodAutoscaler{}
	err := json.Unmarshal([]byte(val), &hpa)
	if err != nil {
		c.String(200, err.Error())
	}
	topicMessage := apiobjects.TopicMessage{}
	topicMessage.ActionType = apiobjects.Delete
	topicMessage.Object = string(val)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Delete(url)
	listwatch.Publish(global.HPARelevantTopic(), string(topicMessageJson))
}
