package ipvs

import "github.com/mqliang/libipvs"

type ServiceNode struct {
	Service   *libipvs.Service
	Endpoints map[string]*EndpointNode
	Visited   bool
}
type EndpointNode struct {
	Endpoint *libipvs.Destination
	//endpoints []*libipvs.Destination
	Visited bool
}

var Services = make(map[string]*ServiceNode)
//var Nodeports = make(map[string]*ServiceNode)

//var Endpoints map[string]EndpointNode