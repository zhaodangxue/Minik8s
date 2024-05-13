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
	PodStatePath       = "/api/podstate"
	PodPathNamespace   = "/api/pod/:namespace"
	PodPathDetail      = "/api/pod/:namespace/:name"
	NodePath           = "/api/node"
	NodePathDetail     = "/api/node/:namespace/:name"
	NodePodBindingPath = "/api/binding/:podnamespace/:podname/:nodename"
	PVPath             = "/api/persistentvolume"
	PVPathNamespace    = "/api/persistentvolume/:namespace"
	PVPathSpecified    = "/api/persistentvolume/:namespace/:name"
	PVPathDetail       = "/api/persistentvolume/:namespace/:name/:storageclass"
	PVCPath            = "/api/persistentvolumeclaim"
	PVCPathNamespace   = "/api/persistentvolumeclaim/:namespace"
	PVCPathSpecified   = "/api/persistentvolumeclaim/:namespace/:name"
	PVCPathDetail      = "/api/persistentvolumeclaim/:namespace/:name/:storageclass"
	PVDynamicAllocate  = "/api/dynamic/allocatePV"
)
