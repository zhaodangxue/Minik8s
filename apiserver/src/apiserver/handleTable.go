package apiserver

import (
	"minik8s/apiserver/src/handler"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc

var PostTable = map[string]HandlerFunc{
	route.TestPostPath: handler.TestPostHandler,
}
var GetTable = map[string]HandlerFunc{
	route.TestGetPath: handler.TestGetHandler,
}
var PutTable = map[string]HandlerFunc{
	route.TestPutPath: handler.TestPutHandler,
}
var DeleteTable = map[string]HandlerFunc{
	route.TestDeletePath: handler.TestDeleteHandler,
}
