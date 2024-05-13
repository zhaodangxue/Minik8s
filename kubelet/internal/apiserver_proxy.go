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

	responseHttp, err := utils.PostWithJson(route.Prefix+route.NodeHealthPath, request)
	if err != nil {
		utils.Error("SendHealthReport PostWithJson error:", err)
	}
	response := api.NodeHealthReportResponse{}
	if err := json.NewDecoder(responseHttp.Body).Decode(&response); err != nil {
		utils.Error("SendHealthReport Decode error:", err)
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
