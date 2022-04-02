package core

import (
	"context"
	"time"
)

type GameStorage interface {
	PutGame(ctx context.Context, g Game) error
	GetGame(ctx context.Context, gameID string) (*Game, error)
}

type QuestionStorage interface {
	GetRandomQuestion(ctx context.Context) (*Question, error)
}

type Clock interface {
	Now() time.Time
}

type TimeoutCalculator interface {
	CalcTimeout(ctx context.Context, countCorrect int) int
}
