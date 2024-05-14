package health

import (
	"minik8s/apiobjects"
	"minik8s/controller/api"
	"minik8s/global"
	"minik8s/utils"
	"time"
)

const HEALTH_REPORT_INTERVAL = 10 * time.Second

type HealthController struct {
	initInfo          api.InitStruct
	ListFuncEnvelops  []api.ListFuncEnvelop
	WatchFuncEnvelops []api.WatchFuncEnvelop
}

func (c *HealthController) Init(init api.InitStruct) {
	c.initInfo = init

	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     HealthReport,
		Interval: HEALTH_REPORT_INTERVAL,
	})

	c.WatchFuncEnvelops = make([]api.WatchFuncEnvelop, 0)
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  c.WatchFunc,
		Topic: global.TestTopic(),
	})
}

func (c *HealthController) GetListFuncEnvelops() []api.ListFuncEnvelop {
	return c.ListFuncEnvelops
}

func (c *HealthController) GetWatchFuncEnvelops() []api.WatchFuncEnvelop {
	return c.WatchFuncEnvelops
}

func HealthReport(controller api.Controller) (err error) {
	utils.Info("HealthReport")
	hc := controller.(*HealthController)
	apiserverUrl := hc.initInfo.GetApiserverUrl()
	utils.Info("HealthReport apiserverUrl:", apiserverUrl)
	// TODO: Send health report to apiserver
	// utils.PutWithJson(hc.initInfo.GetApiserverUrl()+"/health", "health report")
	return
}

func (c *HealthController) WatchFunc(controller api.Controller, message apiobjects.TopicMessage) (err error) {
	utils.Info("HealthController WatchFunc message:", message)
	return
}
