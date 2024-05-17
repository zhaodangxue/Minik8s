package replicaset_handler

import (
	"encoding/json"
	"net/http"
	"time"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func ReplicasetApplyHandler(c *gin.Context) {
	replicaset := apiobjects.Replicaset{}
	err := utils.ReadUnmarshal(c.Request.Body, &replicaset)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if replicaset.ObjectMeta.Namespace == "" {
		replicaset.ObjectMeta.Namespace = global.DefaultNamespace
	}
	replicaset.ObjectMeta.CreationTimestamp = time.Now()
	if replicaset.ObjectMeta.UID == "" {
		replicaset.ObjectMeta.UID = utils.NewUUID()
	}
	url_replicaset := replicaset.GetObjectPath()
	val, _ := etcd.Get(url_replicaset)
	replicaset.Spec.Ready = 0
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		var rs apiobjects.Replicaset
		json.Unmarshal([]byte(val), &rs)
		replicaset.ObjectMeta.UID = rs.ObjectMeta.UID
		topicMessage.ActionType = apiobjects.Update
		replicasetJson, _ := json.Marshal(replicaset)
		topicMessage.Object = string(replicasetJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Put(url_replicaset, string(replicasetJson))
		listwatch.Publish(global.ReplicasetTopic(), string(topicMessageJson))
		c.String(http.StatusOK, "the replicaset is updated")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	replicasetJson, _ := json.Marshal(replicaset)
	topicMessage.Object = string(replicasetJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url_replicaset, string(replicasetJson))
	listwatch.Publish(global.ReplicasetTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "the replicaset is created")
}
