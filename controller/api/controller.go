package api

import (
	"minik8s/apiobjects"
	"time"
)

type Controller interface {
	Init(InitStruct)
	GetListFuncEnvelops() []ListFuncEnvelop
	GetWatchFuncEnvelops() []WatchFuncEnvelop
}

type ListFunc func(controller Controller) error

type ListFuncEnvelop struct {
	Func     ListFunc
	Interval time.Duration
}

type WatchFunc func(controller Controller, message apiobjects.TopicMessage) error

type WatchFuncEnvelop struct {
	Func  WatchFunc
	Topic string
}

// 传入必须的参数初始化必须的参数
type InitStruct struct {
	ApiserverPath   string
}

func (init *InitStruct)GetApiserverUrl() string {
	return "http://" + init.ApiserverPath
}
