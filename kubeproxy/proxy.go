package kubeproxy

/* 主要工作：
1. 监听service资源的创建。创建service
2. 监听service资源的删除。删除service
3. 监听endpoint的创建。设置dest规则。
4. 监听endpoint的删除。删除对应dest规则。
*/

import (
	"minik8s/kubeproxy/ipvs"
	"minik8s/apiobjects"
	"strconv"
)

func Run() {
	ipvs.Init()
	//ipvs.TestConfig()
	//var p proxyServiceHandler
	//var e proxyEndpointHandler

}

/* ========== Start Service Handler ========== */

type proxyServiceHandler struct {
}

func (p proxyServiceHandler) HandleCreate(message []byte) {

}

func (p proxyServiceHandler) HandleDelete(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	for _, p := range svc.Spec.Ports {
		key := svc.Status.ClusterIP + ":" + strconv.Itoa(int(p.Port))
		ipvs.DeleteService(key)
	}

}

func (p proxyServiceHandler) HandleUpdate(message []byte) {
	svc := &apiobjects.Service{}
	svc.UnMarshalJSON(message)

	for _, p := range svc.Spec.Ports {
		ipvs.AddService(svc.Status.ClusterIP, uint16(p.Port))
	}

}

func (p proxyServiceHandler) GetType() string {
	return "service"
}

/* ========== Start Endpoint Handler ========== */

type proxyEndpointHandler struct {
}

func (e proxyEndpointHandler) HandleCreate(message []byte) {
	edpt := &apiobjects.Endpoint{}
	edpt.UnMarshalJSON(message)

	key := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
	ipvs.AddEndpoint(key, edpt.Spec.DestIP, uint16(edpt.Spec.DestPort))
}

func (e proxyEndpointHandler) HandleDelete(message []byte) {
	edpt := &apiobjects.Endpoint{}
	edpt.UnMarshalJSON(message)

	svcKey := edpt.Spec.SvcIP + ":" + strconv.Itoa(int(edpt.Spec.SvcPort))
	dstKey := edpt.Spec.DestIP + ":" + strconv.Itoa(int(edpt.Spec.DestPort))
	ipvs.DeleteEndpoint(svcKey, dstKey)
}

func (e proxyEndpointHandler) HandleUpdate(message []byte) {

}

func (e proxyEndpointHandler) GetType() string {
	return "endpoint"
}