package internal

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func SendPodStatus(pods []*apiobjects.Pod) (err error) {
	// Send changed pods to apiserver
	for _, pod := range pods {
		utils.PutWithJson(route.Prefix + route.NodePath, pod)
	}

	return
}
