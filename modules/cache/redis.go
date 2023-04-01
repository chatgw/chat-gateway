package cachemod

import (
	"context"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func NewRedis() (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	// Redis
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       db,
	})
	err = client.Ping(context.TODO()).Err()

	return client, err
}
