package registry

import (
	"minik8s/controller/api"
	"minik8s/controller/src/health"
)

var ControllerList []api.Controller = []api.Controller{
	&health.HealthController{},
}
