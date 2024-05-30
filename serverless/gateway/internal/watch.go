package internal

import (
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/serverless/gateway/internal/event"
)

var WatchTable = map[string]listwatch.HandlerOnWatch{
	// "Topic": handlerFunc,
	global.FunctionTopic(): FunctionHandlerOnWatch,
	global.EventTopic():    event.EventHandlerOnWatch,
}
