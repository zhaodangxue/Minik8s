package controller

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	//"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	//"minik8s/utils"

	"github.com/go-redis/redis/v8"
)

func (ss *svcServiceHandler)HandleService(msg *redis.Message){
	topicMessage := &apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		//调用ServiceController分配cluster ip，更新serviceList
		svc := &apiobjects.Service{}
		err2 := json.Unmarshal([]byte(topicMessage.Object), svc)
		if err2 != nil {
			fmt.Println(err2)
		}
		svcJson, _ := json.Marshal(svc)
		ss.HandleCreate([]byte(svcJson))
	case apiobjects.Delete:
		//调用ServiceController删除service
		//ss.HandleDelete([]byte(topicMessage.Object))
	case apiobjects.Update:
		//调用ServiceController更新service
		//ss.HandleUpdate([]byte(topicMessage.Object))
	default:
		fmt.Println("error")
	}
}
func (ss *svcEndpointHandler)HandleEndpoints(msg *redis.Message) {
	topicMessage := &apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		//调用ServiceController增加endpoint
		//ss.HandleCreate([]byte(topicMessage.Object))
	case apiobjects.Delete:
		//调用ServiceController删除endpoint
		//ss.HandleDelete([]byte(topicMessage.Object))
	case apiobjects.Update:
		//调用ServiceController更新endpoint
		//ss.HandleUpdate([]byte(topicMessage.Object))
	}
}


func Run() {
	/* service controller */
	var se svcEndpointHandler
	var ss svcServiceHandler

	listwatch.Watch(global.ServiceCmdTopic(), ss.HandleService)
	listwatch.Watch(global.PodStateTopic(), se.HandleEndpoints)

}