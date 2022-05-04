package testutils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/go-redis/redis/v8"
)

func InitRedisClient() (*redis.Client, error) {
	return redis.NewClient(&redis.Options{Addr: "redis:6379"}), nil
}

func ResetRedis(client *redis.Client) error {
	return client.FlushAll(context.Background()).Err()
}

func InsertQuestion(client *redis.Client, questions []core.Question) error {
	b, err := json.Marshal(questions)
	if err != nil {
		return fmt.Errorf("unable to marshal questions test: %w", err)
	}
	err = client.Set(context.Background(), "question", string(b), 0).Err()
	if err != nil {
		return fmt.Errorf("unable to set questions test: %w", err)
	}
	return nil
}
