package hpacontroller

import (
	"minik8s/apiobjects"
	"minik8s/utils"
)

func GetMiddle(x int, min int, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
func ScaleByCPUPercent(rs *apiobjects.Replicaset, min int, max int, target int) (int, int) {
	var replicaset_cpu_usage float32
	replicaset_cpu_usage = rs.Stat.AverageCpuPercent
	//乘以100并转化为int
	var current_cpu_usage int
	var replicas int
	current_cpu_usage = int(replicaset_cpu_usage * 100)
	if current_cpu_usage == 0 {
		replicas = min
	} else {
		ratio := float32(current_cpu_usage) / float32(target)
		replicas = GetMiddle(int(ratio*float32(rs.Spec.Replicas)), min, max)
	}
	utils.Info("current_cpu_usage:", current_cpu_usage, "replicas:", replicas)
	return current_cpu_usage, replicas
}
func ScaleByMemPercent(rs *apiobjects.Replicaset, min int, max int, target int) (int, int) {
	var replicaset_mem_usage float32
	replicaset_mem_usage = rs.Stat.AverageMemPercent
	//乘以100并转化为int
	var current_mem_usage int
	var replicas int
	current_mem_usage = int(replicaset_mem_usage * 100)
	if current_mem_usage == 0 {
		replicas = min
	} else {
		ratio := float32(current_mem_usage) / float32(target)
		replicas = GetMiddle(int(ratio*float32(rs.Spec.Replicas)), min, max)
	}
	utils.Info("current_mem_usage:", current_mem_usage, "replicas:", replicas)
	return current_mem_usage, replicas
}
