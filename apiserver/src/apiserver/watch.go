package apiserver

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

func SyncTest(message *redis.Message) {
	fmt.Printf("msg-payload:%v, msg-channel:%v", message.Payload, message.Channel)
}
