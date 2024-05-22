package listwatch

import (
	"context"
	"minik8s/global"

	"github.com/redis/go-redis/v9"
)

//	var client = redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//		Password: "",
//		DB: 0,
//	})
var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr:     global.Host + ":6379",
	Password: "",
	DB:       0,
})

func Subscribe(topic string) *redis.PubSub {
	return client.Subscribe(ctx, topic)
}
func Publish(topic string, message string) {
	client.Publish(ctx, topic, message)
}
