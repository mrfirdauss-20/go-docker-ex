package gamestrg

import (
	"github.com/ghazlabs/hex-mathrush/internal/core"
)

type gameRow struct {
	ID              string
	PlayerName      string
	Scenario        string
	Score           int
	CountCorrect    int
	QuestionID      *int
	QuestionTimeout int
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
	if g.CurrentQuestion != nil {
		gRow.QuestionTimeout = g.CurrentQuestion.Timeout
	}
	return gRow
}

type gameCompleteRow struct {
	ID              string
	PlayerName      string
	Scenario        string
	Score           int
	CountCorrect    int
	QuestionID      *int
	Problem         *string
	CorrectIndex    *int
	Answers         []string
	QuestionTimeout int
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
		g.CurrentQuestion = &core.TimedQuestion{
			Question: core.Question{
				Problem:      *gRow.Problem,
				Choices:      gRow.Answers,
				CorrectIndex: *gRow.CorrectIndex,
			},
			Timeout: gRow.QuestionTimeout,
		}
	}
	return g, nil
}
