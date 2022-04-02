package core

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type NewGameInput struct {
	PlayerName string `validate:"nonzero"`
}

func (i NewGameInput) Validate() error {
	return validator.Validate(i)
}

type NewGameOutput struct {
	GameID     string `json:"game_id"`
	PlayerName string `json:"player_name"`
	Scenario   string `json:"scenario"`
}

type NewQuestionInput struct {
	GameID string `validate:"nonzero"`
}

func (i NewQuestionInput) Validate() error {
	return validator.Validate(i)
}

type NewQuestionOutput struct {
	GameID   string   `json:"game_id"`
	Scenario string   `json:"scenario"`
	Problem  string   `json:"problem"`
	Choices  []string `json:"choices"`
	Timeout  int      `json:"timeout"`
}

type SubmitAnswerInput struct {
	GameID    string `validate:"nonzero"`
	AnswerIdx int    `validate:"min=0"`
	StartAt   int64  `validate:"nonzero"`
	SentAt    int64  `validate:"nonzero"`
}

func (i SubmitAnswerInput) Validate() error {
	if i.SentAt < i.StartAt {
		return fmt.Errorf("invalid value of `sent_at`")
	}
	return validator.Validate(i)
}

type SubmitAnswerOutput struct {
	GameID     string `json:"game_id"`
	Scenario   string `json:"scenario"`
	AnswerIdx  int    `json:"answer_idx"`
	CorrectIdx int    `json:"correct_idx"`
	Duration   int    `json:"duration"`
	Timeout    int    `json:"timeout"`
	Score      int    `json:"score"`
}

type Game struct {
	GameID          string
	PlayerName      string
	Scenario        string
	Score           int
	CountCorrect    int
	CurrentQuestion *TimedQuestion
}

type Question struct {
	QuestionID string
	Problem    string
	Choices    []string
	CorrectIdx int
}

type TimedQuestion struct {
	Question
	Timeout int
}

type Scenario string

const (
	ScenarioNewQuestion  = "NEW_QUESTION"
	ScenarioSubmitAnswer = "SUBMIT_ANSWER"
	ScenarioGameOver     = "GAME_OVER"
)
