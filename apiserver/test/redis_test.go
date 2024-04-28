package test

import (
	"fmt"
	"minik8s/apiserver/src/apiserver"
	"minik8s/global"
	"minik8s/listwatch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var apiServer apiserver.ApiServer

func TestRedis(t *testing.T) {
	fmt.Println("TestRedis")
	apiServer = apiserver.New()
	flag := make(chan bool)
	go apiServer.RUN(flag)
	time.Sleep(3 * time.Second)
	a := false
	a = <-flag
	assert.Equal(t, true, a)
	if a {
		fmt.Println("apiServer is running")
		listwatch.Publish(global.TestTopic(), "test-111")
	}
	time.Sleep(3 * time.Second)
}
