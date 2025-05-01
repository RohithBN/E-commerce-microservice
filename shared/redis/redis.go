package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var RedisClient *redis.Client

func ConnectRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Replace with your Redis server address
		Password: "",               // No password for local development
		DB:       0,                // Default DB
	})

	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err		
	}
	fmt.Println("Redis connected", pong)
	return nil
}




