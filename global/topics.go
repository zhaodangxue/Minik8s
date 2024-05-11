package global

import "minik8s/apiobjects"

const (
	podUpdateTopic          = "pod-update"           // 已经存在的Pod的内部信息更新
	podStatusTopic          = "pod-status"           // Pod发生增减的时候，通知其他组件
	podStateTopic 		 	= "pod-state"            // Pod状态变化
	testTopic               = "test"                 // 测试
	bindingTopic            = "binding"              // Pod和Node的绑定
	serviceTopic            = "service"              // Service的增删改查
	strategyUpdateTopic     = "strategy-update"      // 策略更新
	schedulerPodUpdateTopic = "scheduler-pod-update" // Scheduler更新Pod信息
)

func PodUpdateTopic(pod apiobjects.Pod) string {
	return podUpdateTopic + "-" + pod.Name + "-" + pod.Namespace
}
func PodStatusTopic() string {
	return podStatusTopic
}
func TestTopic() string {
	return testTopic
}
func BindingTopic() string {
	return bindingTopic
}
func ServiceTopic() string {
	return serviceTopic
}
func StrategyUpdateTopic() string {
	return strategyUpdateTopic
}
func SchedulerPodUpdateTopic() string {
	return schedulerPodUpdateTopic
}

func PodStateTopic() string {
	return podStateTopic
}
