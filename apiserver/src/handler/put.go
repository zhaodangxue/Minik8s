package handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/api"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NodeHealthHandler(c *gin.Context) {

	healthReport := api.NodeHealthReportRequest{}
	err := utils.ReadUnmarshal(c.Request.Body, &healthReport)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Update node
	nodeJson, err := json.Marshal(healthReport.Node)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	etcd.Put(healthReport.Node.GetObjectPath(), string(nodeJson))

	response := api.NodeHealthReportResponse{}
	response.UnmatchedPodPaths = make([]string, 0)

	// Update pods
	for _, pod := range healthReport.Pods {
		// Check binding
		binding, err := etcd.Get(apiobjects.GetBindingPath(pod))
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		var nodePodBinding apiobjects.NodePodBinding
		err = json.Unmarshal([]byte(binding), &nodePodBinding)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		if healthReport.Node.Equals(&nodePodBinding.Node.Object) == false {
			utils.Warn("NodeHealthHandler: node not match, pod=", pod.ObjectMeta.Name, " node=", nodePodBinding.Node.ObjectMeta.Name)
			response.UnmatchedPodPaths = append(response.UnmatchedPodPaths, pod.GetObjectPath())
			continue
		}

		podJson, err := json.Marshal(pod)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		etcd.Put(pod.GetObjectPath(), string(podJson))
		// Publish pod update event
		message := apiobjects.TopicMessage{}
		message.ActionType = apiobjects.Update
		message.Object = string(podJson)
		message_payload, err := json.Marshal(message)
		if err != nil {
			utils.Error("NodeHealthHandler: json.Marshal failed: ", err)
		}
		listwatch.Publish(global.PodStateTopic(), string(message_payload))
	}

	c.JSON(http.StatusOK, response)
}
