package test

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	ctlutils "minik8s/kubectl/utils"
	"minik8s/serverless/gateway/internal"
	serveless_utils "minik8s/serverless/gateway/utils"
	"minik8s/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddNumber(t *testing.T) {
	t.Log("TestAddNumber")
	gateway := internal.GetServerlessGatewayInstance()
	go gateway.RUN()
	time.Sleep(3 * time.Second)
	var jsonData map[string]interface{}
	jsonData = make(map[string]interface{})
	jsonData["a"] = 1
	jsonData["b"] = 2
	_, err := utils.PostWithJson("http://localhost:8081/serverless/test", jsonData)
	assert.Nil(t, err)
}
func TestBuildDAG(t *testing.T) {
	t.Log("TestBuildDAG")
	var data []byte
	var err error
	data, err = ctlutils.LoadFile("./workflow.json")
	assert.Nil(t, err)
	workflow := apiobjects.Workflow{}
	if err = json.Unmarshal(data, &workflow); err != nil {
		fmt.Println(err)
		return
	}
	dag := serveless_utils.Workflow2DAG(&workflow)
	assert.NotNil(t, dag)
}
func TestCompareFunction(t *testing.T) {
	t.Log("TestCompareFunction")
	var a int64
	a = 5
	var jsonMap map[string]interface{}
	jsonMap = make(map[string]interface{})
	jsonMap["a"] = 5
	jsonMap["b"] = 5
	jsonMapJson, _ := json.Marshal(jsonMap)
	flag := serveless_utils.IntegerEqual(string(jsonMapJson), "a", a)
	assert.True(t, flag)
}
func TestChooseBranch(t *testing.T) {
	t.Log("TestChooseBranch")
	var branches []*serveless_utils.Branch
	var branch1 serveless_utils.Branch
	var branch2 serveless_utils.Branch
	branch1.Next = nil
	branch1.Value = int64(5)
	branch1.BranchFunc = serveless_utils.IntegerEqual
	branch1.Variable = "a"
	branch2.Next = nil
	branch2.Value = int64(5)
	branch2.BranchFunc = serveless_utils.IntegerNotEqual
	branch2.Variable = "a"
	branches = append(branches, &branch1)
	branches = append(branches, &branch2)
	var jsonMap1 map[string]interface{}
	jsonMap1 = make(map[string]interface{})
	jsonMap1["a"] = 5
	jsonMap1["b"] = 5
	jsonMapJson1, _ := json.Marshal(jsonMap1)
	var jsonMap2 map[string]interface{}
	jsonMap2 = make(map[string]interface{})
	jsonMap2["a"] = 6
	jsonMap2["b"] = 5
	jsonMapJson2, _ := json.Marshal(jsonMap2)
	node1 := serveless_utils.ChooseBranch(branches, string(jsonMapJson1))
	node2 := serveless_utils.ChooseBranch(branches, string(jsonMapJson2))
	assert.Nil(t, node1)
	assert.Nil(t, node2)
}
