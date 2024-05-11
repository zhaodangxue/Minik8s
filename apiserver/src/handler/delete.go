package handler

import (
	"encoding/json"
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
		Object:     string(namespace + "/" + name),
	}
	topicMessageJson, _ := json.Marshal(topicMessage)
	c.String(http.StatusOK, "delete service namespace:%s name:%s success", namespace, name)
	listwatch.Publish(global.ServiceTopic(), string(topicMessageJson))
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
