package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func FunctionHandler(c *gin.Context) {
	name := c.Param("name")

	value, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	var jsonParam map[string]interface{}
	err = json.Unmarshal(value, &jsonParam)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	svcUrl := global.DefaultNamespace + "/" + "function-" + name + "-service"
	svc := apiobjects.Service{}
	err = utils.GetUnmarshal("http://localhost:8080/api/get/oneservice/"+svcUrl, &svc)
	if err != nil {
		fmt.Println("Get service error: ", err.Error())
		c.String(500, err.Error())
		return
	}
	val := JudgeReplicas(name)
	if val != "success" {
		fmt.Println(val)
	}
	AddQpsCounter(name)
	response, err := utils.PostWithJson("http://"+svc.Status.ClusterIP+":8080", jsonParam)
	if err != nil {
		fmt.Println("Post error: ", err.Error())
		c.String(500, err.Error())
		return
	}
	utils.Info("response: ", response)
	value, err = io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		fmt.Println("Read response error: ", err.Error())
		c.String(500, err.Error())
		return
	}

	var jsonParam2 map[string]interface{}
	err = json.Unmarshal(value, &jsonParam2)
	if err != nil {
		fmt.Println("Unmarshal response error: ", err.Error())
		c.String(500, err.Error())
		return
	}
	c.JSON(200, jsonParam2)
}
