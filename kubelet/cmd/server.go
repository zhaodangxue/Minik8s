package main

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"minik8s/kubelet/internal"
	"time"


	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type kubeletServer struct {
	Node apiobjects.Node
	// NodePodBindings 用于存放Node和Pod的绑定关系
	// key: NodePodBinding.Name
	Bindings map[string]apiobjects.NodePodBinding
	PodCreateChan chan apiobjects.Pod
}

var server kubeletServer = kubeletServer{}

func serverInit(){
	server.Node = apiobjects.Node{
		TypeMeta: apiobjects.TypeMeta{
			ApiVersion: global.ApiVersion,
			Kind: "Node",
		},
		ObjectMeta: apiobjects.ObjectMeta{
			Name: "",
			Namespace: global.SystemNamespace,
			UID: uuid.NewString(),
			Labels: map[string]string{},
			CreationTimestamp: time.Now(),
			DeletionTimestamp: time.Time{},
		},
		Info: apiobjects.NodeInfo{
			Ip: utils.GetLocalIP(),
		},
		Status: apiobjects.NodeStatus{
			State: apiobjects.NodeStateHealthy,
		},
	}
	server.Bindings = make(map[string]apiobjects.NodePodBinding)
	server.PodCreateChan = make(chan apiobjects.Pod, 100)

	// TODO: 解决pod启动问题
	// 获取server的Bindings，或通知apiserver Node重启(通过node的状态变化)

	// TODO: 通知apiserver更新node状态
}

func onBingdingUpdate(message *redis.Message) {
	binding := apiobjects.NodePodBinding{}
	err := json.Unmarshal([]byte(message.Payload), &binding)
	if err != nil{
		utils.Error("kubelet:onBingdingUpdate err=", err)
		return
	}

	// OPT: 可以通过为每个Node设置不同的BindingTopic，减少不必要的消息处理
	if binding.Node.Name != server.Node.Name {
		utils.Warn("kubelet:onBingdingUpdate node not match, binding.Node.Name=", binding.Node.Name)
		return
	}

	if binding.Pod.Status.PodPhase != apiobjects.PodCreated {
		utils.Warn("kubelet:onBingdingUpdate wrong pod phase, binding.Pod.Status.PodPhase=", binding.Pod.Status.PodPhase)
		return
	}

	server.Bindings[binding.Name] = binding
	utils.Info("kubelet:onBingdingUpdate binding=", binding)

	server.PodCreateChan <- binding.Pod
}

func podCreateHandler(pod apiobjects.Pod) {
	// FIXME: 考虑多线程同步
	utils.Info("kubelet:podCreateHandler pod=", pod)

	internal.CreatePod(pod);

	// TODO: 通知apiserver更新pod状态
}

func main() {

	serverInit()

	listwatch.Watch(global.BindingTopic(), onBingdingUpdate)

	for {
		select {
		case pod := <-server.PodCreateChan:
			podCreateHandler(pod)
		}
	}
}