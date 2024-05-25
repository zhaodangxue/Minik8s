package hpacontroller

import (
	"context"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"sync"
	"time"
)

type Worker interface {
	Run()
	SetTarget(target *apiobjects.HorizontalPodAutoscaler)
}
type worker struct {
	target *apiobjects.HorizontalPodAutoscaler
	ctx    context.Context
	mtx    sync.Mutex
}

func (c *worker) SetTarget(target *apiobjects.HorizontalPodAutoscaler) {
	if target != nil {
		c.mtx.Lock()
		c.target = target
		c.mtx.Unlock()
	}
}
func (c *worker) Run() {
	c.SyncLoop()
}
func (c *worker) GetRelevantReplicaset() *apiobjects.Replicaset {
	var rs apiobjects.Replicaset
	name := c.target.Spec.ScaleTargetRef.Name
	namespace := c.target.Spec.ScaleTargetRef.Namespace
	url := route.Prefix + route.ReplicasetPath + "/" + namespace + "/" + name
	err := utils.GetUnmarshal(url, &rs)
	if err != nil {
		utils.Error("get replicaset", name, "in namespace", namespace, "failed")
	}
	return &rs
}
func (c *worker) SyncLoop() {
	Interval := time.Second * time.Duration(c.target.Spec.ScaleInterval)
	tick := time.NewTicker(Interval)
	for c.Sync() {
		select {
		case <-tick.C:
			continue
		case <-c.ctx.Done():
			utils.Info("hpa worker", c.target.Name, "done")
			return
		}

	}
}
func (c *worker) Sync() bool {
	c.mtx.Lock()
	rs := c.GetRelevantReplicaset()
	fmt.Println("Sync")
	if rs == nil {
		utils.Error("HPA get relevant replicaset failed")
		return false
	}
	if rs.ObjectMeta.Name == "" && rs.ObjectMeta.Namespace == "" {
		utils.Error("HPA get relevant replicaset failed")
		return false
	}
	ExpectedCPUPercent := c.target.Spec.Metrics.CPUUtilizationPercentage
	ExpectedMemPercent := c.target.Spec.Metrics.MemoryUtilizationPercentage
	var NumReplicas int
	if ExpectedCPUPercent > 0 {
		current_cpu_usage, replicas := ScaleByCPUPercent(rs, c.target.Spec.MinReplicas, c.target.Spec.MaxReplicas, ExpectedCPUPercent)
		NumReplicas = replicas
		c.target.Stat.CurrnentReplicaseCPUUsage = current_cpu_usage
	} else {
		current_mem_usage, replicas := ScaleByMemPercent(rs, c.target.Spec.MinReplicas, c.target.Spec.MaxReplicas, ExpectedMemPercent)
		NumReplicas = replicas
		c.target.Stat.CurrentReplicaseMemUsage = current_mem_usage
	}
	diff := NumReplicas - rs.Spec.Replicas
	fmt.Printf("expected_num: %d, num_run: %d, diff: %d\n", NumReplicas, rs.Spec.Replicas, diff)
	if diff > 0 {
		go c.UpdateRelevantReplicaset(rs, rs.Spec.Replicas+1)
	} else if diff < 0 {
		go c.UpdateRelevantReplicaset(rs, rs.Spec.Replicas-1)
	}
	c.UpdateHPA()
	c.mtx.Unlock()
	return true
}
func (c *worker) UpdateRelevantReplicaset(rs *apiobjects.Replicaset, num int) {
	rs.Spec.Replicas = num
	url := route.Prefix + route.ReplicasetScale
	val, err := utils.PutWithJson(url, rs)
	if err != nil {
		utils.Error("HPA try to scale replicaset", rs.ObjectMeta.Name, "failed")
	}
	utils.Info(val)
	return
}
func (c *worker) UpdateHPA() {
	url := route.Prefix + route.HorizontalPodAutoscalerPath + "/" + c.target.Namespace + "/" + c.target.Name
	val, err := utils.PutWithJson(url, c.target)
	if err != nil {
		utils.Error("HPA try to update HPA", c.target.Name, "failed")
	}
	utils.Info(val)
	return
}
func NewWorker(ctx context.Context, target *apiobjects.HorizontalPodAutoscaler) Worker {
	return &worker{
		ctx:    ctx,
		target: target,
		mtx:    sync.Mutex{},
	}
}
