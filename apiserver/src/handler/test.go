package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/utils"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

type TestBody struct {
	Body    string  `yaml:"body"`
	Name    string  `yaml:"name"`
	Uid     string  `yaml:"uid"`
	Num     int     `yaml:"num"`
	Percent float64 `yaml:"percent"`
}

func TestHandler(c *gin.Context) {
	fmt.Println("test-success")
}
func TestPostHandler(c *gin.Context) {
	fmt.Println("test-post")
	obj := TestBody{}
	err := utils.ReadUnmarshal(c.Request.Body, &obj)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	var TestJson []byte
	TestJson, err = json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusOK, err.Error())
		return
	}
	etcd.Put("/test/post", string(TestJson))
	c.String(http.StatusOK, "ok")
	fmt.Println("test-post-success")
}
func TestGetHandler(c *gin.Context) {
	fmt.Println("test-get")
	value, _ := etcd.Get("/test/get")
	var TestJson TestBody
	err := json.Unmarshal([]byte(value), &TestJson)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, TestJson)
	fmt.Println("test-get-success")
}
func TestPutHandler(c *gin.Context) {
	fmt.Println("test-put")
	name := c.Param("name")
	uid := c.Param("uid")
	defer c.Request.Body.Close()
	var Input TestBody
	content, _ := io.ReadAll(c.Request.Body)
	err := json.Unmarshal(content, &Input)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	path := path.Join("/test/put", name, uid)
	value, _ := etcd.Get(path)
	var TestJson TestBody
	_ = json.Unmarshal([]byte(value), &TestJson)
	TestJson.Name = Input.Name
	TestJson.Uid = Input.Uid
	TestJson.Num = Input.Num
	TestJson.Percent = Input.Percent
	TestJson.Body = Input.Body
	data, _ := json.Marshal(TestJson)
	etcd.Put(path, string(data))
	c.String(http.StatusOK, "set TestBody uuid:%s name:%s success to TestBody %v", uid, name, TestJson)
	fmt.Println("test-put-success")
}
func TestDeleteHandler(c *gin.Context) {
	fmt.Println("test-delete")
	name := c.Param("name")
	uid := c.Param("uid")
	path := path.Join("/test/delete", name, uid)
	etcd.Delete(path)
	c.String(http.StatusOK, "delete TestBody uuid:%s name:%s success", uid, name)
	fmt.Println("test-delete-success")
}
func TestCtlHandler(c *gin.Context) {
	fmt.Println("test-ctl")
	obj := apiobjects.TestYaml{}
	err := utils.ReadUnmarshal(c.Request.Body, &obj)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	var TestJson []byte
	TestJson, err = json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		c.String(http.StatusOK, err.Error())
		return
	}
	etcd.Put("/test/ctl", string(TestJson))
	c.String(http.StatusOK, "ok")
	fmt.Println("test-ctl-success")
}
func TestCtlGetHandler(c *gin.Context) {
	fmt.Println("test-ctl-get")
	value, _ := etcd.Get("/test/ctl")
	var TestJson apiobjects.TestYaml
	err := json.Unmarshal([]byte(value), &TestJson)
	if err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.JSON(http.StatusOK, TestJson)
	fmt.Println("test-ctl-get-success")
}
func TestCtlDeleteHandler(c *gin.Context) {
	fmt.Println("test-ctl-delete")
	etcd.Delete("/test/ctl")
	c.String(http.StatusOK, "delete TestYaml success")
	fmt.Println("test-ctl-delete-success")
}
