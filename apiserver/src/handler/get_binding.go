package handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NodePodBindingAllHandler(c *gin.Context) {
	var bindings []*apiobjects.NodePodBinding
	values, err := etcd.Get_prefix(route.NodePath)
	if err != nil {
		utils.Error("NodePodBindingAllHandler: Get_prefix failed: ", err)
	}
	for _, value := range values {
		var binding apiobjects.NodePodBinding
		err := json.Unmarshal([]byte(value), &binding)
		if err != nil {
			utils.Error("NodePodBindingAllHandler: json.Unmarshal failed: ", err)
		}
		bindings = append(bindings, &binding)
	}
	c.JSON(http.StatusOK, bindings)
}
