package internal

import (
	"encoding/json"
	"errors"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	events "minik8s/serverless/gateway/internal/event"

	"github.com/redis/go-redis/v9"
)

const EventListInterval = 10

func EventHandlerOnList() (err error) {
	utils.Info("Update event list")
	vals, _ := etcd.Get_prefix(route.EventPath)
	for _, val := range vals {
		event := apiobjects.Event{}
		err = json.Unmarshal([]byte(val), &event)
		if err != nil {
			utils.Error("Unmarshal event failed: ", err)
			continue
		}
		_, ok := events.EventStorageInstance.Events.Load(event.Name)
		if ok {
			continue
		}
		err = createEvent(&event)
		if err != nil {
			utils.Error("Create event failed: ", err)
			continue
		}
	}
	return
}

func generateTrigerFunction(event *apiobjects.Event) func() {
	// CHECK: 考虑event的use after free
	return func() {
		utils.Info("Trigger event: ", event.Name)
		for _, workflow := range event.Workflows {
			name := string(workflow)
			// Start workflow
			go WorkflowExecutor(name, make(map[string]interface{}));
		}
	}
}

func addTimerEvent(event *apiobjects.Event) {
	timerExecutor := events.TimerExecutor{}
	timerExecutor.Init(event.TimeEvent.Cron, generateTrigerFunction(event))
	orival, ok := events.EventStorageInstance.TimerExecutors.Load(event.Name)
	if ok {
		orival.(*events.TimerExecutor).Stop()
	}
	events.EventStorageInstance.TimerExecutors.Store(event.Name, &timerExecutor)
	timerExecutor.Start()
}

func removeEvent(name string) error {
	orival, ok := events.EventStorageInstance.Events.Load(name)
	if !ok {
		return errors.New("Event not found: " + name)
	}
	event := orival.(*apiobjects.Event)
	switch event.Type {
	case apiobjects.EventTypeTimer:
		timerExecutor, ok := events.EventStorageInstance.TimerExecutors.Load(name)
		if !ok {
			return errors.New("Corresponse TimerExecutor not found: " + name)
		}
		timerExecutor.(*events.TimerExecutor).Stop()
		events.EventStorageInstance.TimerExecutors.Delete(name)
	default:
		return errors.New("Unknown event type: " + (string)(event.Type))
	}
	events.EventStorageInstance.Events.Delete(name)
	return nil
}

func createEvent(event *apiobjects.Event) error {
	_, ok := events.EventStorageInstance.Events.Load(event.Name)
	if ok {
		return errors.New("Event already exists: " + event.Name)
	}
	switch event.Type {
	case apiobjects.EventTypeTimer:
		events.EventStorageInstance.Events.Store(event.Name, event)
		addTimerEvent(event)
	default:
		return errors.New("Unknown event type: " + (string)(event.Type))
	}
	return nil
}

// TODO
func EventHandlerOnWatch(msg *redis.Message) {
	utils.Info("EventHandlerOnWatch")

	msgdata := msg.Payload
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msgdata), &topicMessage)
	if err != nil {
		utils.Error("EventHandlerOnWatch: Unmarshal topic message failed: ", err)
		return
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		utils.Info("EventHandlerOnWatch: Create event: ", topicMessage.Object)

		event := apiobjects.Event{}
		err = json.Unmarshal([]byte(topicMessage.Object), &event)
		if err != nil {
			utils.Error("EventHandlerOnWatch: Unmarshal event failed: ", err)
			return
		}
		err = createEvent(&event)
		if err != nil {
			utils.Error("EventHandlerOnWatch: Create event failed: ", err)
			return
		}
	case apiobjects.Delete:
		utils.Info("EventHandlerOnWatch: Delete event: ", topicMessage.Object)

		event := apiobjects.Event{}
		err = json.Unmarshal([]byte(topicMessage.Object), &event)
		if err != nil {
			utils.Error("EventHandlerOnWatch: Unmarshal event failed: ", err)
			return
		}
		err = removeEvent(event.Name)
		if err != nil {
			utils.Error("EventHandlerOnWatch: Remove event failed: ", err)
			return
		}
	default:
		utils.Error("EventHandlerOnWatch: Unknown action type: ", topicMessage.ActionType)
	}
}
