package route

const (
	HttpScheme     = "http://"
	Hostname       = "localhost"
	Port           = ":8080"
	Prefix         = HttpScheme + Hostname + Port
	TestPostPath   = "/api/test/post"
	TestGetPath    = "/api/test/get"
	TestPutPath    = "/api/test/put/:name/:uid"
	TestDeletePath = "/api/test/delete/:name/:uid"
	TestCtlPath    = "/api/test/ctl"
	PodPath        = "/api/pod"
)
