package ipvs

import (
	"fmt"
	"os/exec"
	"strconv"
	"github.com/mqliang/libipvs"
	log "github.com/sirupsen/logrus"
	"net"
	"syscall"
	"time"
	//"github.com/coreos/go-iptables/iptables"
)

 var handler libipvs.IPVSHandle

func Init() {
	h, err := libipvs.New()
	handler = h
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = exec.Command("sysctl", []string{"net.ipv4.vs.conntrack=1"}...).CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
		print("sysctl error")
	}
}

func AddService(ip string, port uint16) {
	serviceIP := ip + ":" + strconv.Itoa(int(port))
	if _, ok := Services[serviceIP]; ok {
		return
	}
	svc := addService(ip, port)
	Services[serviceIP] = &ServiceNode{
		Service:   svc,
		Visited:   true,
		Endpoints: map[string]*EndpointNode{},
	}
	log.Info("[kubeproxy] Add service ", serviceIP)
}

func addService(ip string, port uint16) *libipvs.Service {
	//创建一个service结构体并将其添加到ipvs。
	// 等价于命令 ipvsadm -A -t 10.10.0.1:8410 -s rr
	svc := &libipvs.Service{
		Address:       net.ParseIP(ip),
		AddressFamily: syscall.AF_INET,
		Protocol:      libipvs.Protocol(syscall.IPPROTO_TCP),
		Port:          port,
		SchedName:     libipvs.RoundRobin,
	}
    print(svc.Address.String() + ":" + strconv.Itoa(int(svc.Port)))

	if err := handler.NewService(svc); err != nil {
		fmt.Println(err.Error())
	}
	//绑定ip地址到flannel.1网卡上
	// 等价于命令ip addr add 10.10.0.1/24 dev flannel.1
	args := []string{"addr", "add", ip + "/24", "dev", "flannel.1"}
	_, err := exec.Command("ip", args...).CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}

	//配置iptables:添加SNAT规则
	// 等价于命令iptables -t nat -A POSTROUTING -m ipvs  --vaddr 10.9.0.1 --vport 12 -j MASQUERADE
	args = []string{"-t", "nat", "-A", "POSTROUTING", "-m", "ipvs", "--vaddr", ip, "--vport", strconv.Itoa(int(svc.Port)), "-j", "MASQUERADE"}
	_, err = exec.Command("iptables", args...).CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}

	return svc
}

func DeleteService(key string) {
	log.Info("[kubeproxy] Delete service ", key)
	node := Services[key]
	if node != nil {
		deleteService(node.Service)
	}
	delete(Services, key)
}

func deleteService(svc *libipvs.Service) {
	if err := handler.DelService(svc); err != nil {
		fmt.Println(err.Error())
	}
}

func AddEndpoint(key string, ip string, port uint16) {
	svc, exist := Services[key]
	for !exist {
		time.Sleep(1)
		//to do 可能会死循环
		log.Info("[proxy] Add Endpoint: service doesn't exist!")
		svc, exist = Services[key]
	}
	dst := bindEndpoint(svc.Service, ip, port)
	podIP := ip + ":" + strconv.Itoa(int(port))
	svc.Endpoints[podIP] = &EndpointNode{
		Endpoint: dst,
		Visited:  true,
	}
	log.Info("[kubeproxy] Add endpoint ", podIP, " service:", key)
}

func bindEndpoint(svc *libipvs.Service, ip string, port uint16) *libipvs.Destination {
	dst := libipvs.Destination{
		Address:       net.ParseIP(ip),
		AddressFamily: syscall.AF_INET,
		Port:          port,
	}

	//print(svc.Address.String() + ":" + strconv.Itoa(int(svc.Port)))
	//等价于命令ipvsadm -a -t 10.10.0.1:8410(服务的IP和端口) -r 10.6.4.1:1234(endpoint的IP和端口) -m(使用MASQUERADING模式)                     
	args := []string{"-a", "-t", svc.Address.String() + ":" + strconv.Itoa(int(svc.Port)), "-r", ip + ":" + strconv.Itoa(int(port)), "-m"}
	_, err := exec.Command("ipvsadm", args...).CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}

	return &dst
}

func DeleteEndpoint(svcKey string, dstKey string) {
	if svc, ok := Services[svcKey]; ok {
		dst := svc.Endpoints[dstKey].Endpoint
		unbindEndpoint(svc.Service, dst)
		delete(svc.Endpoints, dstKey)
	}
	log.Info("[kubeproxy] Delete endpoint ", dstKey, " service:", svcKey)
}

func unbindEndpoint(svc *libipvs.Service, dst *libipvs.Destination) {
	if err := handler.DelDestination(svc, dst); err != nil {
		fmt.Println(err.Error())
	}
}

