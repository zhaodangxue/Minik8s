package global

//import "minik8s/apiobjects"

const (
	podStateTopic       = "pod-state"       // 专门用于发布Pod状态变化
	testTopic           = "test"            // 测试
	bindingTopic        = "binding"         // Pod和Node的绑定
	serviceTopic        = "service"         // servicecontroller处理之后对Service真正执行增删改查之后publish，由各个node上的kube-proxy监听
	serviceCmdTopic     = "service-cmd"     // apiserver收到kubctl对Service的增删改查命令,初步处理之后publish，由servicecontroller监听
	endpointTopic       = "endpoint"        // Endpoint的增删
	strategyUpdateTopic = "strategy-update" // 策略更新
	podRelevantTopic    = "pod-relevant"    // Pod相关的信息
	pvRelevantTopic     = "pv-relevant"     // PV相关的信息
	pvcRelevantTopic    = "pvc-relevant"    // PVC相关的信息
)

//	func PodUpdateTopic(pod apiobjects.Pod) string {
//		return podUpdateTopic + "-" + pod.Name + "-" + pod.Namespace
//	}
//
//	func PodStatusTopic() string {
//		return podStatusTopic
//	}
func PodStateTopic() string {
	return podStateTopic
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
func PodRelevantTopic() string {
	return podRelevantTopic
}
func PvRelevantTopic() string {
	return pvRelevantTopic
}
func PvcRelevantTopic() string {
	return pvcRelevantTopic
}
