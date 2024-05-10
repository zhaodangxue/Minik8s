package internal

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func SendPodStatus(changedPods []*apiobjects.Pod) (err error) {
	// Send changed pods to apiserver
	for _, pod := range changedPods {
		utils.PostWithJson(route.Prefix + route.NodePath, pod)
	}
}
