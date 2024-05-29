package event

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func EventCreateHandler(c *gin.Context) {
	var event apiobjects.Event
	err := utils.ReadUnmarshal(c.Request.Body, &event)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if event.Name == "" {
		c.String(200, "Please input the event name")
		return
	}
	url := route.EventPath + "/" + event.Name
	val, _ := etcd.Get(url)
	if val != "" {
		c.String(200, "The event " + event.Name + " is already exist")
		return
	}
	eventJson, _ := json.Marshal(event)
	err = etcd.Put(url, string(eventJson))
	if err != nil {
		c.String(500, "Create event failed")
		return
	}
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(eventJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.EventTopic(), string(topicMessageJson))
	c.String(200, "the event is created")
}

func EventDeleteHandler(c *gin.Context) {
	name := c.Param("name")
	url := route.EventPath + "/" + name
	val, _ := etcd.Get(url)
	if val == "" {
		c.String(200, "The event " + name + " is not found")
		return
	}
	err := etcd.Delete(url)
	if err != nil {
		c.String(500, "Delete event failed")
		return
	}
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Delete
	topicMessage.Object = val
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.EventTopic(), string(topicMessageJson))
	c.String(200, "the function is deleted")
}
