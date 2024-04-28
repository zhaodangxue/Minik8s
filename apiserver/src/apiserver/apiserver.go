package apiserver

import (
	"log"
	"minik8s/global"
	"minik8s/listwatch"

	"github.com/gin-gonic/gin"
)

type ApiServer interface {
	RUN(flag chan bool)
}

func New() ApiServer {
	return &apiServer{
		router: gin.Default(),
	}
}

type apiServer struct {
	router *gin.Engine
}

func (a *apiServer) BindHandler() {
	for URL, handler := range PostTable {
		a.router.POST(URL, handler)
	}
}
func (a *apiServer) Watch() {
	go listwatch.Watch(global.TestTopic(), SyncTest)
}

//	func ipInit(url, ip string, mask int) error {
//		ig := ipgen.New(url, mask)
//		return ig.ClearIfNotInit(ip)
//	}
func (a *apiServer) RUN(flag chan bool) {
	a.router = gin.Default()
	a.BindHandler()
	a.Watch()
	flag <- true
	log.Fatal(a.router.Run(":8080"))
}
