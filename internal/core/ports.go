package core

import (
	"context"
	"time"
)

type GameStorage interface {
	// PutGame is used for putting the game instance on storage. If game is already
	// exists on storage, it will be overwritten.
	PutGame(ctx context.Context, g Game) error

	// GetGame returns game instance stored on storage. Returns nil when game is not
	// found on storage.
	GetGame(ctx context.Context, gameID string) (*Game, error)
}

type QuestionStorage interface {
	// GetRandomQuestion returns random question from storage.
	GetRandomQuestion(ctx context.Context) (*Question, error)
}

type Clock interface {
	// Now returns current time.
	Now() time.Time
}

type TimeoutCalculator interface {
	// CalcTimeout returns proper duration to answer question for given number
	// of correct answers.
	CalcTimeout(ctx context.Context, countCorrect int) int
}
