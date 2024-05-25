package registry

import (
	"minik8s/controller/api"
	hpacontroller "minik8s/controller/src/HPAController"
	"minik8s/controller/src/PVcontroller"
	replicasetcontroller "minik8s/controller/src/ReplicasetController"
	"minik8s/controller/src/health"
	"minik8s/controller/src/node"
	service "minik8s/controller/src/service_controller"
)

var ControllerList []api.Controller = []api.Controller{
	&health.HealthController{},
	&service.ServiceController{},
	&node.NodeController{},
	&PVcontroller.PVcontroller{},
	&replicasetcontroller.ReplicasetController{},
	&hpacontroller.HPAController{},
}
