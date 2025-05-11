package cache

import (
	"context"
	"time"
)

func SetBalance(walletId string, balance int64) error {
	ctx := context.Background()

	user := User{ID: 1, Name: "John Doe"}

	err := client.Set(ctx, "user:1", user, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetBalance(walletId string) (int64, error) {
	ctx := context.Background()

	user := User{ID: 1, Name: "John Doe"}

	err := client.Get(ctx, "user:1").Scan(&user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
