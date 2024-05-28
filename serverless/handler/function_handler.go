package serverless_handler

import (
	"fmt"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func FunctionHandler(c *gin.Context) {
	name := c.Param("name")
	svcUrl := global.DefaultNamespace + "/" + "function-"+name+"-service"
	svc := apiobjects.Service{}
	err := utils.GetUnmarshal("http://localhost:8080/api/get/oneservice/"+svcUrl, &svc)
	if err != nil {
		fmt.Println("error")
	}
	response,err := utils.PostWithJson(svc.Status.ClusterIP+":8080", svc)

	c.JSON(200, gin.H{
		"message": "FunctionHandler",
		"response": response,
	})
}
