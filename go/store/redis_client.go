package store

import (
	"github.com/go-redis/redis/v7"
	"os"
)

// docker run --name my-redis -p 6379:6379 --restart always --detach redis

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_DSN"),
	})
	_, err := Client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
