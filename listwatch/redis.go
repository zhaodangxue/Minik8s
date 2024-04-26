package listwatch

import (
	"context"
	"fmt"
	"minik8s/global"

	"github.com/go-redis/redis/v8"
)

//	var client = redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//		Password: "",
//		DB: 0,
//	})
func Init() {
	var client = redis.NewClient(&redis.Options{
		Addr:     global.Host + ":6379",
		Password: "",
		DB:       0,
	})
	//检查redis是否连接成功
	err := client.Ping(context.Background()).Err()
	if err != nil {
		//连接失败
		fmt.Println("redis connect failed")
		return
	}
	fmt.Println("redis connect success")
}
