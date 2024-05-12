package internal

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func SendPodStatus(pods map[string]*PodWrapper) (err error) {
	// Send changed pods to apiserver
	// TODO: kubelet 使用运行时传入的参数，而不是直接使用全局变量指定apiserver地址
	for _, pod := range pods {
		utils.PutWithJson(route.Prefix + route.PodStatePath, pod.Pod)
	}

	return
}

func SendNodeStatus(node *apiobjects.Node) (err error) {
	// Send node status to apiserver
	_, err = utils.PutWithJson(route.Prefix + route.NodePath, node)
	return
}
