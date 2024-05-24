package internal

import (
	"minik8s/apiobjects"

	cpu "github.com/mackerelio/go-osstat/cpu"
	memory "github.com/mackerelio/go-osstat/memory"
)

// 获取节点的统计信息
func GetNodeStats() (nodeStats *apiobjects.NodeStats, err error) {

	memory, err := memory.Get()
	if err != nil {
		return nil, err
	}
	cpu, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	nodeStats = &apiobjects.NodeStats{
		CpuUsage: apiobjects.NodeCpuUsage{
			Total:  cpu.Total,
			Idle:   cpu.Idle,
			Iowait: cpu.Iowait,
		},
		MemUsage: apiobjects.NodeMemoryUsage{
			UsedBytes:      memory.Used,
			AvailableBytes: memory.Free,
		},
	}

	return
}
