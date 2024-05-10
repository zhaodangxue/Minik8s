package test

import (
	"fmt"
	"minik8s/apiserver/src/apiserver"
	"minik8s/global"
	"minik8s/listwatch"
	"testing"
	"time"
)

var apiServer apiserver.ApiServer

func TestRedis(t *testing.T) {
	fmt.Println("TestRedis")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	fmt.Println("apiServer is running")
	listwatch.Publish(global.TestTopic(), "test-111")
	time.Sleep(3 * time.Second)
}
