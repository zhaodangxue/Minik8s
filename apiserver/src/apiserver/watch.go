package apiserver

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

func SyncTest(message *redis.Message) {
	fmt.Printf("msg-payload:%v, msg-channel:%v", message.Payload, message.Channel)
}
