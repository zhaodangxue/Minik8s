package listwatch

import (
	"fmt"

	"github.com/go-redis/redis/v8"
)

type HandlerOnWatch func(message *redis.Message)

func Watch(topic string, handler HandlerOnWatch) {
	subscribe := client.Subscribe(ctx, topic)
	ch := subscribe.Channel()
	for msg := range ch {
		fmt.Printf("msg-payload:%v, msg-channel:%v", msg.Payload, msg.Channel)
		handler(msg)
	}
}
