package route

import "minik8s/global"

const (
	HttpScheme                           = "http://"
	Hostname                             = global.Host
	Port                                 = ":8080"
	Prefix                               = HttpScheme + Hostname + Port
	TestPostPath                         = "/api/test/post"
	TestGetPath                          = "/api/test/get"
	TestPutPath                          = "/api/test/put/:name/:uid"
	TestDeletePath                       = "/api/test/delete/:name/:uid"
	TestCtlPath                          = "/api/test/ctl"
	PodPath                              = "/api/pod"
	PodStatePath                         = "/api/podstate"
	PodPathNamespace                     = "/api/pod/:namespace"
	PodPathDetail                        = "/api/pod/:namespace/:name"
	NodePath                             = "/api/node"
	NodePathDetail                       = "/api/node/:namespace/:name"
	NodePodBindingPath                   = "/api/binding/:podnamespace/:podname/:nodename"
	ServiceApplyPath                     = "/api/service/apply"
	ServiceCmdDeletePath                 = "/api/service/cmd/delete/:namespace/:name"
	ServiceCreatePath                    = "/api/service"
	NodePortPath                         = "/api/nodeport"
	ServiceUpdatePath                    = "/api/service/update/:namespace/:name"
	ServiceDeletePath                    = "/api/service/delete/:namespace/:name"
	EndpointCreaetPath                   = "/api/endpoint"
	EndpointDeletePath                   = "/api/endpoint/delete/:serviceName/:namespace/:name"
	GetAllPodsPath                       = "/api/get/allpods"
	GetAllServicesPath                   = "/api/get/allservices"
	GetOneServicePath                    = "/api/get/oneservice/:namespace/:name"
	GetAllEndpointsPath                  = "/api/get/allendpoints"
	GetOneEndpointPath                   = "/api/get/oneendpoint/:serviceName/:namespace/:name"
	PVPath                               = "/api/persistentvolume"
	PVPathNamespace                      = "/api/persistentvolume/:namespace"
	PVPathSpecified                      = "/api/persistentvolume/:namespace/:name"
	PVPathDetail                         = "/api/persistentvolume/:namespace/:name/:storageclass"
	PVCPath                              = "/api/persistentvolumeclaim"
	PVCPathNamespace                     = "/api/persistentvolumeclaim/:namespace"
	PVCPathSpecified                     = "/api/persistentvolumeclaim/:namespace/:name"
	PVCPathDetail                        = "/api/persistentvolumeclaim/:namespace/:name/:storageclass"
	NodeHealthPath                       = "/api/nodehealth"
	NodePodBindingAllPath                = "/api/binding"
	NodePodBindingSpecified              = "/api/binding/:namespace/:name"
	PVDynamicAllocate                    = "/api/dynamic/allocatePV"
	ReplicasetPath                       = "/api/replicaset"
	ReplicasetPathNamespace              = "/api/replicaset/:namespace"
	ReplicasetPathSpecified              = "/api/replicaset/:namespace/:name"
	ReplicasetScale                      = "/api/scale/replicaset"
	HorizontalPodAutoscalerPath          = "/api/horizontalpodautoscaler"
	HorizontalPodAutoscalerPathNamespace = "/api/horizontalpodautoscaler/:namespace"
	HorizontalPodAutoscalerPathSpecified = "/api/horizontalpodautoscaler/:namespace/:name"
	SelectPodsByUIDPath                  = "/api/select/:uid"

	DnsApplyPath  = "/api/dns/apply"
	DnsGetAllPath = "/api/dns/get/all"
	DnsDeletePath = "/api/dns/delete/:namespace/:name"

	WorkflowPath          = "/api/workflow"
	WorkflowPathNamespace = "/api/workflow/:namespace"
	WorkflowPathSpecified = "/api/workflow/:namespace/:name"

	FunctionPath          = "/api/function"
	FunctionPathNamespace = "/api/function/:namespace"
	FunctionPathSpecified = "/api/function/:namespace/:name"

	JobPath          = "/api/job"
	JobPathNamespace = "/api/job/:namespace"
	JobPathSpecified = "/api/job/:namespace/:name"

	EventPath          = "/api/event"
	EventPathSpecified = "/api/event/:name"
)
