package main

import (
	"minik8s/global"
	"minik8s/jobserver/internal"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	FullCheckInterval = 15 * time.Second
)

func main() {
	go listwatch.Watch(global.JobTopic(), internal.JobMessageHandler)
	go listwatch.Watch(global.PodStateTopic(), internal.PodStateMessageHandler)
	go utils.CallInterval(internal.ClusterFullCheck, FullCheckInterval)

	router := gin.Default()
	router.POST("/job/run/:namespace/:name", internal.JobRunHandler)
	router.GET("/job/status/:namespace/:name", internal.JobGetStatusHandler)
	utils.Fatal(router.Run(":8082"))
	select {} // block forever
}
