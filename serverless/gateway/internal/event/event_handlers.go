package event

import (
	"encoding/json"
	"errors"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

const EventListInterval = 10

func EventHandlerOnList() (err error) {
	utils.Info("Update event list")
	vals, _ := etcd.Get_prefix(route.EventPath)
	for _, val := range vals {
		event := apiobjects.Event{}
		err := json.Unmarshal([]byte(val), &event)
		if err != nil {
			continue
		}
		err = createEvent(&event)
		if err != nil {
			utils.Error(err)
		}
	}
}

func addTimerEvent(event *apiobjects.Event) {
	timerExecutor := TimerExecutor{}
	// TODO: 替换为真实timer处理函数
	timerExecutor.Init(event.TimeEvent.Cron, func() {})
	orival, ok := EventStorageInstance.TimerExecutors.Load(event.Name)
	if ok {
		orival.(*TimerExecutor).Stop()
	}
	EventStorageInstance.TimerExecutors.Store(event.Name, &timerExecutor)
	timerExecutor.Start()
}

func removeEvent(name string) error {
	orival, ok := EventStorageInstance.Events.Load(name)
	if !ok {
		return errors.New("Event not found: " + name)
	}
	event := orival.(*apiobjects.Event)
	switch event.Type {
	case apiobjects.EventTypeTimer:
		timerExecutor, ok := EventStorageInstance.TimerExecutors.Load(name)
		if !ok {
			return errors.New("Corresponse TimerExecutor not found: " + name)
		}
		timerExecutor.(*TimerExecutor).Stop()
		EventStorageInstance.TimerExecutors.Delete(name)
	default:
		return errors.New("Unknown event type: " + (string)(event.Type))
	}
	EventStorageInstance.Events.Delete(name)
	return nil
}

func createEvent(event *apiobjects.Event) error {
	_, ok := EventStorageInstance.Events.Load(event.Name)
	if ok {
		return errors.New("Event already exists: " + event.Name)
	}
	switch event.Type {
	case apiobjects.EventTypeTimer:
		EventStorageInstance.Events.Store(event.Name, event)
		addTimerEvent(event)
	default:
		return errors.New("Unknown event type: " + (string)(event.Type))
	}
	return nil
}

// TODO
func CreateEventWatcher() {
	
}

// TODO
func DeleteEventWatcher() {
	
}