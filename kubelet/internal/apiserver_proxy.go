package internal

import (
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func SendPodStatus(pods []*PodWrapper) (err error) {
	// Send changed pods to apiserver
	for _, pod := range pods {
		utils.PutWithJson(route.Prefix + route.PodStatePath, pod.Pod)
	}

	return
}
