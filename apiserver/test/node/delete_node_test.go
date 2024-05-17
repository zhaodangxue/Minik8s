package test

import (
	"minik8s/utils"
	"testing"
)

func TestDeleteNode(t *testing.T) {
	t.Log("TestDeleteNode")

 	_, err := utils.Delete("http://localhost:8080/api/node/system/node-2b52c8")
	if err != nil {
		t.Error(err)
	}
}
