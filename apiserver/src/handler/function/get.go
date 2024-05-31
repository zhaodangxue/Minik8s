package function_handler

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

func FunctionGetHandler(c *gin.Context) {
	var functions []*apiobjects.Function
	values, err := etcd.Get_prefix(route.FunctionPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var function apiobjects.Function
		err := json.Unmarshal([]byte(value), &function)
		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{"error": err.Error()})
		}
		functions = append(functions, &function)
	}
	c.JSON(200, functions)
}
