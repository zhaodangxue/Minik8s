package global

import "minik8s/apiobjects"

const (
	podUpdateTopic          = "pod-update"           // 已经存在的Pod的内部信息更新
	podStatusTopic          = "pod-status"           // Pod发生增减的时候，通知其他组件
	testTopic               = "test"                 // 测试
	bindingTopic            = "binding"              // Pod和Node的绑定
	serviceTopic            = "service"              // servicecontroller处理之后对Service真正执行增删改查之后publish，由各个node上的kube-proxy监听
	serviceCmdTopic         = "service-cmd"          // apiserver收到kubctl对Service的增删改查命令,初步处理之后publish，由servicecontroller监听
	endpointTopic           = "endpoint"             // Endpoint的增删
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
func ServiceCmdTopic() string {
	return serviceCmdTopic
}
func EndpointTopic() string {
	return endpointTopic
}

func StrategyUpdateTopic() string {
	return strategyUpdateTopic
}
func SchedulerPodUpdateTopic() string {
	return schedulerPodUpdateTopic
}
