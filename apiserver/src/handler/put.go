package handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)


func PodPutHandler(c *gin.Context){

	pod := apiobjects.Pod{}
	err := utils.ReadUnmarshal(c.Request.Body, &pod)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = global.DefaultNamespace
	}
	url_pod := pod.GetObjectPath()
	val, _ := etcd.Get(url_pod)
	if val == "" {
		c.String(http.StatusBadRequest, "pod not exists")
		return
	}
	podJson, _ := json.Marshal(pod)
	etcd.Put(url_pod, string(podJson))
	utils.Info("PodPutHandler: receive pod: ", pod)

	// FIXME: publish correctly
	topicMessage := apiobjects.TopicMessage{ActionType: apiobjects.Update, Object: string(podJson)}
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.SchedulerPodUpdateTopic(), string(topicMessageJson))

	c.String(http.StatusOK, "ok")
}
