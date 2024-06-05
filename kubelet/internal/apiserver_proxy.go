package internal

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/api"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func SendHealthReport(node *apiobjects.Node, pods map[string]*apiobjects.Pod) (podsToDelete []*apiobjects.Pod, err error) {
	// Send changed pods to apiserver
	// TODO: kubelet 使用运行时传入的参数，而不是直接使用全局变量指定apiserver地址

	request := api.NodeHealthReportRequest{
		Node: *node,
		Pods: make([]*apiobjects.Pod, 0),
	}
	for _, pod := range pods {
		request.Pods = append(request.Pods, pod)
	}

	responseStr, err := utils.PutWithJson(route.Prefix+route.NodeHealthPath, request)
	if err != nil {
		utils.Error("SendHealthReport PutWithJson error:", err)
		return
	}
	response := api.NodeHealthReportResponse{}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		utils.Error("SendHealthReport json.Unmarshal error:", err, "\nresponseStr:", responseStr)
		return
	}
	
	for _, podPath := range response.UnmatchedPodPaths {
		podsToDelete = append(podsToDelete, pods[podPath])
	}

	return
}

func GetAllBindings() (bindings []apiobjects.NodePodBinding, err error) {
	// Get all bindings from apiserver
	err = utils.GetUnmarshal(route.Prefix+route.NodePodBindingAllPath, &bindings)
	if err != nil {
		return
	}
	return
}

func GetPodByPath(podPath string) (pod *apiobjects.Pod, err error) {
	// Get pod from apiserver
	err = utils.GetUnmarshal(route.Prefix+podPath, &pod)
	if err != nil {
		return
	}
	return
}
