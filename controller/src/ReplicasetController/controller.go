package replicasetcontroller

import (
	"encoding/json"
	"time"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/controller/api"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"

	"github.com/redis/go-redis/v9"
)

type Controller interface {
	Run()
}
type ReplicasetController struct {
	initInfo          api.InitStruct
	ListFuncEnvelops  []api.ListFuncEnvelop
	WatchFuncEnvelops []api.WatchFuncEnvelop
	Workers           map[string]Worker
}

func (c *ReplicasetController) Init(init api.InitStruct) {
	c.initInfo = init
	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     c.Sync,
		Interval: 20 * time.Second,
	})
	c.ListFuncEnvelops = append(c.ListFuncEnvelops, api.ListFuncEnvelop{
		Func:     c.Recover,
		Interval: 20 * time.Second,
	})
	c.WatchFuncEnvelops = make([]api.WatchFuncEnvelop, 0)
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  c.WatchReplicaset,
		Topic: global.ReplicasetTopic(),
	})
	c.Workers = make(map[string]Worker)
}

func (c *ReplicasetController) HandleReplicasetCreate(data string) error {
	replicaset := apiobjects.Replicaset{}
	err := json.Unmarshal([]byte(data), &replicaset)
	if err != nil {
		return err
	}
	uid := replicaset.ObjectMeta.UID
	if _, exist := c.Workers[uid]; !exist {
		utils.Info("Replicaset", replicaset.ObjectMeta.Name, "created")
		worker := NewWorker(&replicaset)
		c.Workers[uid] = worker
		go worker.Run()
	}
	return nil
}

func (c *ReplicasetController) HandleReplicasetUpdate(data string) error {
	replicaset := apiobjects.Replicaset{}
	err := json.Unmarshal([]byte(data), &replicaset)
	if err != nil {
		return err
	}
	uid := replicaset.ObjectMeta.UID
	if worker, exist := c.Workers[uid]; exist {
		// TO DO update replicaset
		utils.Info("Replicaset", replicaset.ObjectMeta.Name, "updated")
		worker.ResetTarget(&replicaset)
		worker.SyncCh() <- struct{}{}
	}
	return nil
}
func (c *ReplicasetController) HandleReplicasetScale(data string) error {
	replicaset := apiobjects.Replicaset{}
	err := json.Unmarshal([]byte(data), &replicaset)
	if err != nil {
		return err
	}
	uid := replicaset.ObjectMeta.UID
	if worker, exist := c.Workers[uid]; exist {
		// TO DO scale replicaset
		utils.Info("Replicaset", replicaset.ObjectMeta.Name, "scaled")
		worker.ScaleTarget(&replicaset)
		worker.SyncCh() <- struct{}{}
	}
	return nil
}

func (c *ReplicasetController) HandleReplicasetDelete(data string) error {
	replicaset := apiobjects.Replicaset{}
	err := json.Unmarshal([]byte(data), &replicaset)
	if err != nil {
		return err
	}
	uid := replicaset.ObjectMeta.UID
	if worker, exist := c.Workers[uid]; exist {
		// TO DO delete replicaset
		close(worker.SyncCh())
		worker.Done()
		delete(c.Workers, uid)
		utils.Info("Replicaset", replicaset.ObjectMeta.Name, "deleted")
	}
	return nil
}

func (c *ReplicasetController) GetListFuncEnvelops() []api.ListFuncEnvelop {
	return c.ListFuncEnvelops
}

func (c *ReplicasetController) GetWatchFuncEnvelops() []api.WatchFuncEnvelop {
	return c.WatchFuncEnvelops
}

func (c *ReplicasetController) Sync(controller api.Controller) error {
	var pods []*apiobjects.Pod
	utils.Info("ReplicasetController Sync")
	err := utils.GetUnmarshal(route.Prefix+route.PodPath, &pods)
	for _, worker := range c.Workers {
		worker.SetPods(pods)
	}
	return err
}

func (c *ReplicasetController) Recover(controller api.Controller) error {
	var replicasets []*apiobjects.Replicaset
	err := utils.GetUnmarshal(route.Prefix+route.ReplicasetPath, &replicasets)
	if err != nil {
		return err
	}
	for _, replicaset := range replicasets {
		uid := replicaset.ObjectMeta.UID
		if _, exist := c.Workers[uid]; !exist {
			utils.Info("Replicaset", replicaset.ObjectMeta.Name, "reocvered")
			worker := NewWorker(replicaset)
			c.Workers[uid] = worker
			go worker.Run()
		}
	}
	return nil
}

func (c *ReplicasetController) WatchReplicaset(controller api.Controller, message apiobjects.TopicMessage) error {
	var err error
	switch message.ActionType {
	case apiobjects.Create:
		err = c.HandleReplicasetCreate(message.Object)
	case apiobjects.Update:
		err = c.HandleReplicasetUpdate(message.Object)
	case apiobjects.Delete:
		err = c.HandleReplicasetDelete(message.Object)
	case apiobjects.Scale:
		err = c.HandleReplicasetScale(message.Object)
	}
	return err
}

func (c *ReplicasetController) ReplicasetHandler(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		return
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		c.HandleReplicasetCreate(topicMessage.Object)
	case apiobjects.Update:
		c.HandleReplicasetUpdate(topicMessage.Object)
	case apiobjects.Delete:
		c.HandleReplicasetDelete(topicMessage.Object)
	case apiobjects.Scale:
		c.HandleReplicasetScale(topicMessage.Object)
	}
}

func (c *ReplicasetController) SyncTest() {
	go func() {
		for {
			var pods []*apiobjects.Pod
			err := utils.GetUnmarshal(route.Prefix+route.PodPath, &pods)
			if err != nil {
				utils.Error("Failed to get pods:", err)
			}
			for _, worker := range c.Workers {
				worker.SetPods(pods)
			}
			time.Sleep(20 * time.Second)
		}
	}()
}

func (c *ReplicasetController) Run() {
	c.SyncTest()
	listwatch.Watch(global.ReplicasetTopic(), c.ReplicasetHandler)
}

func NewController() Controller {
	return &ReplicasetController{
		Workers: make(map[string]Worker),
	}
}
