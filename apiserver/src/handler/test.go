package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestHandler(c *gin.Context) {
	fmt.Println("test-success")
}
