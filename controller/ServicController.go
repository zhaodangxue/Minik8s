package controller

import (
	"fmt"
	"minik8s/apiobjects"
	"minik8s/utils"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

/* 主要工作：
1. 监听service资源的创建。一旦service资源创建，为其分配唯一的cluster ip。
2. 遍历pod列表，找到符合selector条件的pod，记录。创建endpoint。
3. 监听pod创建。增加endpoint。
4. 监听pod删除。删除endpoint。
5. 监听pod更新。如果标签更改，删除/增加endpoint。
6. 监听service资源的删除。删除对应endpoint。
*/

var IPMap = [1 << 8]bool{false}
var IPStart = "10.10.0."

var svcToEndpoints = map[string]*[]*apiobjects.Endpoint{}
var svcList = map[string]*apiobjects.Service{}

type svcServiceHandler struct {
}

type svcPodHandler struct {
}

/* ========== Service Handler ========== */

func (s svcServiceHandler) HandleCreate(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	//分配cluster ip，更新serviceList
	svc.Status.ClusterIP = allocateClusterIP()
	svcList[svc.Status.ClusterIP] = svc
	//TODO 发送http给apiserver,更新service,带有分配好的cluster ip
	svcByte,err := svc.MarshalJSON()
	if err != nil {
		fmt.Println("error")
	}
	//response, err :=utils.PostWithString("http://apiserver:8080/api/service/"+svc.Data.Namespace+"/"+svc.Data.Name, string(svcByte))
	response, err :=utils.PostWithString("http://apiserver:8080/api/service", string(svcByte))
	if err != nil{
		print("create service error")
	}
	fmt.Println(response)

	//TODO 遍历pod列表，找到符合selector条件的pod，记录并创建该svc对应的endpoint。
	createEndpointsFromPodList(svc)

	log.Info("[svc controller] Create service. Cluster IP:", svc.Status.ClusterIP)
}

func (s svcServiceHandler) HandleDelete(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)
	delete(svcList, svc.Status.ClusterIP)
	index := strings.SplitN(svc.Status.ClusterIP, ",", -1)
	indexLast, _ := strconv.Atoi(index[len(index)-1])
	print(indexLast)
	IPMap[indexLast] = false
    //删除service
	response, err :=utils.Delete("http://apiserver:8080/api/service/delete/"+svc.Data.Namespace+"/"+svc.Data.Name)
	if err != nil{
		print("delete service error")
	}
	fmt.Println(response)

	//todo 删除对应的endpoints
	for _, edpt := range *svcToEndpoints[svc.Status.ClusterIP] {
		response, err :=utils.Delete("http://apiserver:8080/api/endpoint/"+edpt.Data.Namespace+"/"+edpt.Data.Name)
		if err != nil{
			print("delete endpoints error")
		}
		fmt.Println(response)
	}
	delete(svcToEndpoints, svc.Status.ClusterIP)
	log.Info("[svc controller] Delete service. Cluster IP:", svc.Status.ClusterIP)
}

func (s svcServiceHandler) HandleUpdate(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	oldSvc, ok := svcList[svc.Status.ClusterIP]
	if !ok {
		log.Info("[svc controller] Service not found. ClusterIP:", svc.Status.ClusterIP, "\n")
		return
	}
	//TODO检查标签是否更改。如果是，则删除旧的endpoints并添加新的endpoints
	if !IsLabelEqual(oldSvc.Spec.Selector, svc.Spec.Selector) {
		for _, edpt := range *svcToEndpoints[svc.Status.ClusterIP] {
			response, err :=utils.Delete("http://apiserver:8080/api/v1/endpoint/"+edpt.Data.Namespace+"/"+edpt.Data.Name)
			if err != nil{
				print("delete endpoints error")
			}
			fmt.Println(response)
		}
		createEndpointsFromPodList(svc)
	}

	svcList[svc.Status.ClusterIP] = svc
	log.Info("[svc controller] Update service. Cluster IP:", svc.Status.ClusterIP)
}

func (s svcServiceHandler) GetType() string{
	return "service"
}

/* ========== Pod Handler ========== */
func (s svcPodHandler) HandleCreate(message []byte) {

}

func (s svcPodHandler) HandleDelete(message []byte) {
	//pod := &apiobjects.Pod{}

	//todo删除对应的endpoints

}

func (s svcPodHandler) HandleUpdate(message []byte) {
	pod := &apiobjects.Pod{}
    //todo 获取更新后的Pod
	
    //遍历service列表，检查所有的service，如果有对应这个Pod的endpoint,且Pod更新之后不符合selector条件，则删除对应的endpoints
	//                                如果没有对应这个Pod的endpoint,但Pod更新之后符合selector条件，则创建对应的endpoints
	for _, svc := range svcList {
		exist := isEndpointExist(svcToEndpoints[svc.Status.ClusterIP], pod.Status.PodIP)
		fit := IsLabelEqual(svc.Spec.Selector, pod.Labels)
		if !exist && fit {
			createEndpoints(svcToEndpoints[svc.Status.ClusterIP], svc, pod)
		} else if exist && !fit {
			deleteEndpoints(svc, pod)
		}
		// TODO: handle the change of POD IP
	}
}

func (s svcPodHandler) GetType() string{
	return "pod"
}

/* ========== Util Function ========== */
func allocateClusterIP() string {
	for i, used := range IPMap {
		if i != 0 && !used {
			IPMap[i] = true
			return IPStart + strconv.Itoa(i)
		}
	}
	log.Fatal("[svc controller] Cluster IP used up!")
	return ""
}

func createEndpointsFromPodList(svc *apiobjects.Service) {
	//从apiserver获取pod列表
    //podList, err := utils.GetPodList()    

	var edptList []*apiobjects.Endpoint
	// for _, p := range podList {
    //     //筛选符合selector条件的pod
	// }

	//更新service对应的endpoints
	svcToEndpoints[svc.Status.ClusterIP] = &edptList
}

func IsLabelEqual(a map[string]string, b map[string]string) bool {
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func isEndpointExist(edptList *[]*apiobjects.Endpoint, podIP string) bool {
	for _, edpt := range *edptList {
		if edpt.Spec.DestIP == podIP {
			return true
		}
	}
	return false
}


func createEndpoints(edptList *[]*apiobjects.Endpoint, svc *apiobjects.Service, pod *apiobjects.Pod) {
	logInfo := "[svc controller] Create endpoints."

	for _, port := range svc.Spec.Ports {
		dstPort := findDstPort(port.TargetPort, pod.Spec.Containers)
		spec := apiobjects.EndpointSpec{
			SvcIP:    svc.Status.ClusterIP,
			SvcPort:  port.Port,
			DestIP:   pod.Status.PodIP,
			DestPort: dstPort,
		}
		edpt := &apiobjects.Endpoint{
			Spec: spec,
			Data: apiobjects.MetaData{
				Name:      svc.Data.Name + "-" + pod.Name,
				Namespace: svc.Data.Namespace,
			},
		}
		//TODO 发送http给apiserver,更新edpt

		*edptList = append(*edptList, edpt)
		logInfo += fmt.Sprintf("srcIP:%s:%d, dstIP:%s:%d ; ", svc.Status.ClusterIP, port.Port, pod.Status.PodIP, dstPort)
	}

	log.Info(logInfo)

}

func findDstPort(targetPort string, containers []apiobjects.Container) int32 {
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.Name == targetPort {
				return p.ContainerPort
			}
		}
	}
	log.Fatal("[svc controller] No Match for Target Port!")
	return 0
}

func deleteEndpoints(svc *apiobjects.Service, pod *apiobjects.Pod) {
	logInfo := "[svc controller] Delete endpoints."

	edptList := svcToEndpoints[svc.Status.ClusterIP]
	var newEdptList []*apiobjects.Endpoint
	for key, edpt := range *edptList {
		if edpt.Spec.DestIP == pod.Status.PodIP {
			edpt := (*edptList)[key]
			//TODO 发送http给apiserver,更新edpt

			logInfo += fmt.Sprintf("srcIP:%s:%d, dstIP:%s:%d ; ", edpt.Spec.SvcIP, edpt.Spec.SvcPort, edpt.Spec.DestIP, edpt.Spec.DestPort)
		} else {
			newEdptList = append(newEdptList, edpt)
		}
	}
	svcToEndpoints[svc.Status.ClusterIP] = &newEdptList

	log.Info(logInfo)
}
