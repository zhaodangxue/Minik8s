package global

const (
	nodeStateTopic = "node-state" // 专门用于发布Node状态变化
)

func NodeStateTopic() string {
	return nodeStateTopic
}
