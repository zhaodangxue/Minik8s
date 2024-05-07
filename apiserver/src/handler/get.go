package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NodeGetHandler(c *gin.Context) {
	var nodes []*apiobjects.Node
	values, err := etcd.Get_prefix(route.NodePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var node apiobjects.Node
		err := json.Unmarshal([]byte(value), &node)
		if err != nil {
			fmt.Println(err)
		}
		nodes = append(nodes, &node)
	}
	c.JSON(http.StatusOK, nodes)
}
