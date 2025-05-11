package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func initRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
		DB:       0,
	})

	ctx := context.Background()

	user := User{ID: 1, Name: "John Doe"}

	err := client.Set(ctx, "user:1", user, 5*time.Minute).Err()
	if err != nil {
	}

	err = client.Get(ctx, "user:1").Scan(&user)
	if err != nil {
	}
}
