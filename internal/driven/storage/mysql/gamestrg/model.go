package gamestrg

import (
	"fmt"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/jmoiron/sqlx/types"
)

type gameRow struct {
	ID              string `db:"id"`
	PlayerName      string `db:"player_name"`
	Scenario        string `db:"scenario"`
	Score           int    `db:"score"`
	CountCorrect    int    `db:"count_correct"`
	QuestionID      *int   `db:"question_id"`
	QuestionTimeout int    `db:"question_timeout"`
}

func newGameRow(g core.Game, questionID *int) gameRow {
	gRow := gameRow{
		ID:           g.GameID,
		PlayerName:   g.PlayerName,
		Scenario:     g.Scenario,
		Score:        g.Score,
		CountCorrect: g.CountCorrect,
		QuestionID:   questionID,
	}
	// it is possible for the game to not have current active question
	if g.CurrentQuestion != nil {
		gRow.QuestionTimeout = g.CurrentQuestion.Timeout
	}
	return gRow
}

type gameCompleteRow struct {
	ID              string         `db:"id"`
	PlayerName      string         `db:"player_name"`
	Scenario        string         `db:"scenario"`
	Score           int            `db:"score"`
	CountCorrect    int            `db:"count_correct"`
	QuestionID      *int           `db:"question_id"`
	Problem         *string        `db:"problem"`
	CorrectIndex    *int           `db:"correct_index"`
	Answers         types.JSONText `db:"answers"`
	QuestionTimeout int            `db:"question_timeout"`
}

func (gRow gameCompleteRow) toGame() (*core.Game, error) {
	g := &core.Game{
		GameID:       gRow.ID,
		PlayerName:   gRow.PlayerName,
		Scenario:     gRow.Scenario,
		Score:        gRow.Score,
		CountCorrect: gRow.CountCorrect,
	}
	if gRow.QuestionID != nil {
		var choices []string
		err := gRow.Answers.Unmarshal(&choices)
		if err != nil {
			return nil, fmt.Errorf("unable to parse question due: %w", err)
		}
		g.CurrentQuestion = &core.TimedQuestion{
			Question: core.Question{
				Problem:      *gRow.Problem,
				Choices:      choices,
				CorrectIndex: *gRow.CorrectIndex,
			},
			Timeout: gRow.QuestionTimeout,
		}
	}
	return g, nil
}
