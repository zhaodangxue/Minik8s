package registry

import (
	"minik8s/controller/api"
	"minik8s/controller/src/health"
	service "minik8s/controller/src/service_controller"
)

var ControllerList []api.Controller = []api.Controller{
	&health.HealthController{},
	&service.ServiceController{},
}
