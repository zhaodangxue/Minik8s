package internal

import (
	"minik8s/global"
	"minik8s/listwatch"
)

var WatchTable = map[string]listwatch.HandlerOnWatch{
	// "Topic": handlerFunc,
	global.FunctionTopic(): FunctionHandlerOnWatch,
	global.EventTopic():    EventHandlerOnWatch,
}
