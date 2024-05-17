package replicaset_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func ReplicasetUpdateHandler(c *gin.Context) {
	name := c.Param("name")
	namespace := c.Param("namespace")
	var replicaset apiobjects.Replicaset
	err := utils.ReadUnmarshal(c.Request.Body, &replicaset)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	url := replicaset.GetObjectPath()
	var replicasetJson []byte
	replicasetJson, err = json.Marshal(replicaset)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	etcd.Put(url, string(replicasetJson))
	c.String(200, "the replicaset %s/%s has been updated", namespace, name)
}
