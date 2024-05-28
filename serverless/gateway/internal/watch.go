package internal

import (
	"minik8s/global"
	"minik8s/listwatch"
	serverless_handler "minik8s/serverless/handler"
)

var WatchTable = map[string]listwatch.HandlerOnWatch{
	// "Topic": handlerFunc,
	global.FunctionTopic(): serverless_handler.FunctionHandlerOnWatch,
}
