package kubectl__test

import (
	"fmt"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	command "minik8s/kubectl/src"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var apiServer apiserver.ApiServer

func TestRunApply(t *testing.T) {
	fmt.Println("TestRunApply")
	etcd.Clear()
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunApply_test("./test.yaml")
	assert.Nil(t, err)
}
func TestRunGet(t *testing.T) {
	fmt.Println("TestRunGet")
	apiServer = apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	err := command.RunGet_test("test", "111")
	assert.Nil(t, err)
}
