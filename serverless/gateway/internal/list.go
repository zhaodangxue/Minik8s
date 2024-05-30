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
		Interval: 60 * time.Second,
	},
}
