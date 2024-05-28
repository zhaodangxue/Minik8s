package internal

import (
	serverless_handler "minik8s/serverless/handler"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc

var PostTable = map[string]HandlerFunc{
	"/serverless/function/:name": serverless_handler.FunctionHandler,
	"/serverless/workflow/:name": serverless_handler.WorkflowHandler,
	"/serverless/test":           serverless_handler.TestHandler,
}

var GetTable = map[string]HandlerFunc{}

var PutTable = map[string]HandlerFunc{}

var DeleteTable = map[string]HandlerFunc{}
