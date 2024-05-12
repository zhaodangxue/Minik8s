package main

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/kubelet/internal"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// kubeletServer 用于存放kubelet的状态
//
// Node: 用于存放Node的信息
//
// Bindings: 用于存放Node和Pod的绑定关系
//
// PodCreateChan: 用于存放需要创建的Pod
type kubeletServer struct {
	Node apiobjects.Node
	// Pods 用于存放当前Pod的状态
	// key: Pod的Path
	Pods map[string]*internal.PodWrapper
	// PodCreateChan 用于通知kubelet主循环创建Pod
	PodCreateChan chan apiobjects.Pod

	// PodStatusCheckerChan 用于触发Pod状态检查
	PodStatusCheckerChan chan Empty
	// NodeHealthyReportChan 用于触发上报Node的健康状态
	NodeHealthyReportChan chan Empty
}

const (
	// Pod状态检查定时
	PodStatusCheckInterval = 10 * time.Second
	// Node健康状态上报定时
	NodeHealthyReportInterval = 10 * time.Second
)

// Empty 用于传递空消息
type Empty struct{}

var server kubeletServer = kubeletServer{}

func serverInit() {
	uid := uuid.NewString()
	name := "node-" + uid[:6]
	server.Node = apiobjects.Node{
		// TypeMeta: apiobjects.TypeMeta{
		// 	ApiVersion: global.ApiVersion,
		// 	Kind:       "Node",
		// },
		// ObjectMeta: apiobjects.ObjectMeta{
		// 	Name:              "",
		// 	Namespace:         global.SystemNamespace,
		// 	UID:               uuid.NewString(),
		// 	Labels:            map[string]string{},
		// 	CreationTimestamp: time.Now(),
		// 	DeletionTimestamp: time.Time{},
		// },
		Object: apiobjects.Object{
			TypeMeta: apiobjects.TypeMeta{
				ApiVersion: global.ApiVersion,
				Kind:       "Node",
			},
			ObjectMeta: apiobjects.ObjectMeta{
				Name:              name,
				Namespace:         global.SystemNamespace,
				UID:               uid,
				Labels:            map[string]string{},
				CreationTimestamp: time.Now(),
				DeletionTimestamp: time.Time{},
			},
		},
		Info: apiobjects.NodeInfo{
			Ip: utils.GetLocalIP(),
		},
		Status: apiobjects.NodeStatus{
			State: apiobjects.NodeStateHealthy,
		},
	}
	server.Pods = make(map[string]*internal.PodWrapper)
	server.PodCreateChan = make(chan apiobjects.Pod, 100)
	server.PodStatusCheckerChan = make(chan Empty, 1)
	server.NodeHealthyReportChan = make(chan Empty, 1)

	// TODO: 解决pod启动问题
	// 获取server的Bindings，或通知apiserver Node重启(通过node的状态变化)

	// TODO: 通知apiserver更新node状态
}

func onBingdingUpdate(message *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(message.Payload), &topicMessage)
	if err != nil {
		utils.Error("kubelet:onBingdingUpdate err=", err)
		return
	}

	switch topicMessage.ActionType {
	case apiobjects.Create:
		binding := apiobjects.NodePodBinding{}
		err := json.Unmarshal([]byte(topicMessage.Object), &binding)
		if err != nil {
			utils.Error("kubelet:onBingdingUpdate parsing create binding, err=", err)
			return
		}
		if binding.Node.GetObjectRef() != server.Node.GetObjectRef() {
			utils.Warn("kubelet:onBingdingUpdate node not match, binding.Node.Name=", binding.Node.ObjectMeta.Name)
			return
		}
		utils.Info("kubelet:onBingdingUpdate create pod with binding=", binding)
		server.PodCreateChan <- binding.Pod
	case apiobjects.Update:
		// TODO
		utils.Warn("kubelet:onBingdingUpdate Update not implemented")
	case apiobjects.Delete:
		// TODO
		utils.Warn("kubelet:onBingdingUpdate Delete not implemented")
	default:
		utils.Warn("kubelet:onBingdingUpdate unknown actionType=", topicMessage.ActionType)
	}
}

func podCreateHandler(pod apiobjects.Pod) {
	// FIXME: 考虑多线程同步
	utils.Info("kubelet:podCreateHandler pod=", pod)

	pod.Status.HostIP = server.Node.Info.Ip

	PodSandboxId, err := internal.CreatePod(pod)

	if err != nil {
		utils.Error("kubelet:podCreateHandler CreatePod error:", err)
	}

	podWrapper := internal.PodWrapper{Pod: pod, PodSandboxId: PodSandboxId}

	server.Pods[pod.GetObjectPath()] = &podWrapper

	// 等待pod状态检查线程自动检查
}

// 在指定时间间隔内在信道上发送空消息
func timedInformer(ch chan Empty, interval time.Duration) {
	for {
		ch <- Empty{}
		time.Sleep(interval)
	}
}

// 定时被调用，检查pod状态
func podStatusChecker() {
	utils.Info("kubelet:podStatusChecker")
	podsNotInCluster, err := internal.GetAllPods()
	if err != nil {
		utils.Error("kubelet:podStatusChecker GetAllPods error:", err)
		return
	}

	// Check all pods
	for _, pod := range podsNotInCluster {
		podWrapper, ok := server.Pods[pod.Pod.GetObjectPath()]
		if !ok {
			utils.Warn("kubelet:podStatusChecker running pod not in kubelet internal list: ", pod)
			
			utils.Info("kubelet:podStatusChecker adding pod: ", pod.Pod.GetObjectPath())
			new_pod := new(internal.PodWrapper)
			new_pod.PodSandboxId = pod.PodSandboxId
			server.Pods[pod.Pod.GetObjectPath()] = new_pod
			podWrapper = new_pod
		}
		pod.Pod.Status.HostIP = server.Node.Info.Ip
		podWrapper.Pod = pod.Pod
	}

	// Send all pods to apiserver
	podsNotInCluster, err = internal.SendPodStatus(server.Pods)
	if err != nil {
		utils.Error("kubelet:podStatusChecker SendPodStatus error:", err)
	} else {
		// Delete pods not in cluster
		for _, pod := range podsNotInCluster {
			utils.Info("kubelet:podStatusChecker deleting pod not in cluster: ", pod.Pod.GetObjectPath())
			delete(server.Pods, pod.Pod.GetObjectPath())
			internal.DeletePod(pod);
		}
	}

}

// 定时被调用，上报node的健康状态
func nodeHealthyReport() {
	// TODO: 定时被调用，上报node的健康状态
	utils.Info("kubelet:nodeHealthyReport")
	node := &server.Node
	node.Status.State = apiobjects.NodeStateHealthy
	err := internal.SendNodeStatus(node)
	if err != nil {
		utils.Error("kubelet:nodeHealthyReport SendNodeStatus error:", err)
	}
}

func main() {

	serverInit()

	go listwatch.Watch(global.BindingTopic(), onBingdingUpdate)
	go timedInformer(server.PodStatusCheckerChan, PodStatusCheckInterval)
	go timedInformer(server.NodeHealthyReportChan, NodeHealthyReportInterval)

	for {
		select {
		case pod := <-server.PodCreateChan:
			podCreateHandler(pod)
		case <-server.PodStatusCheckerChan:
			podStatusChecker()
		case <-server.NodeHealthyReportChan:
			nodeHealthyReport()
		}
	}
}
