package queststrg

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

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

func (s *Storage) GetRandomQuestion(ctx context.Context) (*core.Question, error) {
	var qList []core.Question
	val, err := s.redisClient.Get(ctx, "questions").Result()
	if err != nil {
		return nil, fmt.Errorf("unable to get question: %v", err)
	}
	err = json.Unmarshal([]byte(val), &qList)
	if err != nil {
		return nil, fmt.Errorf("unablr to unmarshal questions: %v", err)
	}

	r := rand.New(rand.New(rand.NewSource(time.Now().UnixNano())))
	idx := r.Intn(len(qList))

	return &qList[idx], nil
}
