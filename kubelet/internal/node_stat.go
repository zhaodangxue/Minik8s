package internal

import (
	"minik8s/apiobjects"

	memory "github.com/mackerelio/go-osstat/memory"
	cpu "github.com/mackerelio/go-osstat/cpu"
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
			Total: cpu.Total,
			Idle: cpu.Idle,
			Iowait: cpu.Iowait,
		},
		MemUsage: apiobjects.NodeMemoryUsage{
			UsageBytes: memory.Used,
			AvailableBytes: memory.Free,
		},
	}

	return 
}
