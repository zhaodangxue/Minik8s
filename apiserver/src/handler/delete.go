package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServiceDeleteHandler(c *gin.Context) {
	//svc := apiobjects.Service{}
	//err := utils.ReadUnmarshal(c.Request.Body, &svc)
	namespace := c.Param("namespace")
	name := c.Param("name")
	action := apiobjects.Delete
	val, _ := etcd.Get("/api/service/" + namespace + "/" + name)
	if val == "" {
		c.String(http.StatusOK, "service/"+namespace+"/"+name+"/not found")
		return
	}
	etcd.Delete("/api/service/" + namespace + "/" + name)
	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(val),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	c.String(http.StatusOK, "delete service namespace:%s name:%s success", namespace, name)
	listwatch.Publish(global.ServiceTopic(), string(topicMessageJson))
}

func ServiceCmdDeleteHandler(c *gin.Context) {
	//svc := apiobjects.Service{}
	//err := utils.ReadUnmarshal(c.Request.Body, &svc)
	namespace := c.Param("namespace")
	name := c.Param("name")
	action := apiobjects.Delete
	val, _ := etcd.Get("/api/service/" + namespace + "/" + name)
	if val == "" {
		c.String(http.StatusOK, "service/"+namespace+"/"+name+"/not found")
		return
	}
	//etcd.Delete("/api/service/" + namespace + "/" + name)
	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(val),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	c.String(http.StatusOK, "delete service namespace:%s name:%s cmd:%s success", namespace, name)
	listwatch.Publish(global.ServiceCmdTopic(), string(topicMessageJson))
}

func EndpointDeleteHandler(c *gin.Context) {
	//edpt := apiobjects.Endpoint{}
	//err := utils.ReadUnmarshal(c.Request.Body, &edpt)
	serviceName := c.Param("serviceName")
	namespace := c.Param("namespace")
	name := c.Param("name")
	action := apiobjects.Delete
	val, _ := etcd.Get("/api/endpoint/" + serviceName + "/" + namespace + "/" + name)
	if val == "" {
		c.String(http.StatusOK, "endpoint/"+namespace+"/"+name+"/not found")
		return
	}
	etcd.Delete("/api/endpoint/" + serviceName + "/" + namespace + "/" + name)
	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(serviceName + "/" + namespace + "/" + name),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	c.String(http.StatusOK, "delete endpoint namespace:%s name:%s success", namespace, name)
	listwatch.Publish(global.EndpointTopic(), string(topicMessageJson))
}

func PodDeleteHandler(c *gin.Context) {
	np := c.Param("namespace")
	podName := c.Param("name")
	url := "/api/pod" + "/" + np + "/" + podName
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "pod not found")
		return
	}
	var pod apiobjects.Pod
	err := json.Unmarshal([]byte(val), &pod)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	var binding apiobjects.NodePodBinding
	binding.Pod = pod
	val_binding, _ := etcd.Get(binding.GetBindingPath())
	if val_binding == "" {
		fmt.Println("no binding")
	}
	var msg apiobjects.TopicMessage
	msg.ActionType = apiobjects.Delete
	msg.Object = val_binding
	msgJson, _ := json.Marshal(msg)
	var msg2 apiobjects.TopicMessage
	msg2.ActionType = apiobjects.Delete
	msg2.Object = val
	msgJson2, _ := json.Marshal(msg2)
	etcd.Delete(url)
	etcd.Delete(binding.GetBindingPath())
	listwatch.Publish(global.PodRelevantTopic(), string(msgJson2))
	listwatch.Publish(global.BindingTopic(), string(msgJson))
	ret := "delete podname:" + podName + " namespace:" + np + " success"
	c.String(200, ret)
}
func PVCDeleteHandler(c *gin.Context) {
	np := c.Param("namespace")
	pvcName := c.Param("name")
	url := "/api/persistentvolume" + "/" + np + "/" + pvcName
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "pvc not found")
		return
	}
	var PVCPodBinding apiobjects.PVCPodBinding
	PVCPodBinding.PVCName = pvcName
	PVCPodBinding.PVCNamespace = np
	url_pvc_binding := PVCPodBinding.GetBindingPath()
	val_binding, _ := etcd.Get(url_pvc_binding)
	err := json.Unmarshal([]byte(val_binding), &PVCPodBinding)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if len(PVCPodBinding.Pods) != 0 {
		c.String(200, "pvc is used by pod")
		return
	}
	var pvc apiobjects.PersistentVolumeClaim
	err = json.Unmarshal([]byte(val), &pvc)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	var msg apiobjects.TopicMessage
	msg.ActionType = apiobjects.Delete
	msg.Object = val
	msgJson, _ := json.Marshal(msg)
	etcd.Delete(url)
	listwatch.Publish(global.PvcRelevantTopic(), string(msgJson))
	ret := "delete pvcname:" + pvcName + " namespace:" + np + " success"
	c.String(200, ret)
}
func PVDeleteHandler(c *gin.Context) {
	np := c.Param("namespace")
	pvName := c.Param("name")
	url := "/api/persistentvolumeclaim" + "/" + np + "/" + pvName
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "pv not found")
		return
	}
	var pv apiobjects.PersistentVolume
	err := json.Unmarshal([]byte(val), &pv)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if pv.Status == apiobjects.Bound {
		c.String(200, "pv is used by pvc")
		return
	}
	var msg apiobjects.TopicMessage
	msg.ActionType = apiobjects.Delete
	msg.Object = val
	msgJson, _ := json.Marshal(msg)
	etcd.Delete(url)
	listwatch.Publish(global.PvRelevantTopic(), string(msgJson))
	ret := "delete pvname:" + pvName + " namespace:" + np + " success"
	c.String(200, ret)
}
