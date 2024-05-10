package handler

import (
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
	val, _ := etcd.Get("/api/service/" + namespace + "/" + name)
	if val == "" {
		c.String(http.StatusOK, "service/"+namespace+"/"+name+"/not found")
		return
	}
	etcd.Delete("/api/service/" + namespace + "/" + name)
	c.String(http.StatusOK, "delete service namespace:%s name:%s success", namespace, name)
	listwatch.Publish(global.ServiceDeleteTopic(), namespace+"/"+name)
}
