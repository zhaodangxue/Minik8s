package event

import (
	"minik8s/apiobjects"
	"sync"
)

// CHECK: thread safe?
type EventStorage struct {
	Events sync.Map
	TimerExecutors sync.Map
}

var EventStorageInstance *EventStorage = &EventStorage{
	Events: sync.Map{},
	TimerExecutors: sync.Map{},
}
