package handler

import (
	"encoding/json"
	"fmt"
	//"strconv"

	//"log"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	pv.Spec.PVPath = "/home/kubelet/volumes/pv-" + pv.ObjectMeta.UID
	if pv.Spec.StorageClassName == "" {
		pv.Spec.StorageClassName = "default"
	}
	url_pv := pv.GetObjectPath() + "/" + pv.Spec.StorageClassName
	val, _ := etcd.Get(url_pv)
	var topicMessage apiobjects.TopicMessage
	if val != "" {
		topicMessage.ActionType = apiobjects.Update
		pvJson, _ := json.Marshal(pv)
		topicMessage.Object = string(pvJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Delete_prefix(pv.GetObjectPath())
		etcd.Put(url_pv, string(pvJson))
		listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
		c.String(http.StatusOK, "pv has configed")
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
		topicMessage.ActionType = apiobjects.Update
		pvcJson, _ := json.Marshal(pvc)
		topicMessage.Object = string(pvcJson)
		topicMessageJson, _ := json.Marshal(topicMessage)
		etcd.Put(url_pvc, string(pvcJson))
		listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
		c.String(http.StatusOK, "pvc has configed")
		return
	}
	var PVCPodBinding apiobjects.PVCPodBinding
	PVCPodBinding.PVCName = pvc.ObjectMeta.Name
	PVCPodBinding.PVCNamespace = pvc.ObjectMeta.Namespace
	url_pvc_binding := PVCPodBinding.GetBindingPath()
	PVCPodBindingJson, _ := json.Marshal(PVCPodBinding)
	etcd.Put(url_pvc_binding, string(PVCPodBindingJson))
	topicMessage.ActionType = apiobjects.Create
	pvcJson, _ := json.Marshal(pvc)
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	etcd.Put(url_pvc, string(pvcJson))
	fmt.Printf("receive pvc name: %s namespace: %s uuid: %s", pvc.ObjectMeta.Name, pvc.ObjectMeta.Namespace, pvc.ObjectMeta.UID)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}

func ServiceCreateHandler(c *gin.Context) {
	svc := apiobjects.Service{}
	action := apiobjects.Create
	err := utils.ReadUnmarshal(c.Request.Body, &svc)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if svc.Data.Namespace == "" {
		svc.Data.Namespace = "default"
	}

	url := svc.GetObjectPath()
	val, _ := etcd.Get(url)
	if val != "" {
		c.String(http.StatusOK, "service/"+svc.Data.Namespace+"/"+svc.Data.Name+"/already exists")
		return
	}
	//svc.Data.UID = utils.NewUUID()
	svc.Status.Phase = "CREATED"
	svcJson, _ := json.Marshal(svc)
	etcd.Put(url, string(svcJson))
	fmt.Printf("service create: %s\n", string(svcJson))

	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(svcJson),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.ServiceTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")

	if svc.Spec.Type == apiobjects.ServiceTypeNodePort {
		etcd.Put("/api/nodeport/"+svc.Data.Namespace+"/"+svc.Data.Name, string(svcJson))

		var pods []*apiobjects.Pod
	    values, err := etcd.Get_prefix(route.PodPath)
	    if err != nil {
		    fmt.Println(err)
	    }
	    for _, value := range values {
		    utils.Info("pod value: ", value)
		    var pod apiobjects.Pod
		    err := json.Unmarshal([]byte(value), &pod)
		    if err != nil {
			     fmt.Println(err)
		    }
		    pods = append(pods, &pod)
	    }
		for _, pod := range pods{
			//筛选符合selector条件的pod
			if pod.Status.PodPhase == apiobjects.PodPhase_POD_RUNNING && IsLabelEqual(svc.Spec.Selector, pod.Labels) {
				createEndpoints(&svc, pod)
			}
		}
	}
}
func createEndpoints(svc *apiobjects.Service, pod *apiobjects.Pod) {
	for _, port := range svc.Spec.Ports {
		dstPort := findDstPort(port.TargetPort, pod.Spec.Containers)
		if dstPort == 1314 {
			log.Fatal("[svc controller] No Match for Target Port!")
			return
		}
		spec := apiobjects.EndpointSpec{
			SvcIP:    "HostIP",
			SvcPort:  port.Port,
			DestIP:   pod.Status.PodIP,
			DestPort: dstPort,
		}
		edpt := &apiobjects.Endpoint{
			ServiceName: svc.Data.Name,
			Spec: spec,
			Data: apiobjects.MetaData{
				Name:      svc.Data.Name + "-" + pod.Name + "-port:" + port.TargetPort,
				Namespace: svc.Data.Namespace,
			},
		}
		//TODO 发送http给apiserver,更新edpt
		edptByte, err := edpt.MarshalJSON()
		if err != nil {
			fmt.Println("error")
		}
		etcd.Put("/api/nodeport/"+edpt.Data.Namespace+"/"+edpt.Data.Name, string(edptByte))
		topicMessage := apiobjects.TopicMessage{
			ActionType: apiobjects.Create,
			Object:     string(edptByte),
		}
		topicMessageJson, _ := json.Marshal(topicMessage)
		listwatch.Publish(global.EndpointTopic(), string(topicMessageJson))
	}
}

func IsLabelEqual(a map[string]string, b map[string]string) bool {
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func findDstPort(targetPort string, containers []apiobjects.Container) int32 {
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.Name == targetPort {
				return p.ContainerPort
			}
		}
	}
	log.Fatal("[svc controller] No Match for Target Port!")
	return 1314
}


func ServiceUpdateHandler(c *gin.Context) {
	svc := apiobjects.Service{}
	action := apiobjects.Update
	err := utils.ReadUnmarshal(c.Request.Body, &svc)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	url := svc.GetObjectPath()
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(http.StatusOK, "service/"+svc.Data.Namespace+"/"+svc.Data.Name+"/not found")
		return
	}
	svcJson, _ := json.Marshal(svc)
	etcd.Put(url, string(svcJson))
	fmt.Printf("service update: %s\n", string(svcJson))

	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(svcJson),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.ServiceTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}


func EndpointCreateHandler(c *gin.Context) {
	endpoint := apiobjects.Endpoint{}
	action := apiobjects.Create
	err := utils.ReadUnmarshal(c.Request.Body, &endpoint)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if endpoint.ServiceName == "" {
		c.String(http.StatusOK, "endpoint service name is empty")
		return
	}
	if endpoint.Data.Namespace == "" {
		endpoint.Data.Namespace = "default"
	}
	url := endpoint.GetObjectPath()
	val, _ := etcd.Get(url)
	if val != "" {
		c.String(http.StatusOK, "endpoint/"+endpoint.Data.Namespace+"/"+endpoint.Data.Name+"/already exists")
		return
	}
	endpointJson, _ := json.Marshal(endpoint)
	etcd.Put(url, string(endpointJson))
	fmt.Printf("endpoint create: %s\n", string(endpointJson))

	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(endpointJson),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.EndpointTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}

func ServiceApplyHandler(c *gin.Context) {
	svc := apiobjects.Service{}
	action := apiobjects.Create
	err := utils.ReadUnmarshal(c.Request.Body, &svc)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if svc.Data.Namespace == "" {
		svc.Data.Namespace = global.DefaultNamespace
	}
	url_svc := svc.GetObjectPath()
	val, _ := etcd.Get(url_svc)
	if val != "" {
		c.String(http.StatusOK, "service already exists")
		return
	}
	svc.Data.UID = utils.NewUUID()
	svc.Status.Phase = "CREATING"
	svcJson, _ := json.Marshal(svc)
	//etcd.Put(url_svc, string(svcJson))
	fmt.Printf("receive service name: %s namespace: %s uuid: %s", svc.Data.Name, svc.Data.Namespace, svc.Data.UID)

	topicMessage := apiobjects.TopicMessage{
		ActionType: action,
		Object:     string(svcJson),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.ServiceCmdTopic(), string(topicMessageJson))
	c.String(http.StatusOK, "ok")
}
