package queststrg

import (
	"context"
	"math/rand"
	"time"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"gopkg.in/validator.v2"
)

type Storage struct {
	questions []core.Question
}

type Config struct {
	Questions []core.Question `validate:"min=1"`
}

func (c Config) Validate() error {
	return validator.Validate(c)
}

func New(cfg Config) (*Storage, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &Storage{questions: cfg.Questions}, nil
}

func (s *Storage) GetRandomQuestion(ctx context.Context) (*core.Question, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(s.questions))

	return &s.questions[idx], nil
}
