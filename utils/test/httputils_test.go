package test_utils

import (
	"fmt"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/handler"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var apiServer apiserver.ApiServer

func TestPost(t *testing.T) {
	fmt.Println("TestPost")
	etcd.Clear()
	apiServer = apiserver.New()
	flag := make(chan bool)
	go apiServer.RUN(flag)
	time.Sleep(3 * time.Second)
	a := false
	a = <-flag
	assert.Equal(t, true, a)
	obj := handler.TestBody{
		Body:    "test",
		Name:    "test",
		Uid:     "100",
		Num:     1,
		Percent: 1.0,
	}
	response, err := utils.PostWithJson(route.Prefix+route.TestPostPath, obj)
	fmt.Println(response)
	assert.Nil(t, err)
}
func TestGet(t *testing.T) {
	fmt.Println("TestGet")
	etcd.Clear()
	apiServer = apiserver.New()
	flag := make(chan bool)
	go apiServer.RUN(flag)
	time.Sleep(3 * time.Second)
	a := false
	a = <-flag
	assert.Equal(t, true, a)
	etcd.Put("/test/get", `{"body":"test","name":"test","uid":"100","num":1,"percent":1}`)
	var TestJson handler.TestBody
	err := utils.GetUnmarshal(route.Prefix+route.TestGetPath, &TestJson)
	fmt.Printf("%+v\n", TestJson)
	assert.Nil(t, err)
}
func TestPut(t *testing.T) {
	fmt.Println("TestPut")
	etcd.Clear()
	apiServer = apiserver.New()
	flag := make(chan bool)
	go apiServer.RUN(flag)
	time.Sleep(3 * time.Second)
	a := false
	a = <-flag
	assert.Equal(t, true, a)
	etcd.Put("/test/put/zbm/100", `{"body":"test","name":"zbm","uid":"100","num":1,"percent":1}`)
	obj := handler.TestBody{
		Body:    "test2",
		Name:    "zbm",
		Uid:     "100",
		Num:     2,
		Percent: 2.0,
	}
	url := route.Prefix + "/api/test/put" + "/zbm/100"
	response, err := utils.PutWithJson(url, obj)
	fmt.Println(response)
	assert.Nil(t, err)
}
func TestDelete(t *testing.T) {
	fmt.Println("TestDelete")
	etcd.Clear()
	apiServer = apiserver.New()
	flag := make(chan bool)
	go apiServer.RUN(flag)
	time.Sleep(3 * time.Second)
	a := false
	a = <-flag
	assert.Equal(t, true, a)
	etcd.Put("/test/delete/zbm/100", `{"body":"test","name":"zbm","uid":"100","num":1,"percent":1}`)
	response, err := utils.Delete(route.Prefix + "/api/test/delete" + "/zbm/100")
	fmt.Println(response)
	assert.Nil(t, err)
}
