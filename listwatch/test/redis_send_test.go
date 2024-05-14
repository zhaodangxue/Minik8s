package test

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/listwatch"
	"testing"
)

func TestRedisSend(t *testing.T) {
	t.Log("TestRedisSend")

	message := apiobjects.TopicMessage{
		ActionType: apiobjects.Update,
		Object: "Hello from test!",
	}

	messageStr, err := json.Marshal(message)
	if err != nil {
		t.Error("Failed to marshal message:", err)
		return
	}

	listwatch.Publish(global.TestTopic(), string(messageStr))
}
