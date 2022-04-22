package gamestrg

import (
	"context"
	"fmt"

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

func (s *Storage) PutGame(ctx context.Context, g core.Game) error {
	// start transaction
	tx, err := s.sqlClient.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction due: %w", err)
	}
	// defer commit or rollback transaction
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()
	// resolve question id
	var questionID int
	query := `SELECT id FROM questions WHERE problem = ?`
	err = s.sqlClient.GetContext(ctx, &questionID, query, g.CurrentQuestion.Problem)
	if err != nil {
		return fmt.Errorf("unable to resolve question id due: %w", err)
	}
	// put game in database
	gr := newGameRow(g, questionID)
	query = `
		REPLACE INTO games (
			id, player_name, scenario, score, 
			count_correct, question_id, question_timeout
		)
		VALUES (
			:id, :player_name, :scenario, :score, 
			:count_correct, :question_id, :question_timeout
		)
	`
	_, err = s.sqlClient.NamedExecContext(ctx, query, gr)
	if err != nil {
		return fmt.Errorf("unable to execute query to put game due: %w", err)
	}
	return nil
}

func (s *Storage) GetGame(ctx context.Context, gameID string) (*core.Game, error) {
	// get game along with question
	var gRow gameCompleteRow
	query := `
		SELECT 
			games.id,
			games.player_name,
			games.scenario, 
			games.score,
			games.count_correct,
			questions.problem,
			questions.correct_index,
			questions.answers,
			games.question_timeout 
		FROM 
			games JOIN questions ON questions.id = games.question_id
		WHERE games.id = ?
	`
	err := s.sqlClient.GetContext(ctx, &gRow, query, gameID)
	if err != nil {
		return nil, fmt.Errorf("unable to execute query due: %w", err)
	}
	// construct output
	game, err := gRow.toGame()
	if err != nil {
		return nil, fmt.Errorf("unable to construct game due: %w", err)
	}
	return game, nil
}
