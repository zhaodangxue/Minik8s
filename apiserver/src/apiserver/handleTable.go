package apiserver

import (
	"minik8s/apiserver/src/handler"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc

var PostTable = map[string]HandlerFunc{
	route.TestPostPath: handler.TestPostHandler,
	route.TestCtlPath:  handler.TestCtlHandler,
}
var GetTable = map[string]HandlerFunc{
	route.TestGetPath: handler.TestGetHandler,
	route.TestCtlPath: handler.TestCtlGetHandler,
}
var PutTable = map[string]HandlerFunc{
	route.TestPutPath: handler.TestPutHandler,
}
var DeleteTable = map[string]HandlerFunc{
	route.TestDeletePath: handler.TestDeleteHandler,
	route.TestCtlPath:    handler.TestCtlDeleteHandler,
}
