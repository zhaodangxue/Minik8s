package internal

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc

var PostTable = map[string]HandlerFunc{
	"/serverless/function/:name": FunctionHandler,
	"/serverless/workflow/:name": WorkflowHandler,
	"/serverless/test":           TestHandler,
}

var GetTable = map[string]HandlerFunc{}

var PutTable = map[string]HandlerFunc{}

var DeleteTable = map[string]HandlerFunc{}
