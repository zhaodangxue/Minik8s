package route

const (
	HttpScheme         = "http://"
	Hostname           = "localhost"
	Port               = ":8080"
	Prefix             = HttpScheme + Hostname + Port
	TestPostPath       = "/api/test/post"
	TestGetPath        = "/api/test/get"
	TestPutPath        = "/api/test/put/:name/:uid"
	TestDeletePath     = "/api/test/delete/:name/:uid"
	TestCtlPath        = "/api/test/ctl"
	PodPath            = "/api/pod"
	PodPathNamespace   = "/api/pod/:namespace"
	PodPathDetail      = "/api/pod/:namespace/:name"
	NodePath           = "/api/node"
	NodePathDetail     = "/api/node/:namespace/:name"
	NodePodBindingPath = "/api/binding/:podnamespace/:podname/:nodename"

	ServiceApplyPath   = "/api/service/apply"
	ServiceCreatePath  = "/api/service"
	ServiceUpdatePath  = "/api/service/update/:namespace/:name"
	ServiceDeletePath  = "/api/service/delete/:namespace/:name"
	EndpointCreaetPath = "/api/endpoint"
	EndpointDeletePath = "/api/endpoint/delete/:serviceName/:namespace/:name"

	GetAllPodsPath = "/api/get/allpods"
)
