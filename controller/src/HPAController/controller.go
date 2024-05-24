package hpacontroller

import (
	"context"
	"minik8s/apiobjects"
	"minik8s/controller/api"
	"minik8s/global"
)

var background = context.Background()

type Controller interface {
	Run()
}
type HPAController struct {
	initInfo          api.InitStruct
	ListFuncEnvelops  []api.ListFuncEnvelop
	WatchFuncEnvelops []api.WatchFuncEnvelop
	CancelFunctions   map[string]context.CancelFunc
}

func (c *HPAController) Init(init api.InitStruct) {
	c.initInfo = init
	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.WatchFuncEnvelops = make([]api.WatchFuncEnvelop, 0)
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  c.WatchHPA,
		Topic: global.HPARelevantTopic(),
	})
}
func (c *HPAController) GetListFuncEnvelops() []api.ListFuncEnvelop {
	return c.ListFuncEnvelops
}
func (c *HPAController) GetWatchFuncEnvelops() []api.WatchFuncEnvelop {
	return c.WatchFuncEnvelops
}
func (c *HPAController) HandleHPACreate(data string) error {
	return nil
}
func (c *HPAController) HandleHPAUpdate(data string) error {
	return nil
}
func (c *HPAController) HandleHPADelete(data string) error {
	return nil
}
func (c *HPAController) WatchHPA(controller api.Controller, message apiobjects.TopicMessage) error {
	var err error
	switch message.ActionType {
	case apiobjects.Create:
		err = c.HandleHPACreate(message.Object)
	case apiobjects.Update:
		err = c.HandleHPAUpdate(message.Object)
	case apiobjects.Delete:
		err = c.HandleHPADelete(message.Object)
	}
	return err
}
