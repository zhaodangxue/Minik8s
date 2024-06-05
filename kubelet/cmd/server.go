package main

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/kubelet/internal"
	"minik8s/kubelet/internal/config"
	cri "minik8s/kubelet/internal/cri_proxy"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	Pods map[string]*apiobjects.Pod
	// PodCreateChan 用于通知kubelet主循环创建Pod
	PodCreateChan chan apiobjects.Pod

	// HealthReportChan 用于定时上报Node和Pod的状态
	HealthReportChan chan Empty
}

// Empty 用于传递空消息
type Empty struct{}

var server kubeletServer = kubeletServer{}

func generateServerName() string {
	return "node-" + utils.GetLocalIP()
}

func serverInit() {
	uid := uuid.NewString()
	name := generateServerName()
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
	server.Pods = make(map[string]*apiobjects.Pod)
	server.PodCreateChan = make(chan apiobjects.Pod, 100)
	server.HealthReportChan = make(chan Empty, 1)
}

func onBingdingUpdate(message *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(message.Payload), &topicMessage)
	if err != nil {
		utils.Error("kubelet:onBingdingUpdate err=", err)
		return
	}

	switch topicMessage.ActionType {
	case apiobjects.Update:
		utils.Info("kubelet:onBingdingUpdate update binding=", topicMessage.Object)
		fallthrough
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
	case apiobjects.Delete:
		binding := apiobjects.NodePodBinding{}
		err := json.Unmarshal([]byte(topicMessage.Object), &binding)
		if err != nil {
			utils.Error("kubelet:onBingdingUpdate parsing delete binding, err=", err)
			return
		}
		if binding.Node.GetObjectRef() != server.Node.GetObjectRef() {
			utils.Warn("kubelet:onBingdingUpdate node not match, binding.Node.Name=", binding.Node.ObjectMeta.Name)
			return
		}
		utils.Info("kubelet:onBingdingUpdate delete pod with binding=", binding)
		pod, ok := server.Pods[binding.Pod.GetObjectPath()]
		if !ok {
			utils.Warn("kubelet:onBingdingUpdate pod not found, binding.Pod.Name=", binding.Pod.ObjectMeta.Name)
			return
		}
		cri.DeletePod(pod.Status.SandboxId)
		delete(server.Pods, binding.Pod.GetObjectPath())
		utils.Info("kubelet:onBingdingUpdate delete pod with binding=", binding)
	default:
		utils.Warn("kubelet:onBingdingUpdate unknown actionType=", topicMessage.ActionType)
	}
}

func podCreateHandler(pod apiobjects.Pod) {
	// FIXME: 考虑多线程同步
	utils.Info("kubelet:podCreateHandler pod=", pod)

	pod.Status.HostIP = server.Node.Info.Ip

	err := cri.CreatePod(&pod)

	if err != nil {
		utils.Error("kubelet:podCreateHandler CreatePod error:", err)
	}

	newPod := new(apiobjects.Pod)
	*newPod = pod
	server.Pods[pod.GetObjectPath()] = newPod

	// 等待pod状态检查线程自动检查
}

// 在指定时间间隔内在信道上发送空消息
func timedInformer(ch chan Empty, interval time.Duration) {
	for {
		ch <- Empty{}
		time.Sleep(interval)
	}
}

// 定时被调用，上报Node和Pod的状态
func healthReport() {
	utils.Info("kubelet:healthReport")

	// Remove pods not in kubelet internal list
	podStatuses, err := cri.GetAllPods()
	if err != nil {
		utils.Error("kubelet:healthReport GetAllPods error:", err)
		return
	}
	for _, podStatus := range podStatuses {
		ref := cri.GetObjectRef(podStatus)
		_, ok := server.Pods[ref.GetObjectPath()]
		if !ok {
			utils.Warn("kubelet:healthReport running pod not in kubelet internal list: ", ref.GetObjectPath())
			// Check if the pod is in the cluster
			clusterBinding, err := internal.GetBindingByPath(ref)
			if err == nil {
				utils.Warn("kubelet:healthReport pod not in kubelet internal list but in cluster: ", ref.GetObjectPath())
				// Add the pod to kubelet internal list
				pod := clusterBinding.Pod
				pod.Status.SandboxId = podStatus.Status.Id
				pod.Status.HostIP = server.Node.Info.Ip
				server.Pods[ref.GetObjectPath()] = &pod
			} else {
				// CHECK: 考虑不同的出错可能
				utils.Warn("kubelet:healthReport pod not in kubelet internal list and not in cluster: ", ref.GetObjectPath())
				cri.DeletePod(podStatus.Status.Id)
			}
		}
	}

	// Update node status
	stats, err := internal.GetNodeStats()
	if err != nil {
		utils.Error("kubelet:healthReport GetNodeStats error:", err)
	}
	server.Node.Stats = *stats

	// Update pod status
	for _, pod := range server.Pods {
		cri.UpdatePodStatus(pod)
	}
	podsNotInCluster, err := internal.SendHealthReport(&server.Node, server.Pods)
	if err != nil {
		utils.Error("kubelet:healthReport SendPodStatus error:", err)
	} else {
		// Delete pods not in cluster
		for _, pod := range podsNotInCluster {
			utils.Info("kubelet:healthReport deleting pod not in cluster: ", pod.GetObjectPath())
			delete(server.Pods, pod.GetObjectPath())
			cri.DeletePod(pod.Status.SandboxId)
		}
	}

}

func main() {

	serverInit()

	go listwatch.Watch(global.BindingTopic(), onBingdingUpdate)
	go timedInformer(server.HealthReportChan, config.HealthReportInterval)

	for {
		select {
		case pod := <-server.PodCreateChan:
			podCreateHandler(pod)
		case <-server.HealthReportChan:
			healthReport()
		}
	}
}
