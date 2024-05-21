package replicaset_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"

	"github.com/gin-gonic/gin"
)

func ReplicasetDeleteHandler(c *gin.Context) {
	np := c.Param("namespace")
	name := c.Param("name")
	url := "/api/replicaset/" + np + "/" + name
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "replicaset %s/%s not found", np, name)
		return
	}
	replicaset := apiobjects.Replicaset{}
	err := json.Unmarshal([]byte(val), &replicaset)
	if err != nil {
		c.String(200, err.Error())
	}
	topicMessage := apiobjects.TopicMessage{}
	topicMessage.ActionType = apiobjects.Delete
	topicMessage.Object = string(val)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Delete(url)
	listwatch.Publish(global.ReplicasetTopic(), string(topicMessageJson))
}
