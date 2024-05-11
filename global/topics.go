package global

const (
	podStateTopic 		 	= "pod-state"            // 专门用于发布Pod状态变化
	testTopic               = "test"                 // 测试
	bindingTopic            = "binding"              // Pod和Node的绑定
	serviceTopic            = "service"              // Service的增删改查
	strategyUpdateTopic     = "strategy-update"      // 策略更新
	podRelevantTopic        = "pod-relevant"         // Pod相关的信息
	pvRelevantTopic         = "pv-relevant"          // PV相关的信息
	pvcRelevantTopic        = "pvc-relevant"         // PVC相关的信息
)
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
func PodStateTopic() string {
	return podStateTopic
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
