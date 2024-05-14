package main
import (
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/controller"
)

func Run() {
	/* service controller */
	var ss controller.SvcServiceHandler
	var se controller.SvcEndpointHandler 

	go listwatch.Watch(global.PodStateTopic(), se.HandleEndpoints)
	go listwatch.Watch(global.BindingTopic(), se.HandleBindingEndpoints)
	listwatch.Watch(global.ServiceCmdTopic(), ss.HandleService)


}

func main() {
	/* service controller */
	Run()
}

