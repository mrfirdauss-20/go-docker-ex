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
	var quest []core.Question
	//get questions
	str, err := s.redisClient.Get(ctx, "questions").Result()
	if err != nil {
		return fmt.Errorf("failed to get games: %w", err)
	}
	err = json.Unmarshal([]byte(str), &quest)
	if err != nil {
		return fmt.Errorf("Unable to unmarshaling games: %w", err)
	}

	questId := 0
	if g.CurrentQuestion != nil {

		for quest[questId].Problem != g.CurrentQuestion.Problem {
			questId++
		}
		if questId >= len(quest) {
			return fmt.Errorf("question not found")
		}
	}

	//replace in db
	//get games from db
	var games []core.Game
	str, err = s.redisClient.Get(ctx, "games").Result()
	if err != nil {
		return fmt.Errorf("Unable to get games: %w", err)
	}
	err = json.Unmarshal([]byte(str), &games)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal games: %w", err)
	}
	//replace
	idx := 0
	for games[idx].GameID != g.GameID {
		idx++
	}
	games[idx] = g
	bts, err := json.Marshal(games)
	if err != nil {
		return fmt.Errorf("unable to marshal games: %w", err)
	}
	err = s.redisClient.Set(ctx, "games", string(bts), 0).Err()
	if err != nil {
		return fmt.Errorf("unable to set games: %w", err)
	}
	return nil
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*core.Game, error) {
	var games []core.Game
	str, err := s.redisClient.Get(ctx, "games").Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to get games: %w", err)
	}
	err = json.Unmarshal([]byte(str), &games)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal games: %w", err)
	}
	for _, g := range games {
		if g.GameID == gameID {
			return &g, nil
		}
	}
	return nil, fmt.Errorf("game not found")
}
