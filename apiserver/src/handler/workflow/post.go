package workflow_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func WorkflowApplyHandler(c *gin.Context) {
	workflow := apiobjects.Workflow{}
	err := utils.ReadUnmarshal(c.Request.Body, &workflow)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if workflow.ObjectMeta.Namespace == "" {
		workflow.ObjectMeta.Namespace = global.DefaultNamespace
	}
	workflow.ObjectMeta.CreationTimestamp = time.Now()
	if workflow.ObjectMeta.UID == "" {
		workflow.ObjectMeta.UID = utils.NewUUID()
	}
	url := route.WorkflowPath + "/" + workflow.ObjectMeta.Namespace + "/" + workflow.ObjectMeta.Name
	val, _ := etcd.Get(url)
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		var oldWorkflow apiobjects.Workflow
		json.Unmarshal([]byte(val), &oldWorkflow)
		oldWorkflow.ObjectMeta.UID = workflow.ObjectMeta.UID
		topicMessage.ActionType = apiobjects.Update
		oldWorkflowJson, _ := json.Marshal(oldWorkflow)
		topicMessage.Object = string(oldWorkflowJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Put(url, string(oldWorkflowJson))
		listwatch.Publish(global.WorkFlowTopic(), string(topicMessageJson))
		c.String(200, "the workflow is updated")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	workflowJson, _ := json.Marshal(workflow)
	topicMessage.Object = string(workflowJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url, string(workflowJson))
	listwatch.Publish(global.WorkFlowTopic(), string(topicMessageJson))
	c.String(200, "the workflow is created")
	return
}
