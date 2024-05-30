package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	serveless_utils "minik8s/serverless/gateway/utils"
	"minik8s/utils"
	"net/http"
	"time"

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
	defer c.Request.Body.Close()

	res, err := WorkflowExecutor(name, jsonData)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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
				val := JudgeReplicas(funcName)
				if val != "success" {
					fmt.Println(val)
				}
				AddQpsCounter(funcName)
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
func AddQpsCounter(name string) {
	// TO DO
	if ServerlessGatewayInstance.functions[name] == nil {
		return
	} else {
		currentQps := ServerlessGatewayInstance.functions[name].QPSCounter.Load()
		ServerlessGatewayInstance.functions[name].QPSCounter.Store(currentQps + 1)
	}
}
func JudgeReplicas(name string) string {
	fmt.Println("JudgeReplicas")
	if ServerlessGatewayInstance.functions[name] == nil {
		fmt.Println("ServerlessGatewayInstance.functions[name] is nil")
		return "ServerlessGatewayInstance.functions[name] is nil"
	}
	if ServerlessGatewayInstance.functions[name].ScaleTarget == 0 {
		//这个时候我们必须要去etcd中获取这个function的replicaset并为他扩容
		fmt.Println("ScaleTarget is 0, need to scale replicaset")
		rs, _ := etcd.Get(route.ReplicasetPath + "/" + "default" + "/" + "function-" + name + "-rs")
		if rs == "" {
			return "Replicaset not found with name: " + "function-" + name + "-rs"
		}
		replicaset := apiobjects.Replicaset{}
		if err := json.Unmarshal([]byte(rs), &replicaset); err != nil {
			return err.Error()
		}
		fun, _ := etcd.Get(route.FunctionPath + "/" + "default" + "/" + name)
		if fun == "" {
			return "Function not found with name: " + name
		}
		function := apiobjects.Function{}
		if err := json.Unmarshal([]byte(fun), &function); err != nil {
			return err.Error()
		}
		replicaset.Spec.Replicas = function.Spec.MinReplicas
		url := route.Prefix + route.ReplicasetScale
		_, err := utils.PutWithJson(url, replicaset)
		if err != nil {
			return "Failed to scale replicaset"
		}
		fmt.Println("Scale replicaset success")
		time.Sleep(15 * time.Second)
	}
	return "success"
}
