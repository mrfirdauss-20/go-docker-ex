package queststrg

import (
	"context"
	"math/rand"
	"strconv"
	"time"

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

func (s *Storage) GetRandomQuestion(ctx context.Context) (*core.Question, error) {

	//res, err := rdb.Do(ctx, "set", "key", "value").Result()
	//val2, err := rdb.Get(ctx, "key2").Result()
	var questions []core.Question
	i := 0
	atr, err := s.redisClient.Get(ctx, "questions:0").Result()
	if err != nil {
		panic(err)
	}
	for atr != "" {
		CorIdx, err := s.redisClient.Get(ctx, "questions:"+strconv.Itoa(i)+":correct_index").Result()
		Ans, err := s.redisClient.LRange(ctx, "questions:"+strconv.Itoa(i)+":answers", 0, -1).Result()
		corId, err := strconv.Atoi(CorIdx)
		q := core.Question{
			Problem:      atr,
			CorrectIndex: corId,
			Choices:      Ans,
		}
		questions = append(questions, q)
		i++
		atr, err = s.redisClient.Get(ctx, "questions:"+strconv.Itoa(i)).Result()
		if err != nil {
			continue
		}
	}

	if len(questions) == 0 {
		return nil, nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(questions))
	return &questions[idx], nil
}
