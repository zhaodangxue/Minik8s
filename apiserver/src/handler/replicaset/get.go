package replicaset_handler

import (
	"encoding/json"
	"fmt"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

func ReplicasetGetWithNamespaceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	var replicasets []*apiobjects.Replicaset
	values, err := etcd.Get_prefix(route.ReplicasetPath + "/" + namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var replicaset apiobjects.Replicaset
		err := json.Unmarshal([]byte(value), &replicaset)
		if err != nil {
			fmt.Println(err)
		}
		replicasets = append(replicasets, &replicaset)
	}
	c.JSON(200, replicasets)
}

func ReplicasetGetHandler(c *gin.Context) {
	var replicasets []*apiobjects.Replicaset
	values, err := etcd.Get_prefix(route.ReplicasetPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var replicaset apiobjects.Replicaset
		err := json.Unmarshal([]byte(value), &replicaset)
		if err != nil {
			fmt.Println(err)
		}
		replicasets = append(replicasets, &replicaset)
	}
	c.JSON(200, replicasets)
}
