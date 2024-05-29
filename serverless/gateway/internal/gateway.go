package internal

import (
	"fmt"
	"log"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type ServerlessGateway struct {
	router *gin.Engine
	functions map[string]FunctionWrapper
}

func (a *ServerlessGateway) Init() {
	a.router = gin.Default()
}

func (a *ServerlessGateway) BindHandler() {
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

func (a *ServerlessGateway) Watch() {
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

func (a *ServerlessGateway) List() {
	for _, handler := range ListTable {
		listFuncGenerator(handler.Func, handler.Interval)
	}
}

func (a *ServerlessGateway) RUN() {
	a.router = gin.Default()
	a.BindHandler()
	a.Watch()
	fmt.Println("serverlessGateway is running")
	log.Fatal(a.router.Run(":8081"))
}

// Single Instance
var serverlessGatewayInstance *ServerlessGateway

func GetServerlessGatewayInstance() *ServerlessGateway {
	if serverlessGatewayInstance == nil {
		serverlessGatewayInstance = &ServerlessGateway{}
		serverlessGatewayInstance.Init()
	}
	return serverlessGatewayInstance
}