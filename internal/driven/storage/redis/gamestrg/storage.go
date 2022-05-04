package gamestrg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/go-redis/redis/v8"
	"gopkg.in/validator.v2"
)

type Storage struct {
	redisClient *redis.Client
}

type Config struct {
	RedisClient *redis.Client `validate:"nonnil"`
}

func New(cfg Config) (*Storage, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, err
	}
	s := &Storage{redisClient: cfg.RedisClient}
	return s, nil
}

func (s *Storage) PutGame(ctx context.Context, g core.Game) error {
	key := fmt.Sprintf("game:%s", g.GameID)
	str, err := json.Marshal(g)
	if err != nil {
		return fmt.Errorf("Unable to marshal game: %w", err)
	}
	err = s.redisClient.Set(ctx, key, string(str), 0).Err()
	if err != nil {
		return fmt.Errorf("Unable to put game in redis: %w", err)
	}
	return nil
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*core.Game, error) {
	str, err := s.redisClient.Get(ctx, fmt.Sprintf("game:%s", gameID)).Result()
	if err != nil {
		return nil, fmt.Errorf("unable to get game from redis: %w", err)
	}
	var g core.Game
	err = json.Unmarshal([]byte(str), &g)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal game: %w", err)
	}
	return &g, nil
}
