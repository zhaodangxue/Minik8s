package kubectl__test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	command "minik8s/kubectl/src"
	ctlutils "minik8s/kubectl/utils"

	"github.com/stretchr/testify/assert"
)

var apiServer apiserver.ApiServer

func TestRunApply(t *testing.T) {
	fmt.Println("TestRunApply")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_Cmd("./test.yaml")
	assert.Nil(t, err)
}

func TestRunGet(t *testing.T) {
	fmt.Println("TestRunGet")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_Cmd("test", "111")
	assert.Nil(t, err)
}

func TestRunApplyPod(t *testing.T) {
	fmt.Println("TestRunApplyPod")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_Cmd("./pod-example.yaml")
	assert.Nil(t, err)
}

func TestRunGetPod(t *testing.T) {
	fmt.Println("TestRunGetPod")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_Cmd("pod", "pod")
	assert.Nil(t, err)
}

func TestRunDelete(t *testing.T) {
	fmt.Println("TestRunDelete")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunDelete_Cmd("pod", "pod")
	assert.Nil(t, err)
}

func TestRunDescribe(t *testing.T) {
	fmt.Println("TestRunDescribe")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	etcd.Clear()
	var pod apiobjects.Pod
	pod.ObjectMeta.Name = "pod"
	pod.ObjectMeta.Namespace = "default"
	pod.Spec.Containers = []apiobjects.Container{
		{
			Name:  "nginx",
			Image: "nginx",
		},
	}
	var node apiobjects.Node
	node.ObjectMeta.Name = "node"
	node.Info.Ip = "10.10.10.1"
	binding := apiobjects.NodePodBinding{
		Node: node,
		Pod:  pod,
	}
	url := binding.GetBindingPath()
	bindingJson, err := json.Marshal(binding)
	assert.Nil(t, err)
	etcd.Put(url, string(bindingJson))
	err = command.RunDescribe_Cmd("pod", "pod")
	assert.Nil(t, err)
}

func TestRunApplyPV(t *testing.T) {
	fmt.Println("TestRunApplyPV")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_Cmd("./persistent_volumn.yaml")
	assert.Nil(t, err)
}

func TestRunApplyPVC(t *testing.T) {
	fmt.Println("TestRunApplyPVC")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_Cmd("./persistent_volumn_claim.yaml")
	assert.Nil(t, err)
}

func TestRunGetPV(t *testing.T) {
	fmt.Println("TestRunGetPV")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_Cmd("pv", "nfsserver")
	assert.Nil(t, err)
}

func TestRunGetPVC(t *testing.T) {
	fmt.Println("TestRunGetPVC")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_Cmd("pvc", "pvc-sc-example")
	assert.Nil(t, err)
}

func TestApplyReplicaset(t *testing.T) {
	fmt.Println("TestApplyReplicaset")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_Cmd("./replicaset-example.yaml")
	assert.Nil(t, err)
}
func TestApplyHPA(t *testing.T) {
	fmt.Println("TestApplyHPA")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_Cmd("./hpa-example.yaml")
	assert.Nil(t, err)
}
func TestGetHPA(t *testing.T) {
	fmt.Println("TestGetHPA")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_Cmd("hpa", "hpa")
	assert.Nil(t, err)
}
func TestApplyWorkflow(t *testing.T) {
	fmt.Println("TestApplyWorkflow")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	var data []byte
	var err error
	data, err = ctlutils.LoadFile("./workflow.json")
	assert.Nil(t, err)
	workflow := apiobjects.Workflow{}
	if err = json.Unmarshal(data, &workflow); err != nil {
		fmt.Println(err)
		return
	}
	err = command.AddWorkflowToApiServer(workflow)
	assert.Nil(t, err)
}
func TestDeleteWorkflow(t *testing.T) {
	fmt.Println("TestDeleteWorkflow")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.DeleteWorkflowFromApiServer("example")
	assert.Nil(t, err)
}
