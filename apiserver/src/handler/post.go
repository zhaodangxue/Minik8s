package handler

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NodePodBindingHandler(c *gin.Context) {
	binding := apiobjects.NodePodBinding{}
	err := utils.ReadUnmarshal(c.Request.Body, &binding)
	action := apiobjects.Create
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	var bindingJson []byte
	bindingJson, err = json.Marshal(binding)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	url_binding := "/api/binding" + "/" + binding.Pod.Namespace + "/" + binding.Pod.Name
	val, _ := etcd.Get(url_binding)
	if val != "" {
		var binding1 apiobjects.NodePodBinding
		json.Unmarshal([]byte(val), &binding1)
		url_binding = "/api/binding" + "/" + binding1.Pod.Namespace + "/" + binding1.Pod.Name
		etcd.Delete(url_binding)
		action = apiobjects.Update
	}
	path := binding.GetBindingPath()
	etcd.Put(path, string(bindingJson))
	fmt.Printf("action: %v", action)
	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(bindingJson),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.BindingTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}
func PodApplyHandler(c *gin.Context) {
	pod := apiobjects.Pod{}
	err := utils.ReadUnmarshal(c.Request.Body, &pod)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if pod.ObjectMeta.Namespace == "" {
		pod.ObjectMeta.Namespace = global.DefaultNamespace
	}
	url_pod := pod.GetObjectPath()
	val, _ := etcd.Get(url_pod)
	if val != "" {
		c.String(http.StatusOK, "pod already exists")
		return
	}
	pod.ObjectMeta.UID = utils.NewUUID()
	pod.CreationTimestamp = time.Now()
	podJson, _ := json.Marshal(pod)
	etcd.Put(url_pod, string(podJson))
	fmt.Printf("receive pod name: %s namespace: %s uuid: %s", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.ObjectMeta.UID)
	listwatch.Publish(global.SchedulerPodUpdateTopic(), string(podJson))
	c.String(http.StatusOK, "ok")
}

func ServiceCreateHandler(c *gin.Context) {
	svc := apiobjects.Service{}
	err := utils.ReadUnmarshal(c.Request.Body, &svc)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if svc.Data.Namespace == ""{
       svc.Data.Namespace = "default";
	}

	url:= svc.GetObjectPath()
	val,_ := etcd.Get(url)
	if val != ""{
		c.String(http.StatusOK,"service/" +svc.Data.Namespace+"/"+svc.Data.Name+ "/already exists")
		return
	}
	//svc.Data.UID = utils.NewUUID();
	svcJson,_ := json.Marshal(svc)
	etcd.Put(url,string(svcJson))
	fmt.Printf("service create: %s\n",string(svcJson))
    listwatch.Publish(global.ServiceUpdateTopic(),string(svcJson))
    c.String(http.StatusOK,"ok")
}

func ServiceUpdateHandler(c *gin.Context) {
	svc := apiobjects.Service{}
	err := utils.ReadUnmarshal(c.Request.Body, &svc)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	url:= svc.GetObjectPath()
	val,_ := etcd.Get(url)
	if val == ""{
		c.String(http.StatusOK,"service/" +svc.Data.Namespace+"/"+svc.Data.Name+ "/not found")
		return
	}
	svcJson,_ := json.Marshal(svc)
	etcd.Put(url,string(svcJson))
	fmt.Printf("service update: %s\n",string(svcJson))
	listwatch.Publish(global.ServiceUpdateTopic(),string(svcJson))
	c.String(http.StatusOK,"ok")
}
