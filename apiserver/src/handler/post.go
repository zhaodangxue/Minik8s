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
