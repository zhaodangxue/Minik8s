package apiobjects

type Node struct {
	TypeMeta
	ObjectMeta
	Info NodeInfo
	Status NodeStatus
}

type NodeInfo struct {
	Ip string
}

type NodeStatus struct {
	State string
}

const (
	NodeStateHealthy = "Healthy"
	NodeStateUnhealthy = "Unhealthy"
)
