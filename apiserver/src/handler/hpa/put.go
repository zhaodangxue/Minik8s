package hpa_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func HPAUpdateHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	var hpa apiobjects.HorizontalPodAutoscaler
	err := utils.ReadUnmarshal(c.Request.Body, &hpa)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	url := hpa.GetObjectPath()
	var hpaJson []byte
	hpaJson, err = json.Marshal(hpa)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	etcd.Put(url, string(hpaJson))
	c.String(200, "the hpa %s/%s has been updated", namespace, name)
}
