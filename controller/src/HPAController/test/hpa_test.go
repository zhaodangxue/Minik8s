package hpa_test

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	hpacontroller "minik8s/controller/src/HPAController"
	replicasetcontroller "minik8s/controller/src/ReplicasetController"
	"minik8s/utils"
	"testing"
	"time"

	ctlutils "minik8s/kubectl/utils"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestHPAAdd(t *testing.T) {
	t.Log("TestHPAAdd")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./replicaset-example.yaml")
	assert.Nil(t, err)
	var replicaset apiobjects.Replicaset
	err = yaml.Unmarshal(content, &replicaset)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller_rs := replicasetcontroller.NewController()
	go controller_rs.Run()
	controller_hpa := hpacontroller.NewController()
	go controller_hpa.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.ReplicasetPath, replicaset)
	assert.Nil(t, err)
	time.Sleep(90 * time.Second)
	content, err = ctlutils.LoadFile("./hpa-example.yaml")
	assert.Nil(t, err)
	var hpa apiobjects.HorizontalPodAutoscaler
	err = yaml.Unmarshal(content, &hpa)
	assert.Nil(t, err)
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.HorizontalPodAutoscalerPath, hpa)
	assert.Nil(t, err)
	time.Sleep(60 * time.Second)
}
func TestHPADelete(t *testing.T) {
	t.Log("TestHPADelete")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./replicaset-example.yaml")
	assert.Nil(t, err)
	var replicaset apiobjects.Replicaset
	err = yaml.Unmarshal(content, &replicaset)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller_rs := replicasetcontroller.NewController()
	go controller_rs.Run()
	controller_hpa := hpacontroller.NewController()
	go controller_hpa.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.ReplicasetPath, replicaset)
	assert.Nil(t, err)
	time.Sleep(90 * time.Second)
	content, err = ctlutils.LoadFile("./hpa-example.yaml")
	assert.Nil(t, err)
	var hpa apiobjects.HorizontalPodAutoscaler
	err = yaml.Unmarshal(content, &hpa)
	assert.Nil(t, err)
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.HorizontalPodAutoscalerPath, hpa)
	assert.Nil(t, err)
	time.Sleep(55 * time.Second)
	url_delete := route.Prefix + hpa.GetObjectPath()
	_, err = utils.Delete(url_delete)
	assert.Nil(t, err)
	time.Sleep(25 * time.Second)
}
func TestHPAUpdate(t *testing.T) {
	t.Log("TestHPAAdd")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./replicaset-example.yaml")
	assert.Nil(t, err)
	var replicaset apiobjects.Replicaset
	err = yaml.Unmarshal(content, &replicaset)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller_rs := replicasetcontroller.NewController()
	go controller_rs.Run()
	controller_hpa := hpacontroller.NewController()
	go controller_hpa.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.ReplicasetPath, replicaset)
	assert.Nil(t, err)
	time.Sleep(90 * time.Second)
	content, err = ctlutils.LoadFile("./hpa-example.yaml")
	assert.Nil(t, err)
	var hpa apiobjects.HorizontalPodAutoscaler
	err = yaml.Unmarshal(content, &hpa)
	assert.Nil(t, err)
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.HorizontalPodAutoscalerPath, hpa)
	assert.Nil(t, err)
	time.Sleep(65 * time.Second)
	hpa.Spec.Metrics.CPUUtilizationPercentage = 50
	_, err = utils.PostWithJson(route.Prefix+route.HorizontalPodAutoscalerPath, hpa)
	assert.Nil(t, err)
	time.Sleep(65 * time.Second)
}
func TestLink(t *testing.T) {
	t.Log("TestLink")
	etcd.Clear()
	content, err := ctlutils.LoadFile("./hpa-example.yaml")
	assert.Nil(t, err)
	var hpa apiobjects.HorizontalPodAutoscaler
	err = yaml.Unmarshal(content, &hpa)
	assert.Nil(t, err)
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller_hpa := hpacontroller.NewController()
	go controller_hpa.Run()
	assert.Nil(t, err)
	_, err = utils.PostWithJson(route.Prefix+route.HorizontalPodAutoscalerPath, hpa)
	assert.Nil(t, err)
	time.Sleep(20 * time.Second)
}
