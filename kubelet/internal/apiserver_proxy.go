package internal

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/api"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func SendPodStatus(pods map[string]*apiobjects.Pod) (podsToDelete []*apiobjects.Pod, err error) {
	// Send changed pods to apiserver
	// TODO: kubelet 使用运行时传入的参数，而不是直接使用全局变量指定apiserver地址

	podsToDelete = make([]*apiobjects.Pod, 0)

	for _, pod := range pods {
		responseStr, err := utils.PutWithJson(route.Prefix + route.PodStatePath, pod)
		if err != nil {
			utils.Error("SendPodStatus: PutWithJson failed: ", err)
			continue
		}
		httpError := api.HttpError{}
		err = json.Unmarshal([]byte(responseStr), &httpError)
		if err != nil {
			utils.Error("SendPodStatus: json.Unmarshal failed: ", err)
			continue
		}
		if httpError.Code == api.ApiserverErrorCode_UPDATE_POD_NOT_FOUND {
			podsToDelete = append(podsToDelete, pod)
		} else if httpError.Code != api.ApiserverErrorCode_NO_ERROR {
			utils.Error("SendPodStatus: httpError: ", httpError)
		}
	}

	return
}

func SendNodeStatus(node *apiobjects.Node) (err error) {
	// Send node status to apiserver
	_, err = utils.PutWithJson(route.Prefix + route.NodePath, node)
	return
}
