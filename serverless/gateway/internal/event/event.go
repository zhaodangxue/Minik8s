package event

import (
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
