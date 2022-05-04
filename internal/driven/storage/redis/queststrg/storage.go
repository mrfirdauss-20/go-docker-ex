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

var questList []core.Question

func New(cfg Config) (*Storage, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, err
	}
	s := &Storage{redisClient: cfg.RedisClient}
	//insert qurestion berupa list of Question struct
	questList = []core.Question{
		{
			Problem:      "1 + 2",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 3,
		},
		{
			Problem:      "1 + 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 2,
		},
		{
			Problem:      "2 + 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 3,
		},
		{
			Problem:      "3 - 2",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 1,
		},
		{
			Problem:      "2 - 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 1,
		},
		{
			Problem:      "2 + 2 - 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 3,
		},
		{
			Problem:      "1 + 1 + 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 3,
		},
		{
			Problem:      "3 - 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 2,
		},
		{
			Problem:      "2 + 1 - 2",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 1,
		},
		{
			Problem:      "1 + 1 - 1",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 1,
		},
	}
	str, err := json.Marshal(questList)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal question: %v", err)
	}
	err = s.redisClient.Set(context.Background(), "questions", string(str), 0).Err()
	if err != nil {
		return nil, fmt.Errorf("unable to set question: %v", err)
	}
	return s, nil
}

func (s *Storage) GetRandomQuestion(ctx context.Context) (*core.Question, error) {

	//res, err := rdb.Do(ctx, "set", "key", "value").Result()
	//val2, err := rdb.Get(ctx, "key2").Result()
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
