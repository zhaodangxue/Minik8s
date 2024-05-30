package apiobjects

type EventType string

const (
	EventTypeTimer EventType = "timer"
)

type WorkflowRef string

// Event 事件
type Event struct {
	Name string `json:"name"`
	// 事件类型
	Type EventType `json:"type"`
	// 必须考虑悬挂引用的情况
	Workflows []WorkflowRef `json:"workflows"`
	// 定时器事件, 只有当Type为timer时有效
	TimeEvent *TimeEvent `json:"timeEvent,omitempty"`
}

type TimerType string

const (
	// 一次性定时器
	// 只在StartTime时触发一次
	TimerTypeOnce TimerType = "once"
	// 循环定时器
	// 从StartTime开始，每隔Interval秒触发一次
	// StartTime小于当前时间时，立即触发第一次
	TimerTypeLoop TimerType = "loop"
)

type TimeEvent struct {
	Cron string `json:"cron"`
}
