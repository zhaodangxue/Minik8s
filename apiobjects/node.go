package apiobjects

type Node struct {
	Object `yaml:",inline"`
	Info   NodeInfo   `yaml:"info"`
	Status NodeStatus `yaml:"status"`
	Stats  NodeStats  `yaml:"-"`
}

type NodeInfo struct {
	Ip string `yaml:"ip"`
}

type NodeStatus struct {
	State string `yaml:"state"`
}

type NodeStats struct {
	CpuUsage NodeCpuUsage
	MemUsage NodeMemoryUsage
}

type NodeCpuUsage struct {
	Total  uint64
	Idle   uint64
	Iowait uint64
}

func (cpu *NodeCpuUsage) GetCpuUsage() float32 {
	return float32(cpu.Total-cpu.Idle-cpu.Iowait) / float32(cpu.Total)
}

type NodeMemoryUsage struct {
	UsedBytes      uint64
	AvailableBytes uint64
}

func (memory *NodeMemoryUsage) GetMemPercent() float32 {
	return float32(memory.UsedBytes) / float32(memory.AvailableBytes+memory.UsedBytes)
}

const (
	NodeStateHealthy   = "Healthy"
	NodeStateUnhealthy = "Unhealthy"
)
