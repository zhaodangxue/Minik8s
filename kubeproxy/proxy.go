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
	"minik8s/apiserver/src/route"
	"minik8s/utils"

	//"minik8s/utils"

	//"minik8s/global"
	"minik8s/kubeproxy/ipvs"
	//"minik8s/listwatch"
	"strconv"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

// func main() {
// 	ipvs.Init()
// 	//ipvs.TestConfig()
// 	var p proxyServiceHandler
// 	var e proxyEndpointHandler
// 	listwatch.Watch(global.ServiceTopic(), p.HandleService)
// 	go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
// }

/* ========== Service Handler ========== */

type ProxyServiceHandler struct {
}

func (p ProxyServiceHandler) HandleService(msg *redis.Message) {
	utils.Info("Proxy Handle Service")
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
func (e ProxyEndpointHandler) HandleEndpoints(msg *redis.Message) {
	utils.Info("Proxy Handle Endpoints")
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
		edpt := &apiobjects.Endpoint{}
		err := json.Unmarshal([]byte(topicMessage.Object), edpt)
		if err != nil {
			fmt.Println(err)
		}
		edptJson, _ := json.Marshal(edpt)
		e.HandleDelete([]byte(edptJson))
	case apiobjects.Update:
		//调用ServiceController更新endpoint
		//ss.HandleUpdate([]byte(topicMessage.Object))
	}
}

func (p ProxyServiceHandler) HandleCreate(message []byte) {
}

func (p ProxyServiceHandler) HandleDelete(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	for _, p := range svc.Spec.Ports {
		key := svc.Status.ClusterIP + ":" + strconv.Itoa(int(p.Port))
		ipvs.DeleteService(key)

		if svc.Spec.Type == apiobjects.ServiceTypeNodePort {
			var IP string
			if utils.GetLocalIP() == "192.168.1.12" {
				IP = "10.119.13.186"
				
			}else if utils.GetLocalIP() == "192.168.1.14" {
				IP = "10.119.13.140"
				
			}else if utils.GetLocalIP() == "192.168.1.15"{
				IP = "10.119.13.254"
			}
			key := IP + ":" + strconv.Itoa(int(p.Port))
			ipvs.DeleteService(key)
		}
	}
}

func (p ProxyServiceHandler) HandleUpdate(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	for _, p := range svc.Spec.Ports {
		ipvs.AddService(svc.Status.ClusterIP, uint16(p.Port))

		if svc.Spec.Type == apiobjects.ServiceTypeNodePort {
			var IP string
			if utils.GetLocalIP() == "192.168.1.12" {
				IP = "10.119.13.186"
				
			}else if utils.GetLocalIP() == "192.168.1.14" {
				IP = "10.119.13.140"
				
			}else if utils.GetLocalIP() == "192.168.1.15"{
				IP = "10.119.13.254"
			}
			ipvs.AddService(IP, uint16(p.Port))
		}
	}
}

func (p ProxyServiceHandler) GetType() string {
	return "proxyserviceHandler"
}

/* ========== Endpoint Handler ========== */

type ProxyEndpointHandler struct {
}

func (e ProxyEndpointHandler) HandleCreate(message []byte) {
	edpt := &apiobjects.Endpoint{}
	edpt.UnMarshalJSON(message)
	if edpt.Spec.SvcIP == "HostIP" {
		var IP string
		if utils.GetLocalIP() == "192.168.1.12" {
			IP = "10.119.13.186"
			
		}else if utils.GetLocalIP() == "192.168.1.14" {
			IP = "10.119.13.140"
			
		}else if utils.GetLocalIP() == "192.168.1.15"{
			IP = "10.119.13.254"
		}
		log.Info("[proxy] Add HostIP Endpoint: svcIP:", edpt.Spec.SvcIP, "SvcPort:", edpt.Spec.SvcPort, "DestIP:", edpt.Spec.DestIP, "DestPort:", edpt.Spec.DestPort)
		edpt.Spec.SvcIP = IP
		ipvs.AddService(edpt.Spec.SvcIP, uint16(edpt.Spec.SvcPort))
		key := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
		ipvs.AddEndpoint(key, edpt.Spec.DestIP, uint16(edpt.Spec.DestPort))
		return
	}

	key := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
	log.Info("[proxy] Add Endpoint: svcIP:", edpt.Spec.SvcIP, "SvcPort:", edpt.Spec.SvcPort)
	ipvs.AddEndpoint(key, edpt.Spec.DestIP, uint16(edpt.Spec.DestPort))
}

func (e ProxyEndpointHandler) HandleDelete(message []byte) {
	edpt := &apiobjects.Endpoint{}
	edpt.UnMarshalJSON(message)
	if edpt.Spec.SvcIP == "HostIP" {
		var IP string
		if utils.GetLocalIP() == "192.168.1.12" {
			IP = "10.119.13.186"
			
		}else if utils.GetLocalIP() == "192.168.1.14" {
			IP = "10.119.13.140"
			
		}else if utils.GetLocalIP() == "192.168.1.15"{
			IP = "10.119.13.254"
		}
		edpt.Spec.SvcIP = IP
		svcKey := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
		dstKey := edpt.Spec.DestIP + ":" + strconv.Itoa(int(edpt.Spec.DestPort))
		ipvs.DeleteEndpoint(svcKey, dstKey)
	}

	svcKey := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
	dstKey := edpt.Spec.DestIP + ":" + strconv.Itoa(int(edpt.Spec.DestPort))
	ipvs.DeleteEndpoint(svcKey, dstKey)
}

func (e ProxyEndpointHandler) HandleUpdate(message []byte) {

}

func (e ProxyEndpointHandler) GetType() string {
	return "proxyendpointHandler"
}

func CheckAllServiceAndEndpoint(){
	utils.Info("CheckAllServiceandEndpoint")
	svc_list := []*apiobjects.Service{}
	err := utils.GetUnmarshal(route.Prefix+"/api/get/allservices", &svc_list)
	if err != nil {
		utils.Error("get svc list error")
	}

	// edptList := []*apiobjects.Endpoint{}
	// err = utils.GetUnmarshal(route.Prefix + route.GetAllEndpointsPath, &edptList)
	// if err != nil {
	// 	utils.Info("[ServiceController] get all endpoints error")
	// }
	for _, svc := range svc_list {
		for _, p := range svc.Spec.Ports {
			ipvs.AddService(svc.Status.ClusterIP, uint16(p.Port))
	
			if svc.Spec.Type == apiobjects.ServiceTypeNodePort {
				var IP string
				if utils.GetLocalIP() == "192.168.1.12" {
					IP = "10.119.13.186"
					
				}else if utils.GetLocalIP() == "192.168.1.14" {
					IP = "10.119.13.140"
					
				}else if utils.GetLocalIP() == "192.168.1.15"{
					IP = "10.119.13.254"
				}
				ipvs.AddService(IP, uint16(p.Port))
			}
		}
	}

}
