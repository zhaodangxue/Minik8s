package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"github.com/gin-gonic/gin"
)

func NodeGetHandler(c *gin.Context) {
	var nodes []*apiobjects.Node
	values, err := etcd.Get_prefix(route.NodePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var node apiobjects.Node
		err := json.Unmarshal([]byte(value), &node)
		if err != nil {
			fmt.Println(err)
		}
		nodes = append(nodes, &node)
	}
	c.JSON(http.StatusOK, nodes)
}

func PodGetWithNamespaceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	var pods []*apiobjects.Pod
	values, err := etcd.Get_prefix(route.PodPath + "/" + namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var pod apiobjects.Pod
		err := json.Unmarshal([]byte(value), &pod)
		if err != nil {
			fmt.Println(err)
		}
		pods = append(pods, &pod)
	}
	c.JSON(http.StatusOK, pods)
}

func PodGetDetailHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	podName := c.Param("name")
	url := "/api/binding" + "/" + namespace + "/" + podName
	val, _ := etcd.Get(url)
	var binding apiobjects.NodePodBinding
	if val == "" {
		c.JSON(http.StatusOK, binding)
	}
	err := json.Unmarshal([]byte(val), &binding)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, binding)
}

func GetAllPodsHandler(c *gin.Context) {
	var pods []*apiobjects.Pod
	values, err := etcd.Get_prefix(route.PodPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		//utils.Info("pod value: ", value)
		var pod apiobjects.Pod
		err := json.Unmarshal([]byte(value), &pod)
		if err != nil {
			fmt.Println(err)
		}
		pods = append(pods, &pod)
	}
	c.JSON(http.StatusOK, pods)
}

func GetAllServicesHandler(c *gin.Context) {
	var services []*apiobjects.Service
	values, err := etcd.Get_prefix(route.ServiceCreatePath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var service apiobjects.Service
		err := json.Unmarshal([]byte(value), &service)
		if err != nil {
			fmt.Println(err)
		}
		services = append(services, &service)
	}
	c.JSON(http.StatusOK, services)
}

func GetOneServiceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	url := route.ServiceCreatePath + "/" + namespace + "/" + name
	val, _ := etcd.Get(url)
	var service apiobjects.Service
	if val == "" {
		c.JSON(http.StatusOK, service)
	}
	err := json.Unmarshal([]byte(val), &service)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, service)
}

func GetAllEndpointsHandler(c *gin.Context) {
	var endpoints []*apiobjects.Endpoint
	values, err := etcd.Get_prefix(route.EndpointCreaetPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var endpoint apiobjects.Endpoint
		err := json.Unmarshal([]byte(value), &endpoint)
		if err != nil {
			fmt.Println(err)
		}
		endpoints = append(endpoints, &endpoint)
	}
	c.JSON(http.StatusOK, endpoints)
}

func GetOneEndpointHandler(c *gin.Context) {
	serviceName := c.Param("serviceName")
	namespace := c.Param("namespace")
	name := c.Param("name")
	url := route.EndpointCreaetPath + "/" + serviceName + "/" + namespace + "/" + name
	val, _ := etcd.Get(url)
	var endpoint apiobjects.Endpoint
	if val == "" {
		c.JSON(http.StatusOK, endpoint)
	}
	err := json.Unmarshal([]byte(val), &endpoint)
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, endpoint)
}

func PVGetWithNamespaceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	var pvs []*apiobjects.PersistentVolume
	values, err := etcd.Get_prefix(route.PVPath + "/" + namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var pv apiobjects.PersistentVolume
		err := json.Unmarshal([]byte(value), &pv)
		if err != nil {
			fmt.Println(err)
		}
		pvs = append(pvs, &pv)
	}
	c.JSON(http.StatusOK, pvs)
}

func PVCGetWithNamespaceHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	var pvcs []*apiobjects.PersistentVolumeClaim
	values, err := etcd.Get_prefix(route.PVCPath + "/" + namespace)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var pvc apiobjects.PersistentVolumeClaim
		err := json.Unmarshal([]byte(value), &pvc)
		if err != nil {
			fmt.Println(err)
		}
		pvcs = append(pvcs, &pvc)
	}
	c.JSON(http.StatusOK, pvcs)
}

func PodGetHandler(c *gin.Context) {
	var pods []*apiobjects.Pod
	values, err := etcd.Get_prefix(route.PodPath)
	if err != nil {
		fmt.Println(err)
	}
	for _, value := range values {
		var pod apiobjects.Pod
		err := json.Unmarshal([]byte(value), &pod)
		if err != nil {
			fmt.Println(err)
		}
		pods = append(pods, &pod)
	}
	c.JSON(http.StatusOK, pods)
}

func DnsGetAllHandler(c *gin.Context) {
	var dnsRecords []apiobjects.DNSRecord = getAllDNSRecords()
	c.JSON(http.StatusOK, dnsRecords)
}
