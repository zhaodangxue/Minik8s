package scheduler

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	sched_utils "minik8s/scheduler/utils"
	"minik8s/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

const BindingCheckInterval = 5 * time.Second

type Scheduler interface {
	Start()
	Schedule(pod *apiobjects.Pod) error
}
type scheduler struct {
	labelSelector    sched_utils.LabelSelector
	strategySelector sched_utils.StrategySelector
}

func New() Scheduler {
	return &scheduler{
		labelSelector:    sched_utils.NewLabel(),
		strategySelector: sched_utils.NewStrategy(sched_utils.RandomStrategy),
	}
}
func (s *scheduler) GetAllNodesFromApiServer() (nodes []*apiobjects.Node) {
	err := utils.GetUnmarshal(route.Prefix+route.NodePath, &nodes)
	if err != nil {
		utils.Error(err)
	}
	utils.Info("get all nodes from apiserver")
	return
}
func (s *scheduler) SendScheduleInfoToApiServer(pod *apiobjects.Pod, node *apiobjects.Node) {
	binding := apiobjects.NodePodBinding{
		Node: *node,
		Pod:  *pod,
	}
	url := route.Prefix + "/api/binding" + "/" + pod.Namespace + "/" + pod.Name + "/" + node.ObjectMeta.Name
	_, err := utils.PostWithJson(url, binding)
	if err != nil {
		utils.Error(err)
	}
	_, err = utils.PutWithJson(route.Prefix+route.PodPath+"/"+pod.Namespace+"/"+pod.Name, pod)
	if err != nil {
		utils.Error(err)
	}
}
func (s *scheduler) selectStrategy(strategy byte) {
	s.strategySelector = sched_utils.NewStrategy(strategy)
}
func (s *scheduler) handleStrategyChange(msg *redis.Message) {
	strategy := msg.Payload
	var strategyType byte
	switch strategy {
	case "RandomStrategy":
		strategyType = sched_utils.RandomStrategy
	case "MininumCpuStrategy":
		strategyType = sched_utils.MininumCpuStrategy
	case "MininumMemStrategy":
		strategyType = sched_utils.MininumMemStrategy
	default:
		utils.Error("unknow strategy")
	}
	s.selectStrategy(strategyType)
	fmt.Printf("strategy change to %s\n", strategy)
}
func (s *scheduler) Schedule(pod *apiobjects.Pod) error {
	nodes := s.GetAllNodesFromApiServer()
	if len(nodes) == 0 {
		return fmt.Errorf("no available node")
	}
	filtedNodes := s.labelSelector.LabelSelector(pod, nodes)
	if len(filtedNodes) == 0 {
		return fmt.Errorf("no available node satisfy label selector")
	}
	node := s.strategySelector.StrategySelector(filtedNodes)
	if node == nil {
		return fmt.Errorf("no available node satisfy strategy")
	}
	s.SendScheduleInfoToApiServer(pod, node)
	// binding := apiobjects.NodePodBinding{
	// 	Node: *node,
	// 	Pod:  *pod,
	// }
	// updateMsg, _ := json.Marshal(binding)
	// topics := global.PodUpdateTopic(*pod)
	// listwatch.Publish(topics, string(updateMsg))
	fmt.Printf("schedule pod %s to node %s\n", pod.ObjectMeta.Name, node.ObjectMeta.Name)
	return nil
}
func (s *scheduler) doSchedule(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		utils.Error(err)
	}
	if topicMessage.ActionType == apiobjects.Delete {
		return
	}
	pod := &apiobjects.Pod{}
	err = json.Unmarshal([]byte(topicMessage.Object), pod)
	if err != nil {
		utils.Error(err)
	}
	err = s.Schedule(pod)
	if err != nil {
		utils.Error(err)
	}
}
func (s *scheduler) reSchedule(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		utils.Error(err)
	}
	if topicMessage.ActionType == apiobjects.Delete {
		nodepodbinding := &apiobjects.NodePodBinding{}
		err = json.Unmarshal([]byte(topicMessage.Object), nodepodbinding)
		if err != nil {
			utils.Error(err)
		}
		pod := &nodepodbinding.Pod
		err = s.Schedule(pod)
		if err != nil {
			utils.Error(err)
		}
	}
	return
}
func (s *scheduler) CheckPodBinding() {
	for{
		utils.Debug("check pod binding")
		var pods []*apiobjects.Pod
		err := utils.GetUnmarshal(route.Prefix+route.PodPath, &pods)
		if err != nil {
			utils.Error(err)
		}
		for _, pod := range pods {
			url := route.Prefix + "/api/binding" + "/" + pod.Namespace + "/" + pod.Name
			nodepodbinding := &apiobjects.NodePodBinding{}
			err := utils.GetUnmarshal(url, nodepodbinding)
			if err != nil {
				utils.Error(err)
			}
			if nodepodbinding.Pod.ObjectMeta.Name == "" {
				err = s.Schedule(pod)
			}
		}
		time.Sleep(BindingCheckInterval)
	}
}
func (s *scheduler) Start() {
	go listwatch.Watch(global.StrategyUpdateTopic(), s.handleStrategyChange)
	go listwatch.Watch(global.PodRelevantTopic(), s.doSchedule)
	go s.CheckPodBinding()
	listwatch.Watch(global.BindingTopic(), s.reSchedule)
}
