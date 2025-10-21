package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func (c *RedisClient) Close() error {
	return c.Client.Close()
}

func NewRedisClient() *RedisClient {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // no password by default
		DB:       0,                // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected")
	return &RedisClient{Client: rdb, Ctx: ctx}
}
