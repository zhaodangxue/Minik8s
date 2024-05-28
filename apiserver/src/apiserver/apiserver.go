package apiserver

import (
	"fmt"
	"log"
	"minik8s/global"
	"minik8s/listwatch"

	"github.com/gin-gonic/gin"
)

type ApiServer interface {
	RUN()
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
	for URL, handler := range GetTable {
		a.router.GET(URL, handler)
	}
	for URL, handler := range PutTable {
		a.router.PUT(URL, handler)
	}
	for URL, handler := range DeleteTable {
		a.router.DELETE(URL, handler)
	}
}
func (a *apiServer) Watch() {
	go listwatch.Watch(global.TestTopic(), SyncTest)
}

//	func ipInit(url, ip string, mask int) error {
//		ig := ipgen.New(url, mask)
//		return ig.ClearIfNotInit(ip)
//	}
func (a *apiServer) RUN() {
	a.router = gin.Default()
	a.BindHandler()
	a.Watch()
	fmt.Println("apiServer is running")
	log.Fatal(a.router.Run(":8080"))
}
