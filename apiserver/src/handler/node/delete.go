package node

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NodeDeleteHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	
	// CHECK: ；这里可能会随着ObjectPath的生成方式改变而出错
	nodePath := "/api/node/" + namespace + "/" + name
	value, err := etcd.Get(nodePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "get node namespace:%s name:%s failed", namespace, name)
		return
	}
	err = etcd.Delete(nodePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "delete node namespace:%s name:%s failed", namespace, name)
		return
	}
	node := apiobjects.Node{}
	err = json.Unmarshal([]byte(value), &node)
	if err != nil {
		c.String(http.StatusInternalServerError, "unmarshal node failed")
		return
	}

	// 删除Node相关的资源
	
	// 删除Binding
	// 获取所有的Binding
	values, err := etcd.Get_prefix("/api/binding")
	if err != nil {
		c.String(http.StatusInternalServerError, "get binding failed")
		return
	}
	for _, value := range values {
		binding := apiobjects.NodePodBinding{}
		err := json.Unmarshal([]byte(value), &binding)
		if err != nil {
			utils.Warn("unmarshal binding failed")
			c.String(http.StatusInternalServerError, "unmarshal binding failed")
			// TODO: fault tolerance
			return
		}
		// 删除与该Node相关的Binding
		if binding.Node.GetObjectPath() == nodePath {
			err := etcd.Delete(binding.GetBindingPath())
			if err != nil {
				utils.Warn("delete binding failed")
				c.String(http.StatusInternalServerError, "delete binding failed")
				// TODO: fault tolerance
				return
			}
		}
		// 删除Pod中的状态
		val, err := etcd.Get(binding.Pod.GetObjectPath())
		if err != nil {
			utils.Warn("get pod failed")
			c.String(http.StatusInternalServerError, "get pod failed")
			return
		}
		valJson := []byte(val)
		pod := apiobjects.Pod{}
		err = json.Unmarshal(valJson, &pod)
		if err != nil {
			utils.Warn("unmarshal pod failed")
			c.String(http.StatusInternalServerError, "unmarshal pod failed")
			return
		}
		cleanUpPod(&pod)
		valJson, err = json.Marshal(pod)
		if err != nil {
			utils.Warn("marshal pod failed")
			c.String(http.StatusInternalServerError, "marshal pod failed")
			return
		}
		etcd.Put(pod.GetObjectPath(), string(valJson))
		// 发布Binding删除事件
		bindingJson, err := json.Marshal(binding)
		if err != nil {
			utils.Warn("marshal binding failed")
			c.String(http.StatusInternalServerError, "marshal binding failed")
			return
		}
		message := apiobjects.TopicMessage{
			ActionType: apiobjects.Delete,
			Object:     string(bindingJson),
		}
		messageJson, err := json.Marshal(message)
		if err != nil {
			utils.Warn("marshal message failed")
			c.String(http.StatusInternalServerError, "marshal message failed")
			return
		}
		listwatch.Publish(global.BindingTopic(), string(messageJson))
	}

	// TODO: 发布Node删除事件
}

func cleanUpPod(pod *apiobjects.Pod) {
	pod.Status.PodIP = ""
	pod.Status.PodPhase = apiobjects.PodPhase_POD_CREATED
	pod.Status.HostIP = ""
	pod.Status.SandboxId = ""
}
