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

	"github.com/gin-gonic/gin"
)

func NodePodBindingHandler(c *gin.Context) {
	binding := apiobjects.NodePodBinding{}
	err := utils.ReadUnmarshal(c.Request.Body, &binding)
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
		url_binding = "/api/binding" + "/" + binding1.Pod.Namespace + "/" + binding1.Pod.Name + "/" + binding1.Node.ObjectMeta.Name
		etcd.Delete(url_binding)
	}
	path := binding.GetBindingPath()
	etcd.Put(path, string(bindingJson))
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
