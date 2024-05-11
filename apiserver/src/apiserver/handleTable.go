package apiserver

import (
	"minik8s/apiserver/src/handler"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc

var PostTable = map[string]HandlerFunc{
	route.TestPostPath:       handler.TestPostHandler,
	route.TestCtlPath:        handler.TestCtlHandler,
	route.NodePodBindingPath: handler.NodePodBindingHandler,
	route.ServiceCreatePath:  handler.ServiceCreateHandler,
	route.ServiceUpdatePath:  handler.ServiceUpdateHandler,
	route.EndpointCreaetPath: handler.EndpointCreateHandler,
	route.ServiceApplyPath:   handler.ServiceApplyHandler,
	//route.EndpointUpdatePath: handler.EndpointUpdateHandler,
	route.PodPath: handler.PodApplyHandler,
	route.PVPath:  handler.PVApplyHandler,
	route.PVCPath: handler.PVCApplyHandler,
}
var GetTable = map[string]HandlerFunc{
	route.TestGetPath:      handler.TestGetHandler,
	route.TestCtlPath:      handler.TestCtlGetHandler,
	route.NodePath:         handler.NodeGetHandler,
	route.PodPathNamespace: handler.PodGetWithNamespaceHandler,
	route.PodPathDetail:    handler.PodGetDetailHandler,
	route.GetAllPodsPath:   handler.GetAllPodsHandler,
	route.PVPathNamespace:  handler.PVGetWithNamespaceHandler,
	route.PVCPathNamespace: handler.PVCGetWithNamespaceHandler,
}
var PutTable = map[string]HandlerFunc{
	route.TestPutPath: handler.TestPutHandler,
}
var DeleteTable = map[string]HandlerFunc{
	route.TestDeletePath:     handler.TestDeleteHandler,
	route.TestCtlPath:        handler.TestCtlDeleteHandler,
	route.ServiceDeletePath:  handler.ServiceDeleteHandler,
	route.EndpointDeletePath: handler.EndpointDeleteHandler,
	//route.TestDeletePath:     handler.TestDeleteHandler,
	//route.TestCtlPath:        handler.TestCtlDeleteHandler,
	//route.PodPathDetail:      handler.PodDeleteHandler,
	route.PVCPathSpecified: handler.PVCDeleteHandler,
	route.PVPathSpecified:  handler.PVDeleteHandler,
}
