package main

import (
	"minik8s/global"
	"minik8s/kubeproxy"
	"minik8s/kubeproxy/ipvs"
	"minik8s/listwatch"
	"time"
)

var ProxyInterval = 10 * time.Second

func main() {
	ipvs.Init()
	//ipvs.TestConfig()
	var p kubeproxy.ProxyServiceHandler
	var e kubeproxy.ProxyEndpointHandler
	go timedInformer(kubeproxy.CheckAllServiceAndEndpoint, ProxyInterval)
	go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
	listwatch.Watch(global.ServiceTopic(), p.HandleService)
	//go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
}

// 在指定时间间隔内在信道上发送空消息
func timedInformer(Proxfunc func(), interval time.Duration) {
	for {
		Proxfunc()
		time.Sleep(interval)
	}
}