package apiserver

import (
	"minik8s/apiserver/src/handler"
	"minik8s/apiserver/src/route"

	"github.com/gin-gonic/gin"
)

type HandlerFunc = gin.HandlerFunc

var PostTable = map[string]HandlerFunc{
	route.TestPath: handler.TestHandler,
}
