package sched_utils

import (
	"math/rand"
	"minik8s/apiobjects"
)

// 第一个是我们的标签选择器接口
type LabelSelector interface {
	LabelSelector(pod *apiobjects.Pod, nodes []*apiobjects.Node) (filted []*apiobjects.Node)
}
type LabelSelectorImpl struct {
}

func NewLabel() LabelSelector {
	return &LabelSelectorImpl{}
}

// 第一个参数是pod中的label，第二个参数是指定node中的label
func isMatch(labels map[string]string, selector map[string]string) bool {
	for key, value := range selector {
		if val, ok := labels[key]; !ok || val != value {
			return false
		}
	}
	return true
}
func (ls *LabelSelectorImpl) LabelSelector(pod *apiobjects.Pod, nodes []*apiobjects.Node) (filted []*apiobjects.Node) {
	nodeSelector := pod.Spec.NodeSelector
	for _, n := range nodes {
		if isMatch(n.ObjectMeta.Labels, nodeSelector) {
			filted = append(filted, n)
		}
	}
	return
}

// 第二个是我们的选择策略接口，要求从多个node中选择一个node作为我们绑定的node
const (
	RandomStrategy byte = iota
	MininumCpuStrategy
	MininumMemStrategy
)

type StrategySelector interface {
	StrategySelector(nodes []*apiobjects.Node) *apiobjects.Node
}
type RandomStrategySelector struct {
}

func (rs *RandomStrategySelector) StrategySelector(nodes []*apiobjects.Node) *apiobjects.Node {
	sum := len(nodes)
	if sum == 0 {
		return nil
	}
	index := rand.Intn(sum)
	return nodes[index]
}

type MininumCpuStrategySelector struct {
}

func (mcs *MininumCpuStrategySelector) StrategySelector(nodes []*apiobjects.Node) *apiobjects.Node {
	if len(nodes) == 0 {
		return nil
	}
	min := nodes[0].Status.CpuPercent
	minNode := nodes[0]
	for _, node := range nodes {
		if node.Status.CpuPercent < min {
			min = node.Status.CpuPercent
			minNode = node
		}
	}
	return minNode
}

type MininumMemStrategySelector struct {
}

func (mms *MininumMemStrategySelector) StrategySelector(nodes []*apiobjects.Node) *apiobjects.Node {
	if len(nodes) == 0 {
		return nil
	}
	min := nodes[0].Status.MemPercent
	minNode := nodes[0]
	for _, node := range nodes {
		if node.Status.MemPercent < min {
			min = node.Status.MemPercent
			minNode = node
		}
	}
	return minNode
}
func NewStrategy(strategy byte) StrategySelector {
	switch strategy {
	case RandomStrategy:
		return &RandomStrategySelector{}
	case MininumCpuStrategy:
		return &MininumCpuStrategySelector{}
	case MininumMemStrategy:
		return &MininumMemStrategySelector{}
	default:
		panic("Unknown strategy")
	}
}
