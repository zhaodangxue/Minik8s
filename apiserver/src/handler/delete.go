package handler

import (
	"encoding/json"
	"fmt"

	//"log"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	apiserver_utils "minik8s/apiserver/src/utils"
	"minik8s/global"
	"minik8s/listwatch"
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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


    svc := apiobjects.Service{}
	err := json.Unmarshal([]byte(val), &svc)
	if err != nil {
		log.Error("[Service Delete Handler] error unmarshalling service object: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if svc.Spec.Type == apiobjects.ServiceTypeNodePort {
		err := etcd.Delete("/api/nodeport/service/"+svc.Data.Namespace+"/"+svc.Data.Name)
		if err != nil {
			log.Error("[Service Delete Handler] error deleting node port: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = etcd.Delete_prefix("api/nodeport/endpoint/"+svc.Data.Name)
		if err != nil {
			log.Error("[Service Delete Handler] error deleting node port endpoints: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
}

func ServiceCmdDeleteHandler(c *gin.Context) {
	//svc := apiobjects.Service{}
	//err := utils.ReadUnmarshal(c.Request.Body, &svc)
	namespace := c.Param("namespace")
	name := c.Param("name")
	action := apiobjects.Delete
	val, _ := etcd.Get("/api/service/" + namespace + "/" + name)
	if val == "" {
		c.String(http.StatusBadRequest, "service/"+namespace+"/"+name+"/not found")
		return
	}
	//etcd.Delete("/api/service/" + namespace + "/" + name)
	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     val,
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
		Object:     string(val),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	c.String(http.StatusOK, "delete endpoint namespace:%s name:%s success", namespace, name)
	listwatch.Publish(global.EndpointTopic(), string(topicMessageJson))

	url := "/api/nodeport/endpoint/" + serviceName + "/" + namespace+"/"+name
	val2,_ := etcd.Get(url)
	if val2 == "" {
		return
	}
	etcd.Delete(url)
	topicMessage2 := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(val2),
	}
	topicMessageJson2, _ := json.Marshal(topicMessage2)
	listwatch.Publish(global.EndpointTopic(), string(topicMessageJson2))
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
	url := "/api/persistentvolumeclaim" + "/" + np + "/" + pvcName
	val, _ := etcd.Get_prefix(url)
	if val[0] == "" {
		c.String(200, "pvc not found")
		return
	}
	var pvc apiobjects.PersistentVolumeClaim
	err := json.Unmarshal([]byte(val[0]), &pvc)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if len(pvc.PodBinding) != 0 {
		c.String(200, "pvc is used by pod")
	}
	var msg apiobjects.TopicMessage
	msg.ActionType = apiobjects.Delete
	msg.Object = val[0]
	msgJson, _ := json.Marshal(msg)
	etcd.Delete(url + "/" + pvc.Spec.StorageClassName)
	listwatch.Publish(global.PvcRelevantTopic(), string(msgJson))
	ret := "delete pvcname:" + pvcName + " namespace:" + np + " success"
	c.String(200, ret)
}
func PVDeleteHandler(c *gin.Context) {
	np := c.Param("namespace")
	pvName := c.Param("name")
	url := "/api/persistentvolume" + "/" + np + "/" + pvName
	val, _ := etcd.Get_prefix(url)
	if val[0] == "" {
		c.String(200, "pv not found")
		return
	}
	var pv apiobjects.PersistentVolume
	err := json.Unmarshal([]byte(val[0]), &pv)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if pv.Status == apiobjects.PVBound {
		c.String(200, "pv is used by pvc")
		return
	}
	var msg apiobjects.TopicMessage
	msg.ActionType = apiobjects.Delete
	msg.Object = val[0]
	msgJson, _ := json.Marshal(msg)
	if pv.Dynamic == 1 {
		err = apiserver_utils.DeletePVPath(pv.ObjectMeta.Name)
		if err != nil {
			c.String(200, err.Error())
			return
		}
	}
	etcd.Delete(url + "/" + pv.Spec.StorageClassName)
	listwatch.Publish(global.PvRelevantTopic(), string(msgJson))
	ret := "delete pvname:" + pvName + " namespace:" + np + " success"
	c.String(200, ret)
}

func DnsDeleteHandler(c *gin.Context) { 
	namespace := c.Param("namespace")
	name := c.Param("name")

	if namespace == "" || name == ""{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "namespace / name is empty",
		})
		return
	}
	url := "/api/dns/" + namespace + "/" + name

	val, err := etcd.Get(url)
	if val == "" || err != nil{
		log.Error("[Dns Delete Handler] dns record not found")
		c.String(http.StatusBadRequest, "dns/"+namespace+"/"+name+"/not found")
		return
	}

	//delete the DNSRecord and the path in the etcd
	err = etcd.Delete(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	dnsRecord := apiobjects.DNSRecord{}
	err = json.Unmarshal([]byte(val), &dnsRecord)
	if err != nil {
		log.Error("[Dns Delete Handler] error unmarshalling dns object: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := generateHost(dnsRecord.Host)
	val2, err := etcd.Get(result)
	if val2 == "" || err != nil{
		log.Error("[Dns Delete Handler] dns record not found2")
		c.String(http.StatusBadRequest, "dns/"+namespace+"/"+name+"/not found2")
		return
	}
	err = etcd.Delete(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// update the nginx config
	err = updateNginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})

}
