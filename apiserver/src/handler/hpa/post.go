package hpa_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func HPAApplyHandler(c *gin.Context) {
	horizontalPodAutoscaler := apiobjects.HorizontalPodAutoscaler{}
	err := utils.ReadUnmarshal(c.Request.Body, &horizontalPodAutoscaler)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if horizontalPodAutoscaler.Spec.ScaleTargetRef.Namespace == "" {
		horizontalPodAutoscaler.Spec.ScaleTargetRef.Namespace = global.DefaultNamespace
	}
	if horizontalPodAutoscaler.ObjectMeta.Namespace == "" {
		horizontalPodAutoscaler.ObjectMeta.Namespace = global.DefaultNamespace
	}
	horizontalPodAutoscaler.ObjectMeta.CreationTimestamp = time.Now()
	if horizontalPodAutoscaler.ObjectMeta.UID == "" {
		horizontalPodAutoscaler.ObjectMeta.UID = utils.NewUUID()
	}
	url_horizontalPodAutoscaler := horizontalPodAutoscaler.GetObjectPath()
	val, _ := etcd.Get(url_horizontalPodAutoscaler)
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		var hpa apiobjects.HorizontalPodAutoscaler
		json.Unmarshal([]byte(val), &hpa)
		horizontalPodAutoscaler.ObjectMeta.UID = hpa.ObjectMeta.UID
		topicMessage.ActionType = apiobjects.Update
		horizontalPodAutoscalerJson, _ := json.Marshal(horizontalPodAutoscaler)
		topicMessage.Object = string(horizontalPodAutoscalerJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Put(url_horizontalPodAutoscaler, string(horizontalPodAutoscalerJson))
		listwatch.Publish(global.HPARelevantTopic(), string(topicMessageJson))
		c.String(200, "the horizontalPodAutoscaler is updated")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	horizontalPodAutoscalerJson, _ := json.Marshal(horizontalPodAutoscaler)
	topicMessage.Object = string(horizontalPodAutoscalerJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url_horizontalPodAutoscaler, string(horizontalPodAutoscalerJson))
	listwatch.Publish(global.HPARelevantTopic(), string(topicMessageJson))
	c.String(200, "the horizontalPodAutoscaler is created")
	return
}
