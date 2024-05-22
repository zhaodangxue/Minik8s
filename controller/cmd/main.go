/*
这个包是Controller Manager的启动入口

Controller Manager会初始化所有注册的Controller，然后启动一个goroutine，定期的调用Controller的ListFunc，
同时为每个Controller的WatchFunc注册一个消息监听器，当消息到达时，调用Controller的WatchFunc
*/
package main

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/controller/api"
	"minik8s/controller/src/registry"
	"minik8s/listwatch"
	"minik8s/utils"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func getArgs() string {
	if len(os.Args) != 2 {
		panic("Usage: controller-manager <apiserver-path>")
	}
	return os.Args[1]
}

// topicHandlerGenerator 生成一个消息监听器的处理函数
func topicHandlerGenerator(controller api.Controller, handler api.WatchFunc) listwatch.HandlerOnWatch {
	return func(message *redis.Message) {
		var topicMessage apiobjects.TopicMessage
		err := json.Unmarshal([]byte(message.Payload), &topicMessage)
		if err != nil {
			utils.Error("Failed to unmarshal message payload:", err)
		}
		err = handler(controller, topicMessage)
		if err != nil {
			utils.Error("Failed to handle message:", err)
		}
	}
}

// 为List函数生成一个定期调用的goroutine
func listFuncGenerator(controller api.Controller, listFunc api.ListFunc, interval time.Duration) {
	go func() {
		for {
			err := listFunc(controller)
			if err != nil {
				utils.Error("Err occur when calling list func ", controller, "err:", err)
			}
			utils.Debug("ListFunc done")
			time.Sleep(interval)
		}
	}()
}

func main() {
	apiserverPath := getArgs()

	// 初始化Controller
	for _, controller := range registry.ControllerList {
		controller.Init(api.InitStruct{
			ApiserverPath: apiserverPath,
		})
	}

	// 添加消息监听器
	for _, controller := range registry.ControllerList {
		for _, watchFuncEnvelop := range controller.GetWatchFuncEnvelops() {
			go listwatch.Watch(watchFuncEnvelop.Topic, topicHandlerGenerator(controller, watchFuncEnvelop.Func))
		}
	}

	// 启动定期调用List函数的goroutine
	for _, controller := range registry.ControllerList {
		for _, listFuncEnvelop := range controller.GetListFuncEnvelops() {
			listFuncGenerator(controller, listFuncEnvelop.Func, listFuncEnvelop.Interval)
		}
	}

	// 阻塞主goroutine
	select {}
}
