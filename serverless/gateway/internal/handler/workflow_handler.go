package serverless_handler

import (
	"encoding/json"
	"io"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	serveless_utils "minik8s/serverless/gateway/utils"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WorkflowHandler(c *gin.Context) {
	name := c.Param("name")
	params, _ := io.ReadAll(c.Request.Body)
	var jsonData map[string]interface{}
	if err := json.Unmarshal(params, &jsonData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	jsonData["WorkflowId"] = utils.NewUUID()
	defer c.Request.Body.Close()
	data, _ := etcd.Get(route.WorkflowPath + "/" + "default" + "/" + name)
	if data == "" {
		c.JSON(404, gin.H{
			"message": "Workflow not found",
		})
		return
	}
	workflow := apiobjects.Workflow{}
	if err := json.Unmarshal([]byte(data), &workflow); err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
		return
	}
	dag := serveless_utils.Workflow2DAG(&workflow)
	res, err := WorkflowTrigger(jsonData, dag)
	if err != nil {
		c.JSON(404, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, res)
}
func TestHandler(c *gin.Context) {
	data, _ := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	jsonData["WorkflowId"] = utils.NewUUID()
	c.JSON(200, jsonData)
}
func WorkflowTrigger(params map[string]interface{}, dag *serveless_utils.DAG) (map[string]interface{}, error) {
	// 1. get the start node
	node := dag.Root
	var err error
	for node != nil {
		if node.Type == serveless_utils.FunctionType {
			funcName := node.Name
			serviceUrl := GetServiceUrl(funcName)
			if serviceUrl == "" {
				return nil, err
			} else {
				params, err = utils.PostWithJsonReturnJson(serviceUrl, params)
				if err != nil {
					return nil, err
				}
				node = node.Next
			}
		} else if node.Type == serveless_utils.BranchType {
			if node.Branches == nil {
				return nil, err
			}
			paramsjson, err := json.Marshal(params)
			if err != nil {
				return nil, err
			}
			next := serveless_utils.ChooseBranch(node.Branches, string(paramsjson))
			if next == nil {
				return nil, err
			}
			node = next
		}
	}
	return params, nil
}
func GetServiceUrl(name string) string {
	serviceName := "function-" + name + "-service"
	data, _ := etcd.Get("/api/service" + "/" + "default" + "/" + serviceName)
	if data == "" {
		return ""
	}
	service := apiobjects.Service{}
	if err := json.Unmarshal([]byte(data), &service); err != nil {
		return ""
	}
	return "http://" + service.Status.ClusterIP + ":8080"
}
