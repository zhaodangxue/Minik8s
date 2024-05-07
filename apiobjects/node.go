package apiobjects

type Node struct {
	Object `yaml:",inline"`
	Info   NodeInfo   `yaml:"info"`
	Status NodeStatus `yaml:"status"`
}

type NodeInfo struct {
	Ip string `yaml:"ip"`
}

type NodeStatus struct {
	State      string  `yaml:"state"`
	CpuPercent float64 `yaml:"-"`
	MemPercent float64 `yaml:"-"`
}

const (
	NodeStateHealthy   = "Healthy"
	NodeStateUnhealthy = "Unhealthy"
)
