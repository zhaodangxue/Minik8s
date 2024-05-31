package apiserver

import (
	"minik8s/apiserver/src/handler"
	"minik8s/apiserver/src/handler/event"
	function_handler "minik8s/apiserver/src/handler/function"
	hpa_handler "minik8s/apiserver/src/handler/hpa"
	node "minik8s/apiserver/src/handler/node"
	persistentvolume_handler "minik8s/apiserver/src/handler/persistentvolume"
	replicaset_handler "minik8s/apiserver/src/handler/replicaset"
	workflow_handler "minik8s/apiserver/src/handler/workflow"
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
	// route.EndpointUpdatePath: handler.EndpointUpdateHandler,
	route.PodPath:                     handler.PodApplyHandler,
	route.PVPath:                      handler.PVApplyHandler,
	route.PVCPath:                     handler.PVCApplyHandler,
	route.PVCPathDetail:               persistentvolume_handler.PVCApplyDetailHandler,
	route.PVPathDetail:                persistentvolume_handler.PVApplyDetailHandler,
	route.PVDynamicAllocate:           persistentvolume_handler.PVDynamicAllocateHandler,
	route.ReplicasetPath:              replicaset_handler.ReplicasetApplyHandler,
	route.DnsApplyPath:                handler.DnsApplyHandler,
	route.HorizontalPodAutoscalerPath: hpa_handler.HPAApplyHandler,
	route.WorkflowPath:                workflow_handler.WorkflowApplyHandler,
	route.FunctionPath:                function_handler.FunctionApplyHandler,
	route.EventPath:                   event.EventCreateHandler,
}

var GetTable = map[string]HandlerFunc{
	route.TestGetPath:         handler.TestGetHandler,
	route.TestCtlPath:         handler.TestCtlGetHandler,
	route.NodePath:            handler.NodeGetHandler,
	route.PodPath:             handler.PodGetHandler,
	route.PodPathNamespace:    handler.PodGetWithNamespaceHandler,
	route.PodPathDetail:       handler.PodGetDetailHandler,
	route.GetAllPodsPath:      handler.GetAllPodsHandler,
	route.GetOneServicePath:   handler.GetOneServiceHandler,
	route.GetAllServicesPath:  handler.GetAllServicesHandler,
	route.GetOneEndpointPath:  handler.GetOneEndpointHandler,
	route.GetAllEndpointsPath: handler.GetAllEndpointsHandler,
	route.PVPathNamespace:     handler.PVGetWithNamespaceHandler,
	route.PVCPathNamespace:    handler.PVCGetWithNamespaceHandler,

	route.NodePodBindingAllPath:                handler.NodePodBindingAllHandler,
	route.NodePodBindingSpecified:              handler.NodePodBindingSpecifiedHandler,
	route.PVPathSpecified:                      persistentvolume_handler.PVGetSpecifiedHandler,
	route.PVCPathSpecified:                     persistentvolume_handler.PVCGetSpecifiedHandler,
	route.PVCPath:                              persistentvolume_handler.PVCSGetHandler,
	route.ReplicasetPath:                       replicaset_handler.ReplicasetGetHandler,
	route.ReplicasetPathNamespace:              replicaset_handler.ReplicasetGetWithNamespaceHandler,
	route.ReplicasetPathSpecified:              replicaset_handler.ReplicasetGetSpecifiedHandler,
	route.DnsGetAllPath:                        handler.DnsGetAllHandler,
	route.HorizontalPodAutoscalerPathNamespace: hpa_handler.HPAGetWithNamespaceHandler,
	route.HorizontalPodAutoscalerPath:          hpa_handler.HPAGetHandler,
	route.FunctionPath:                         function_handler.FunctionGetHandler,
}

var PutTable = map[string]HandlerFunc{
	route.TestPutPath:                          handler.TestPutHandler,
	route.NodeHealthPath:                       handler.NodeHealthHandler,
	route.ReplicasetPathSpecified:              replicaset_handler.ReplicasetUpdateHandler,
	route.ReplicasetScale:                      replicaset_handler.ReplicasetScaleHandler,
	route.HorizontalPodAutoscalerPathSpecified: hpa_handler.HPAUpdateHandler,
}

var DeleteTable = map[string]HandlerFunc{
	route.TestDeletePath:                       handler.TestDeleteHandler,
	route.TestCtlPath:                          handler.TestCtlDeleteHandler,
	route.ServiceDeletePath:                    handler.ServiceDeleteHandler,
	route.ServiceCmdDeletePath:                 handler.ServiceCmdDeleteHandler,
	route.EndpointDeletePath:                   handler.EndpointDeleteHandler,
	route.PodPathDetail:                        handler.PodDeleteHandler,
	route.PVCPathSpecified:                     handler.PVCDeleteHandler,
	route.PVPathSpecified:                      handler.PVDeleteHandler,
	route.NodePathDetail:                       node.NodeDeleteHandler,
	route.ReplicasetPathSpecified:              replicaset_handler.ReplicasetDeleteHandler,
	route.DnsDeletePath:                        handler.DnsDeleteHandler,
	route.HorizontalPodAutoscalerPathSpecified: hpa_handler.HPADeleteHandler,
	route.WorkflowPathSpecified:                workflow_handler.WorkflowDeleteHandler,
	route.FunctionPathSpecified:                function_handler.FunctionDeleteHandler,
	route.EventPathSpecified:                   event.EventDeleteHandler,
}
