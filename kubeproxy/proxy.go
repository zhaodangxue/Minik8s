package kubeproxy

/* 主要工作：
1. 监听service资源的创建。创建service
2. 监听service资源的删除。删除service
3. 监听endpoint的创建。设置dest规则。
4. 监听endpoint的删除。删除对应dest规则。
*/

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/kubeproxy/ipvs"
	"minik8s/listwatch"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func Run() {
	ipvs.Init()
	//ipvs.TestConfig()
	var p proxyServiceHandler
	var e proxyEndpointHandler
	listwatch.Watch(global.ServiceTopic(), p.HandleService)
	listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)

}

/* ========== Service Handler ========== */

type proxyServiceHandler struct {
}

func (p proxyServiceHandler) HandleService(msg *redis.Message) {
	topicMessage := &apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		//调用Handleupdate
		svc := &apiobjects.Service{}
		err2 := json.Unmarshal([]byte(topicMessage.Object), svc)
		if err2 != nil {
			fmt.Println(err2)
		}
		svcJson, _ := json.Marshal(svc)
		p.HandleUpdate([]byte(svcJson))
	case apiobjects.Delete:
		//调用HandleDelete
		svc := &apiobjects.Service{}
		err2 := json.Unmarshal([]byte(topicMessage.Object), svc)
		if err2 != nil {
			fmt.Println(err2)
		}
		svcJson, _ := json.Marshal(svc)
		p.HandleDelete([]byte(svcJson))
	default:
		fmt.Println("error")
	}
}
func (e proxyEndpointHandler) HandleEndpoints(msg *redis.Message) {
	topicMessage := &apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		//调用HandleCreate
		edpt := &apiobjects.Endpoint{}
		err := json.Unmarshal([]byte(topicMessage.Object), edpt)
		if err != nil {
			fmt.Println(err)
		}
		edptJson, _ := json.Marshal(edpt)
		e.HandleCreate([]byte(edptJson))
	case apiobjects.Delete:
		//调用ServiceController删除endpoint
		//ss.HandleDelete([]byte(topicMessage.Object))
	case apiobjects.Update:
		//调用ServiceController更新endpoint
		//ss.HandleUpdate([]byte(topicMessage.Object))
	}
}

func (p proxyServiceHandler) HandleCreate(message []byte) {
}

func (p proxyServiceHandler) HandleDelete(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	for _, p := range svc.Spec.Ports {
		key := svc.Status.ClusterIP + ":" + strconv.Itoa(int(p.Port))
		ipvs.DeleteService(key)
	}

}

func (p proxyServiceHandler) HandleUpdate(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	for _, p := range svc.Spec.Ports {
		ipvs.AddService(svc.Status.ClusterIP, uint16(p.Port))
	}
}

func (p proxyServiceHandler) GetType() string {
	return "proxyserviceHandler"
}

/* ========== Endpoint Handler ========== */

type proxyEndpointHandler struct {
}

func (e proxyEndpointHandler) HandleCreate(message []byte) {
	edpt := &apiobjects.Endpoint{}
	edpt.UnMarshalJSON(message)

	key := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
	ipvs.AddEndpoint(key, edpt.Spec.DestIP, uint16(edpt.Spec.DestPort))
}

func (e proxyEndpointHandler) HandleDelete(message []byte) {
	edpt := &apiobjects.Endpoint{}
	edpt.UnMarshalJSON(message)

	svcKey := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
	dstKey := edpt.Spec.DestIP + ":" + strconv.Itoa(int(edpt.Spec.DestPort))
	ipvs.DeleteEndpoint(svcKey, dstKey)
}

func (e proxyEndpointHandler) HandleUpdate(message []byte) {

}

func (e proxyEndpointHandler) GetType() string {
	return "proxyendpointHandler"
}