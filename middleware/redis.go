package middleware

import (
	"log"

	"github.com/go-redis/redis/v8"
)

func GetRedisClient() *(redis.Client) {
	opt, err := redis.ParseURL("redis://127.0.0.1:6379/0")
	if err != nil {
		log.Fatal(err)
	}
	return redis.NewClient(opt)
}
