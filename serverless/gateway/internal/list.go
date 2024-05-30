package internal

import (
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
		Func:     EventHandlerOnList,
		Interval: EventListInterval * time.Second,
	},
}
