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

func PodStatePutHandler(c *gin.Context) {

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
	listwatch.Publish(global.PodStateTopic(), string(topicMessageJson))

	c.String(http.StatusOK, "ok")
}

func NodeHealthHandler(c *gin.Context) {
	node := apiobjects.Node{}
	err := utils.ReadUnmarshal(c.Request.Body, &node)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	url_node := node.GetObjectPath()
	nodeJson, err := json.Marshal(node)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	etcd.Put(url_node, string(nodeJson))
	utils.Info("NodeHealthHandler: receive node: ", node)
	// CHECK: No publish for node health. Should we publish it?
}
