package internal

import (
	"minik8s/serverless/gateway/internal/event"
	"time"
)

type ListFunc func() error

type ListFuncEnvelop struct {
	Func     ListFunc
	Interval time.Duration
}

var ListTable = []ListFuncEnvelop{
	{
		Func:     FunctionHandlerOnList,
		Interval: 30 * time.Second,
	},
	{
		Func:     event.EventHandlerOnList,
		Interval: event.EventListInterval * time.Second,
	},
}
