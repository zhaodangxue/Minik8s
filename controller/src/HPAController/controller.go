package hpacontroller

import (
	"context"
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/controller/api"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/redis/go-redis/v9"
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
	workers           map[string]Worker
}

func (c *HPAController) Init(init api.InitStruct) {
	c.initInfo = init
	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     c.Recover,
		Interval: 20 * time.Second,
	})
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
	hpa := apiobjects.HorizontalPodAutoscaler{}
	err := json.Unmarshal([]byte(data), &hpa)
	if err != nil {
		return err
	}
	uid := hpa.ObjectMeta.UID
	_, CancelFuncExists := c.CancelFunctions[uid]
	_, WorkerExists := c.workers[uid]
	if WorkerExists && CancelFuncExists {
		return nil
	}
	ctx, cancel := context.WithCancel(background)
	worker := NewWorker(ctx, &hpa)
	c.workers[uid] = worker
	c.CancelFunctions[uid] = cancel
	utils.Info("HPA", hpa.ObjectMeta.Name, "created")
	go worker.Run()
	return nil
}
func (c *HPAController) HandleHPAUpdate(data string) error {
	hpa := apiobjects.HorizontalPodAutoscaler{}
	err := json.Unmarshal([]byte(data), &hpa)
	if err != nil {
		return err
	}
	uid := hpa.ObjectMeta.UID
	if worker, exist := c.workers[uid]; exist {
		worker.SetTarget(&hpa)
		utils.Info("HPA", hpa.ObjectMeta.Name, "updated")
	}
	return nil
}
func (c *HPAController) HandleHPADelete(data string) error {
	hpa := apiobjects.HorizontalPodAutoscaler{}
	err := json.Unmarshal([]byte(data), &hpa)
	if err != nil {
		return err
	}
	uid := hpa.ObjectMeta.UID
	if cancel, exist := c.CancelFunctions[uid]; exist {
		delete(c.CancelFunctions, uid)
		delete(c.workers, uid)
		cancel()
		utils.Info("HPA", hpa.ObjectMeta.Name, "deleted")
	}
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
func (c *HPAController) Recover(controller api.Controller) error {
	var hpas []*apiobjects.HorizontalPodAutoscaler
	err := utils.GetUnmarshal(route.Prefix+route.HorizontalPodAutoscalerPath, &hpas)
	if err != nil {
		return err
	}
	for _, hpa := range hpas {
		uid := hpa.ObjectMeta.UID
		if _, exist := c.workers[uid]; !exist {
			utils.Info("HPA", hpa.ObjectMeta.Name, "recovered")
			ctx, cancel := context.WithCancel(background)
			worker := NewWorker(ctx, hpa)
			c.workers[uid] = worker
			c.CancelFunctions[uid] = cancel
			go worker.Run()
		}
	}
	return nil
}
func (c *HPAController) HPAHandler(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		return
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		c.HandleHPACreate(topicMessage.Object)
	case apiobjects.Update:
		c.HandleHPAUpdate(topicMessage.Object)
	case apiobjects.Delete:
		c.HandleHPADelete(topicMessage.Object)
	}
}
func (c *HPAController) Run() {
	listwatch.Watch(global.HPARelevantTopic(), c.HPAHandler)
}
func NewController() Controller {
	return &HPAController{
		CancelFunctions: make(map[string]context.CancelFunc),
		workers:         make(map[string]Worker),
	}
}
