package kubectl__test

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	command "minik8s/kubectl/src"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var apiServer apiserver.ApiServer

func TestRunApply(t *testing.T) {
	fmt.Println("TestRunApply")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_test("./test.yaml")
	assert.Nil(t, err)
}
func TestRunGet(t *testing.T) {
	fmt.Println("TestRunGet")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_test("test", "111")
	assert.Nil(t, err)
}
func TestRunApplyPod(t *testing.T) {
	fmt.Println("TestRunApplyPod")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_test("./pod-example.yaml")
	assert.Nil(t, err)
}
func TestRunGetPod(t *testing.T) {
	fmt.Println("TestRunGetPod")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_test("pod", "pod")
	assert.Nil(t, err)
}
func TestRunDelete(t *testing.T) {
	fmt.Println("TestRunDelete")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunDelete_test("pod", "pod")
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
	err = command.RunDescribe_test("pod", "pod")
	assert.Nil(t, err)
}
