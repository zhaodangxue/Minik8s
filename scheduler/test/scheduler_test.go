package sche__test

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	"minik8s/global"
	"minik8s/listwatch"
	scheduler "minik8s/scheduler/src"
	sched_utils "minik8s/scheduler/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLabelSelector(t *testing.T) {
	var pod apiobjects.Pod
	pod.Spec.NodeSelector = map[string]string{
		"app": "nginx",
		"env": "prod",
	}
	var nodes []*apiobjects.Node
	var node1 apiobjects.Node
	var node2 apiobjects.Node
	node1.ObjectMeta.Labels = map[string]string{
		"app": "nginx",
	}
	node1.ObjectMeta.Name = "node1"
	node2.ObjectMeta.Labels = map[string]string{
		"app": "nginx",
		"env": "prod",
	}
	node2.ObjectMeta.Name = "node2"
	nodes = append(nodes, &node1)
	nodes = append(nodes, &node2)
	ls := sched_utils.NewLabel()
	filted := ls.LabelSelector(&pod, nodes)
	assert.Equal(t, 1, len(filted))
	assert.Equal(t, "node2", filted[0].ObjectMeta.Name)
}
func TestStrategySelector(t *testing.T) {
	var nodes []*apiobjects.Node
	var node1 apiobjects.Node
	var node2 apiobjects.Node
	node1.ObjectMeta.Name = "node1"
	node2.ObjectMeta.Name = "node2"
	nodes = append(nodes, &node1)
	nodes = append(nodes, &node2)
	rs := sched_utils.NewStrategy(sched_utils.RandomStrategy)
	node := rs.StrategySelector(nodes)
	assert.NotNil(t, node)
}
func TestScheduler(t *testing.T) {
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	fmt.Println("apiServer is running")
	etcd.Clear()
	var node1 apiobjects.Node
	var node2 apiobjects.Node
	node1.ObjectMeta.Name = "node1"
	node1.ObjectMeta.Namespace = "minik8s-system"
	node1.TypeMeta.Kind = "Node"
	node2.ObjectMeta.Name = "node2"
	node2.ObjectMeta.Namespace = "minik8s-system"
	node2.TypeMeta.Kind = "Node"
	node1.Status.CpuPercent = 0.5
	node2.Status.CpuPercent = 0.3
	url_node1 := node1.Object.GetObjectPath()
	url_node2 := node2.Object.GetObjectPath()
	node1_JSON, _ := json.Marshal(node1)
	node2_JSON, _ := json.Marshal(node2)
	etcd.Put(url_node1, string(node1_JSON))
	etcd.Put(url_node2, string(node2_JSON))
	fmt.Println("node1 and node2 are added")
	scheduler := scheduler.New()
	go scheduler.Start()
	fmt.Println("scheduler is running")
	var pod apiobjects.Pod
	pod.ObjectMeta.Name = "pod1"
	pod.ObjectMeta.Namespace = "default"
	pod.ObjectMeta.UID = "123"
	pod.Spec.Containers = []apiobjects.Container{
		{
			Name:  "nginx",
			Image: "nginx",
		},
	}
	pod.TypeMeta.Kind = "Pod"
	pod_Json, _ := json.Marshal(pod)
	url_pod := pod.Object.GetObjectPath()
	etcd.Put(url_pod, string(pod_Json))
	listwatch.Publish(global.SchedulerPodUpdateTopic(), string(pod_Json))
	time.Sleep(5 * time.Second)
	url_binding := "/api/binding"
	val, _ := etcd.Get_prefix(url_binding)
	assert.Equal(t, 1, len(val))
	time.Sleep(5 * time.Second)
	listwatch.Publish(global.StrategyUpdateTopic(), "MininumCpuStrategy")
	pod.Spec.Containers[0].Name = "redis"
	pod.Spec.Containers[0].Image = "redis"
	pod_Json, _ = json.Marshal(pod)
	etcd.Put(url_pod, string(pod_Json))
	listwatch.Publish(global.SchedulerPodUpdateTopic(), string(pod_Json))
	time.Sleep(5 * time.Second)
	url_binding = "/api/binding"
	val, _ = etcd.Get_prefix(url_binding)
	assert.Equal(t, 1, len(val))
	var binding apiobjects.NodePodBinding
	json.Unmarshal([]byte(val[0]), &binding)
	assert.Equal(t, "node2", binding.Node.ObjectMeta.Name)
}
