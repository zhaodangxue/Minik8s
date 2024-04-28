package global

const podUpdateTopic = "pod-update"
const podStatusTopic = "pod-status"
const testTopic = "test"

func PodUpdateTopic(name string) string {
	return podUpdateTopic + "-" + name
}
func PodStatusTopic() string {
	return podStatusTopic
}
func TestTopic() string {
	return testTopic
}
