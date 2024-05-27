package hpa_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

func HPAGetWithNamespaceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	var horizontalPodAutoscalers []*apiobjects.HorizontalPodAutoscaler
	values, err := etcd.Get_prefix(route.HorizontalPodAutoscalerPath + "/" + namespace)
	if err != nil {
		c.String(200, err.Error())
	}
	for _, value := range values {
		var horizontalPodAutoscaler apiobjects.HorizontalPodAutoscaler
		err := json.Unmarshal([]byte(value), &horizontalPodAutoscaler)
		if err != nil {
			c.String(200, err.Error())
		}
		horizontalPodAutoscalers = append(horizontalPodAutoscalers, &horizontalPodAutoscaler)
	}
	c.JSON(200, horizontalPodAutoscalers)
}
func HPAGetHandler(c *gin.Context) {
	var horizontalPodAutoscalers []*apiobjects.HorizontalPodAutoscaler
	values, err := etcd.Get_prefix(route.HorizontalPodAutoscalerPath)
	if err != nil {
		c.String(200, err.Error())
	}
	for _, value := range values {
		var horizontalPodAutoscaler apiobjects.HorizontalPodAutoscaler
		err := json.Unmarshal([]byte(value), &horizontalPodAutoscaler)
		if err != nil {
			c.String(200, err.Error())
		}
		horizontalPodAutoscalers = append(horizontalPodAutoscalers, &horizontalPodAutoscaler)
	}
	c.JSON(200, horizontalPodAutoscalers)
}
