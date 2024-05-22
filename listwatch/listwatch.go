package listwatch

import (
	"minik8s/utils"

	"github.com/redis/go-redis/v9"
)

type HandlerOnWatch func(message *redis.Message)

func Watch(topic string, handler HandlerOnWatch) {
	subscribe := client.Subscribe(ctx, topic)
	ch := subscribe.Channel()
	for msg := range ch {
		utils.Debug("msg-payload:", msg.Payload, "\n msg-channel:", msg.Channel)
		handler(msg)
	}
}
