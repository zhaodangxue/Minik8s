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
	pod.ObjectMeta.UID = utils.NewUUID()
	pod.CreationTimestamp = time.Now()
	for _, volume := range pod.Spec.Volumes {
		if volume.NFS != nil {
			volume.NFS.BindingPath = "/home/kubelet/volumes/" + utils.NewUUID()
		}
		if volume.PersistentVolumeClaim != nil {
			if volume.PersistentVolumeClaim.ClaimNamespace == "" {
				volume.PersistentVolumeClaim.ClaimNamespace = global.DefaultNamespace
			}
		}
	}
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		topicMessage.ActionType = apiobjects.Update
		podJson, _ := json.Marshal(pod)
		topicMessage.Object = string(podJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Put(url_pod, string(podJson))
		listwatch.Publish(global.PodRelevantTopic(), string(topicMessageJson))
		c.String(http.StatusOK, "pod has configed")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	podJson, _ := json.Marshal(pod)
	topicMessage.Object = string(podJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url_pod, string(podJson))
	fmt.Printf("receive pod name: %s namespace: %s uuid: %s", pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.ObjectMeta.UID)
	listwatch.Publish(global.PodRelevantTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}
func PVApplyHandler(c *gin.Context) {
	pv := apiobjects.PersistentVolume{}
	err := utils.ReadUnmarshal(c.Request.Body, &pv)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if pv.ObjectMeta.Namespace == "" {
		pv.ObjectMeta.Namespace = global.DefaultNamespace
	}
	pv.ObjectMeta.UID = utils.NewUUID()
	pv.CreationTimestamp = time.Now()
	pv.Dynamic = 0
	if pv.Spec.StorageClassName == "" {
		pv.Spec.StorageClassName = "default"
	}
	url_pv := pv.GetObjectPath() + "/" + pv.Spec.StorageClassName
	val, _ := etcd.Get(url_pv)
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		// topicMessage.ActionType = apiobjects.Update
		// pvJson, _ := json.Marshal(pv)
		// topicMessage.Object = string(pvJson)
		// topicMessageJson, _ := json.Marshal(topicMessage)
		// etcd.Delete_prefix(pv.GetObjectPath())
		// etcd.Put(url_pv, string(pvJson))
		// listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
		c.String(http.StatusOK, "pv has already exist,please delete it first")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	pvJson, _ := json.Marshal(pv)
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url_pv, string(pvJson))
	fmt.Printf("receive pv name: %s namespace: %s uuid: %s", pv.ObjectMeta.Name, pv.ObjectMeta.Namespace, pv.ObjectMeta.UID)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}
func PVCApplyHandler(c *gin.Context) {
	pvc := apiobjects.PersistentVolumeClaim{}
	err := utils.ReadUnmarshal(c.Request.Body, &pvc)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if pvc.ObjectMeta.Namespace == "" {
		pvc.ObjectMeta.Namespace = global.DefaultNamespace
	}
	pvc.ObjectMeta.UID = utils.NewUUID()
	pvc.CreationTimestamp = time.Now()
	if pvc.Spec.StorageClassName == "" {
		pvc.Spec.StorageClassName = "default"
	}
	url_pvc := pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
	val, _ := etcd.Get(url_pvc)
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		// topicMessage.ActionType = apiobjects.Update
		// pvcJson, _ := json.Marshal(pvc)
		// topicMessage.Object = string(pvcJson)
		// topicMessageJson, _ := json.Marshal(topicMessage)
		// etcd.Delete_prefix(pvc.GetObjectPath())
		// etcd.Put(url_pvc, string(pvcJson))
		// listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
		c.String(http.StatusOK, "the pvc has already exist,please delete it first")
		return
	}
	topicMessage.ActionType = apiobjects.Create
	pvcJson, _ := json.Marshal(pvc)
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url_pvc, string(pvcJson))
	fmt.Printf("receive pvc name: %s namespace: %s uuid: %s", pvc.ObjectMeta.Name, pvc.ObjectMeta.Namespace, pvc.ObjectMeta.UID)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}
