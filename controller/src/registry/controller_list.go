package registry

import (
	"minik8s/controller/api"
	"minik8s/controller/src/PVcontroller"
	"minik8s/controller/src/health"
	"minik8s/controller/src/node"
)

var ControllerList []api.Controller = []api.Controller{
	&health.HealthController{},
	&node.NodeController{},
	&PVcontroller.PVcontroller{},
}
