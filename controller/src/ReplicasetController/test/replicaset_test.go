package replicaset_test

import (
	"testing"
	"time"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	replicasetcontroller "minik8s/controller/src/ReplicasetController"
	ctlutils "minik8s/kubectl/utils"
	"minik8s/utils"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestReplicasetAdd(t *testing.T) {
	t.Log("TestReplicasetAdd")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./replicaset-example.yaml")
	assert.Nil(t, err)
	var replicaset apiobjects.Replicaset
	err = yaml.Unmarshal(content, &replicaset)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := replicasetcontroller.NewController()
	go controller.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.ReplicasetPath, replicaset)
	assert.Nil(t, err)
	time.Sleep(95 * time.Second)
}

func TestReplicasetDelete(t *testing.T) {
	t.Log("TestReplicasetDelete")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./replicaset-example.yaml")
	assert.Nil(t, err)
	var replicaset apiobjects.Replicaset
	err = yaml.Unmarshal(content, &replicaset)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := replicasetcontroller.NewController()
	go controller.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.ReplicasetPath, replicaset)
	assert.Nil(t, err)
	time.Sleep(90 * time.Second)
	// 注意需要我们手动在etcd中设置那些加入的pod的状态为running
	url_delete := route.Prefix + replicaset.GetObjectPath()
	_, err = utils.Delete(url_delete)
	assert.Nil(t, err)
	time.Sleep(15 * time.Second)
}

func TestReplicasetUpdate(t *testing.T) {
	t.Log("TestReplicasetUpdate")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./replicaset-example.yaml")
	assert.Nil(t, err)
	var replicaset apiobjects.Replicaset
	err = yaml.Unmarshal(content, &replicaset)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := replicasetcontroller.NewController()
	go controller.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.ReplicasetPath, replicaset)
	assert.Nil(t, err)
	time.Sleep(80 * time.Second)
	// 注意需要我们手动在etcd中设置那些加入的pod的状态为running
	url_update := route.Prefix + route.ReplicasetPath
	replicaset.Spec.Replicas = 2
	replicaset.Spec.Ready = 0
	_, err = utils.PostWithJson(url_update, replicaset)
	assert.Nil(t, err)
	time.Sleep(110 * time.Second)
}
