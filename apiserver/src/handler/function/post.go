package function_handler

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

func FunctionApplyHandler(c *gin.Context) {
	function := apiobjects.Function{}
	err := utils.ReadUnmarshal(c.Request.Body, &function)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if function.ObjectMeta.Namespace == "" {
		function.ObjectMeta.Namespace = global.DefaultNamespace
	}
	function.ObjectMeta.CreationTimestamp = time.Now()
	if function.ObjectMeta.UID == "" {
		function.ObjectMeta.UID = utils.NewUUID()
	}
	url := route.FunctionPath + "/" + function.ObjectMeta.Namespace + "/" + function.ObjectMeta.Name
	val, _ := etcd.Get(url)
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		var oldFunction apiobjects.Function
		json.Unmarshal([]byte(val), &oldFunction)
		function.ObjectMeta.UID = oldFunction.ObjectMeta.UID
		topicMessage.ActionType = apiobjects.Update
		FunctionJson, _ := json.Marshal(function)
		topicMessage.Object = string(FunctionJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Put(url, string(FunctionJson))
		listwatch.Publish(global.FunctionTopic(), string(topicMessageJson))
		c.String(200, "the function is updated")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	functionJson, _ := json.Marshal(function)
	topicMessage.Object = string(functionJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url, string(functionJson))
	listwatch.Publish(global.FunctionTopic(), string(topicMessageJson))
	c.String(200, "the function is created")
}
