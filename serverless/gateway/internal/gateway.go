package internal

import (
	"fmt"
	"log"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerlessGateway interface {
	RUN()
}

func New() ServerlessGateway {
	return &serverlessGateway{
		router: gin.Default(),
	}
}

type serverlessGateway struct {
	router *gin.Engine
}

func (a *serverlessGateway) BindHandler() {
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

func (a *serverlessGateway) Watch() {
	for topic, handler := range WatchTable {
		go listwatch.Watch(topic, handler)
	}
}

func listFuncGenerator(listFunc ListFunc, interval time.Duration) {
	go func() {
		for {
			err := listFunc()
			if err != nil {
				utils.Error("Err occur when calling list func err: ", err)
			}
			utils.Debug("ListFunc done")
			time.Sleep(interval)
		}
	}()
}

func (a *serverlessGateway) List() {
	for _, handler := range ListTable {
		listFuncGenerator(handler.Func, handler.Interval)
	}
}

func (a *serverlessGateway) RUN() {
	a.router = gin.Default()
	a.BindHandler()
	a.Watch()
	fmt.Println("serverlessGateway is running")
	log.Fatal(a.router.Run(":8081"))
}
