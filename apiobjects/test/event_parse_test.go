package api_test

import (
	"encoding/json"
	"minik8s/apiobjects"
	"os"
	"testing"
)


func TestEventYaml(t *testing.T) {
	var event apiobjects.Event
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dir)
	data, err := os.ReadFile("../examples/event/event.json")
	if err != nil {
		t.Error(err)
		return
	}
	err = json.Unmarshal(data, &event)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(event)
}
