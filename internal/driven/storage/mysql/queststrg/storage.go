package queststrg

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/jmoiron/sqlx"
	"gopkg.in/validator.v2"
)

type Storage struct {
	sqlClient *sqlx.DB
}

type Config struct {
	SQLClient *sqlx.DB `validate:"nonnil"`
}

func New(cfg Config) (*Storage, error) {
	err := validator.Validate(cfg)
	if err != nil {
		return nil, err
	}
	s := &Storage{sqlClient: cfg.SQLClient}
	return s, nil
}

func (s *Storage) GetRandomQuestion(ctx context.Context) (*core.Question, error) {
	// fetch data from database
	var rows []questionRow
	query := `SELECT problem, correct_index, answers FROM questions`
	err := s.sqlClient.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	// parse rows into questions
	var questions []core.Question
	for _, row := range rows {
		q, err := row.toQuestion()
		if err != nil {
			continue
		}
		questions = append(questions, *q)
	}
	// select random question
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(questions))

	return &questions[idx], nil
}
