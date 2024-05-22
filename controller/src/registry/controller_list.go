package registry

import (
	"minik8s/controller/api"
	"minik8s/controller/src/PVcontroller"
	replicasetcontroller "minik8s/controller/src/ReplicasetController"
	"minik8s/controller/src/health"
	"minik8s/controller/src/node"
)

var ControllerList []api.Controller = []api.Controller{
	&health.HealthController{},
	&node.NodeController{},
	&PVcontroller.PVcontroller{},
	&replicasetcontroller.ReplicasetController{},
}
