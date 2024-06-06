package main

import (
	"minik8s/global"
	"minik8s/kubeproxy"
	"minik8s/kubeproxy/ipvs"
	"minik8s/listwatch"
	"minik8s/utils"
	"os"
	"time"
)

var ProxyInterval = 10 * time.Second

func main() {
	kubeproxy.ServerUrl = getMasterPath()

	ipvs.Init()
	//ipvs.TestConfig()
	var p kubeproxy.ProxyServiceHandler
	var e kubeproxy.ProxyEndpointHandler
	go timedInformer(kubeproxy.CheckAllServiceAndEndpoint, ProxyInterval)
	go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
	listwatch.Watch(global.ServiceTopic(), p.HandleService)
	//go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
}

func getMasterPath() string {
	// 从参数中获取服务器url
	if (len(os.Args) != 2) {
		utils.Error("Usage: ./kubeproxy <server-url>")
		utils.Error("Example: ./kubeproxy 192.168.1.12:8080")
		os.Exit(1)
	}
	serverUrl := "http://" + os.Args[1]
	return serverUrl
}

// 在指定时间间隔内在信道上发送空消息
func timedInformer(Proxfunc func(), interval time.Duration) {
	for {
		Proxfunc()
		time.Sleep(interval)
	}
}