package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ConnectToRedis(connString string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
		DB:       0,
	})

	status := client.Ping(context.Background())
	if err := status.Err(); err != nil {
		return nil, err
	}

	return client, nil
}
