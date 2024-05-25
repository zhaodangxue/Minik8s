package registry

import (
	"minik8s/controller/api"
	"minik8s/controller/src/PVcontroller"
	replicasetcontroller "minik8s/controller/src/ReplicasetController"
	"minik8s/controller/src/health"
	service "minik8s/controller/src/service_controller"
	"minik8s/controller/src/node"
)

var ControllerList []api.Controller = []api.Controller{
	&health.HealthController{},
	&service.ServiceController{},
	&node.NodeController{},
	&PVcontroller.PVcontroller{},
	&replicasetcontroller.ReplicasetController{},
}
