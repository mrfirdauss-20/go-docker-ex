package gamestrg

import (
	"context"
	"strconv"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/go-redis/redis/v8"
	"gopkg.in/validator.v2"
)

type Storage struct {
	redisClient *redis.Client
}

type Config struct {
	redisClient *redis.Client `validate:"nonnil"`
}

func New(cfg Config) (*Storage, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, err
	}
	s := &Storage{redisClient: cfg.redisClient}
	return s, nil
}

func (s *Storage) PutGame(ctx context.Context, g core.Game) error {
	var cursor uint64
	var n int
	var keys []string
	var err error
	keys, cursor, err = s.redisClient.Scan(ctx, cursor, "game:*", 10).Result()
	n = len(keys)
	err = s.redisClient.Set(ctx, "game:"+strconv.Itoa(n), g.GameID, 0).Err()
	if err != nil {
		return err
	}
	err = s.redisClient.Set(ctx, "game:"+strconv.Itoa(n)+":player_name", g.PlayerName, 0).Err()
	if err != nil {
		return err
	}
	err = s.redisClient.Set(ctx, "game:"+strconv.Itoa(n)+":scenario", g.Scenario, 0).Err()
	if err != nil {
		return err
	}
	err = s.redisClient.Set(ctx, "game:"+strconv.Itoa(n)+":score", g.Score, 0).Err()
	if err != nil {
		return err
	}
	err = s.redisClient.Set(ctx, "game:"+strconv.Itoa(n)+":count_correct", g.CountCorrect, 0).Err()
	if err != nil {
		return err
	}

	var qs []string
	qs, cursor, err = s.redisClient.Scan(ctx, cursor, "question:*:problem", 10).Result()

	i := 0
	iter := ""
	for i < len(qs) && g.CurrentQuestion.Problem != iter {
		iter, err = s.redisClient.Get(ctx, "question:"+strconv.Itoa(i)+":problem").Result()
		i++
	}

	err = s.redisClient.Set(ctx, "game:"+strconv.Itoa(n)+":question_id", i-1, 0).Err()
	if err != nil {
		return err
	}
	return err
}
