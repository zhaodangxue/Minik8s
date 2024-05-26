package replicasetcontroller

import (
	"fmt"
	"sync"
	"time"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/utils"
)

type Worker interface {
	Run()
	SyncCh() chan<- struct{}
	ResetTarget(target *apiobjects.Replicaset)
	ScaleTarget(target *apiobjects.Replicaset)
	SetPods(pods []*apiobjects.Pod)
	GetMtx() *sync.Mutex
	Done()
}
type worker struct {
	syncCh     chan struct{}
	target     *apiobjects.Replicaset
	pods       []*apiobjects.Pod
	mtx        sync.Mutex // 这个锁控制pod的访问
	global_mtx sync.Mutex
}

func (c *worker) AddPodToApiserver() {
	podTemplate := c.target.Spec.Template
	pod := apiobjects.ToPod(&podTemplate)
	pod.Namespace = c.target.Namespace
	pod.ObjectMeta.UID = utils.NewUUID()
	pod.ObjectMeta.Name = c.target.Name + "-" + pod.ObjectMeta.UID
	pod.AddLabel(global.ReplicasetLabel, c.target.ObjectMeta.UID)
	url := route.Prefix + route.PodPath
	utils.Info("replicaset worker create pod", pod.ObjectMeta.Name)
	_, err := utils.PostWithJson(url, pod)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *worker) DeletePodToApiserver(name, namespace string) {
	fmt.Printf("delete pod %s in namespace %s\n", name, namespace)
	url := route.Prefix + route.PodPath + "/" + namespace + "/" + name
	_, err := utils.Delete(url)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *worker) GetPodsByReplicasetUID() []*apiobjects.Pod {
	c.mtx.Lock()
	var Pods []*apiobjects.Pod
	for _, pod := range c.pods {
		val, exist := pod.Labels[global.ReplicasetLabel]
		if exist && val == c.target.ObjectMeta.UID {
			Pods = append(Pods, pod)
		}
	}
	c.mtx.Unlock()
	return Pods
}

func (c *worker) NumPods(pods []*apiobjects.Pod) (int, float32, float32) {
	c.mtx.Lock()
	count := len(pods)
	var cpu float32
	var mem float32
	for _, pod := range pods {
		cpu += pod.Stats.CpuUsage.GetCpuUsage()
		mem += pod.Stats.MemoryUsage.GetMemUsage()
	}
	var AverageCPUPercent float32
	var AverageMemPercent float32
	if count == 0 {
		AverageCPUPercent = 0
		AverageMemPercent = 0
	} else {
		AverageCPUPercent = cpu / float32(count)
		AverageMemPercent = mem / float32(count)
	}
	c.mtx.Unlock()
	return count, AverageCPUPercent, AverageMemPercent
}

func (c *worker) UpdateReplicasetReady(rs *apiobjects.Replicaset) {
	url := c.target.GetObjectPath()
	val, err := utils.PutWithJson(route.Prefix+url, rs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val)
}

func (c *worker) SyncLoop() bool {
	c.global_mtx.Lock()
	expected_num := c.target.Spec.Replicas
	pods := c.GetPodsByReplicasetUID()
	num_pods, AverageCPUPercent, AverageMemUsage := c.NumPods(pods)
	diff := expected_num - num_pods
	fmt.Printf("expected_num: %d, num_run: %d, diff: %d\n", expected_num, num_pods, diff)
	if diff > 0 {
		go c.AddPodToApiserver()
	}
	if diff < 0 {
		go c.DeletePodToApiserver(pods[0].Name, pods[0].Namespace)
	}
	c.target.Spec.Ready = num_pods
	c.target.Stat.AverageCpuPercent = AverageCPUPercent
	c.target.Stat.AverageMemUsage = AverageMemUsage
	c.UpdateReplicasetReady(c.target)
	timeout := time.NewTimer(20 * time.Second)
	c.global_mtx.Unlock()
	select {
	case _, open := <-c.syncCh:
		if !open {
			return false
		}
		return true
	case <-timeout.C:
		return true
	}
}

func (c *worker) Run() {
	// TODO
	for c.SyncLoop() {
	}
}

func (c *worker) Done() {
	var pods []*apiobjects.Pod
	err := utils.GetUnmarshal(route.Prefix+route.PodPath, &pods)
	if err != nil {
		fmt.Println(err)
	}
	c.mtx.Lock()
	c.pods = pods
	c.mtx.Unlock()
	pods = c.GetPodsByReplicasetUID()
	for _, pod := range pods {
		utils.Info("replicaset worker delete pod", pod.Name)
		c.DeletePodToApiserver(pod.Name, pod.Namespace)
	}
}

func (c *worker) SyncCh() chan<- struct{} {
	return c.syncCh
}

func (c *worker) ResetTarget(target *apiobjects.Replicaset) {
	pods := c.GetPodsByReplicasetUID()
	c.global_mtx.Lock()
	c.target = target
	c.global_mtx.Unlock()
	for _, pod := range pods {
		c.DeletePodToApiserver(pod.Name, pod.Namespace)
	}
	utils.Info("replicaset worker reset target", target.Name)
}

func (c *worker) SetPods(pods []*apiobjects.Pod) {
	c.mtx.Lock()
	c.pods = pods
	c.mtx.Unlock()
}

func (c *worker) GetMtx() *sync.Mutex {
	return &c.mtx
}
func (c *worker) ScaleTarget(target *apiobjects.Replicaset) {
	c.global_mtx.Lock()
	c.target = target
	c.global_mtx.Unlock()
	utils.Info("replicaset worker scale target", target.Name)
}
func NewWorker(target *apiobjects.Replicaset) Worker {
	return &worker{
		syncCh:     make(chan struct{}, 3),
		target:     target,
		mtx:        sync.Mutex{},
		global_mtx: sync.Mutex{},
		pods:       []*apiobjects.Pod{},
	}
}
