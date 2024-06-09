package main

import (
	"minik8s/global"
	"minik8s/jobserver/internal"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"
)

const (
	FullCheckInterval = 15 * time.Second
)

func main(){
	go listwatch.Watch(global.JobTopic(), internal.JobMessageHandler)
	go listwatch.Watch(global.PodStateTopic(), internal.PodStateMessageHandler)
	go utils.CallInterval(internal.FullCheck, FullCheckInterval)
	select {} // block forever
}
