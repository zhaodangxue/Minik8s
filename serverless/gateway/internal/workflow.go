package internal

import (
	"encoding/json"
	"errors"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	serveless_utils "minik8s/serverless/gateway/utils"
)

func WorkflowExecutor(name string, jsonData map[string]interface{}) (res map[string]interface{},err error) {
	data, _ := etcd.Get(route.WorkflowPath + "/" + "default" + "/" + name)
	if data == "" {
		err = errors.New("WorkflowExecutor: workflow not found")
		return
	}
	workflow := apiobjects.Workflow{}
	if err = json.Unmarshal([]byte(data), &workflow); err != nil {
		return
	}
	dag := serveless_utils.Workflow2DAG(&workflow)
	res, err = WorkflowTrigger(jsonData, dag)
	return
}
