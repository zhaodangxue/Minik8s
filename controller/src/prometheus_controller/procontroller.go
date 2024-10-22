package prometheuscontroller

import (
	"fmt"
	"minik8s/apiobjects"
	"minik8s/prometheus"
	"time"

	"minik8s/apiserver/src/route"
	"minik8s/controller/api"
	"minik8s/utils"
	"strconv"
	"reflect"

	log "github.com/sirupsen/logrus"
)


var podList = map[string]string{}
var nodeList = map[string]*apiobjects.Node{}


const ProController_REPORT_INTERVAL = 15 * time.Second

type PrometheusController struct {
	initInfo          api.InitStruct
	ListFuncEnvelops  []api.ListFuncEnvelop
	WatchFuncEnvelops []api.WatchFuncEnvelop
}

func (c *PrometheusController) Init(init api.InitStruct) {
	c.initInfo = init
	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     CheckAllNodeAndPod,
		Interval: ProController_REPORT_INTERVAL,
	})

	c.WatchFuncEnvelops = make([]api.WatchFuncEnvelop, 0)
	// c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
	// 	Func:  c.se.HandleEndpoints,
	// 	Topic: global.PodStateTopic(),
	// })
}

func (c *PrometheusController) GetListFuncEnvelops() []api.ListFuncEnvelop {
	return c.ListFuncEnvelops
}
func (c *PrometheusController) GetWatchFuncEnvelops() []api.WatchFuncEnvelop {
	return c.WatchFuncEnvelops
}

func CheckAllNodeAndPod(controller api.Controller)(error) {
	utils.Info("[Prometheus Controller] CheckAllNodeAndPod")
   tmp_nodeList := make(map[string]*apiobjects.Node)
   tmp_podList := make(map[string]string)

	utils.Info("CheckAllService")
	pod_list := []*apiobjects.Pod{}
	err := utils.GetUnmarshal(route.Prefix +"/api/get/allpods", &pod_list)
	if err != nil {
		fmt.Println("get pod list error")
	}
	node_list := []*apiobjects.Node{}
	err = utils.GetUnmarshal("http://localhost:8080/api/node", &node_list)
	if err != nil {
		fmt.Println("get svc list error")
	}

	for _, node := range node_list {
		//utils.Info("[Prometheus Controller] find Node")
		if node.Status.State == apiobjects.NodeStateHealthy {
			tmp_nodeList[node.Info.Ip] = node
		}
	}

	for _, pod := range pod_list{
		//utils.Info("[Prometheus Controller] get pod")
		if pod.Status.PodPhase == apiobjects.PodPhase_POD_RUNNING {
			flag := false
			for key, value := range pod.Labels {
				if key == "log" && value == "prometheus" {
					flag = true
					break
				}
			}
			if flag == false {
				continue
			}
			dstPort := findDstPort("prometheus", pod.Spec.Containers)
			if dstPort == 1314 {
				log.Fatal("[pro controller] No Match Port for Prometheus!")
				continue
			}
			tmp_podList[pod.Status.PodIP] = strconv.Itoa(int(dstPort))
		}
	}
	//打印一下tmp_nodeList和tmp_podList
	// for key, value := range nodeList {
	// 	utils.Info("[Prometheus Controller] nodeList: ", key, value)
	// }
	// for key, value := range podList {
	// 	utils.Info("[Prometheus Controller] podList: ", key, value)
	// }
	// for key, value := range tmp_nodeList {
	// 	utils.Info("[Prometheus Controller] tmp_nodeList: ", key, value)
	// }
	// for key, value := range tmp_podList {
	// 	utils.Info("[Prometheus Controller] tmp_podList: ", key, value)
	// }
    // isNodeUpdate := false
	// for key, _ := range tmp_nodeList {
	// 	_, exist := nodeList[key]
	// 	if !exist {
	// 		isNodeUpdate = true
	// 		break
	// 	}
	// }

	// isPodUpdate := false
	// for key, v1 := range tmp_podList {
	// 	v2, exist := podList[key]
	// 	if !exist {
	// 		isPodUpdate = true
	// 		break
	// 	}
	// 	if v1 != v2 {
	// 		isPodUpdate = true
	// 		break
	// 	}
	// }
	isNodeUpdate := !reflect.DeepEqual(tmp_nodeList, nodeList)
	isPodUpdate := !reflect.DeepEqual(tmp_podList, podList)
	if isPodUpdate {
		utils.Info("[Prometheus Controller] pod update")
	}

	if isNodeUpdate || isPodUpdate {
		// update prometheus config
		utils.Info("[Prometheus Controller] update config")
		podList = tmp_podList
		nodeList = tmp_nodeList
		
		configs := []string{}
		for key, _ := range nodeList {
			configs = append(configs, key+":9100")
		}
		for key, value := range podList {
			//utils.Info("[Prometheus Controller] get pod")
			configs = append(configs, key+":"+value)
		}
		prometheus.GenerateProConfig(configs)
		err := prometheus.ReloadPrometheus()
		if err != nil {
			utils.Error("reload prometheus error")
			return err
		}
	}
	return nil
}

func findDstPort(targetPort string, containers []apiobjects.Container) int32 {
	for _, c := range containers {
		for _, p := range c.Ports {
			if p.Name == targetPort {
				return p.ContainerPort
			}
		}
	}
	log.Fatal("[pro controller] No Match Port for Prometheus!")
	return 1314
}


func IsMapEqual(a map[string]string, b map[string]string) bool {
	for k, v := range a {
		_,exist := b[k]
		if !exist {
			utils.Info("[Prometheus Controller] map not equal")
			return false
		}
		if b[k] != v {
			utils.Info("[Prometheus Controller] map not equal")
			return false
		}
	}
	return true
}