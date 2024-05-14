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

	listwatch.Watch(global.ServiceCmdTopic(), ss.HandleService)
	 go listwatch.Watch(global.PodStateTopic(), se.HandleEndpoints)

}

func main() {
	/* service controller */
	Run()
}

