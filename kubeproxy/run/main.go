package main

import "minik8s/kubeproxy/ipvs"
import "minik8s/listwatch"
import "minik8s/global"
import "minik8s/kubeproxy"

func main() {
	ipvs.Init()
	//ipvs.TestConfig()
	var p kubeproxy.ProxyServiceHandler
	var e kubeproxy.ProxyEndpointHandler
	go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
	listwatch.Watch(global.ServiceTopic(), p.HandleService)
	//go listwatch.Watch(global.EndpointTopic(), e.HandleEndpoints)
}