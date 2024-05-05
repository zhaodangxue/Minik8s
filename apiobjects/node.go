package apiobjects

type Node struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata"`
	Info       NodeInfo   `yaml:"info"`
	Status     NodeStatus `yaml:"status"`
}

type NodeInfo struct {
	Ip string `yaml:"ip"`
}

type NodeStatus struct {
	State string `yaml:"state"`
}

const (
	NodeStateHealthy   = "Healthy"
	NodeStateUnhealthy = "Unhealthy"
)
