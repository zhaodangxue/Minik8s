package serverless__test

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
	gateway := internal.New()
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
