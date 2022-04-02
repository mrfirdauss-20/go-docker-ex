package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

type Service interface {
	NewGame(ctx context.Context, input NewGameInput) (*NewGameOutput, error)
	NewQuestion(ctx context.Context, input NewQuestionInput) (*NewQuestionOutput, error)
	SubmitAnswer(ctx context.Context, input SubmitAnswerInput) (*SubmitAnswerOutput, error)
}

type ServiceConfig struct {
	GameStorage       GameStorage       `validate:"nonnil"`
	QuestionStorage   QuestionStorage   `validate:"nonnil"`
	TimeoutCalculator TimeoutCalculator `validate:"nonnil"`
	Clock             Clock             `validate:"nonnil"`
	AddScore          int               `validate:"min=1"`
}

func (c ServiceConfig) Validate() error {
	return validator.Validate(c)
}

func NewService(cfg ServiceConfig) (Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	s := &service{
		gameStorage:       cfg.GameStorage,
		questionStorage:   cfg.QuestionStorage,
		timeoutCalculator: cfg.TimeoutCalculator,
		clock:             cfg.Clock,
		addScore:          cfg.AddScore,
	}
	return s, nil
}

type service struct {
	gameStorage       GameStorage
	questionStorage   QuestionStorage
	timeoutCalculator TimeoutCalculator
	clock             Clock
	addScore          int
}

func (s *service) NewGame(ctx context.Context, input NewGameInput) (*NewGameOutput, error) {
	// validate input
	err := input.Validate()
	if err != nil {
		return nil, err
	}
	// initialize game instance
	game := Game{
		GameID:          uuid.NewString(),
		PlayerName:      input.PlayerName,
		Scenario:        ScenarioNewQuestion,
		Score:           0,
		CountCorrect:    0,
		CurrentQuestion: nil,
	}
	// save game instance in storage
	err = s.gameStorage.PutGame(ctx, game)
	if err != nil {
		return nil, fmt.Errorf("unable to create game in storage due: %w", err)
	}
	// prepare output
	output := &NewGameOutput{
		GameID:     game.GameID,
		PlayerName: game.PlayerName,
		Scenario:   game.Scenario,
	}
	return output, nil
}

func (s *service) NewQuestion(ctx context.Context, input NewQuestionInput) (*NewQuestionOutput, error) {
	// validate input
	err := input.Validate()
	if err != nil {
		return nil, err
	}
	// get game instance
	game, err := s.gameStorage.GetGame(ctx, input.GameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game instance due: %w", err)
	}
	if game == nil {
		return nil, ErrGameNotFound
	}
	// check if game scenario is NEW_QUESTION
	if game.Scenario != ScenarioNewQuestion {
		return nil, ErrInvalidScenario
	}
	// get random question from storage
	question, err := s.questionStorage.GetRandomQuestion(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get random question due: %w", err)
	}
	// set question to game instance
	game.CurrentQuestion = &TimedQuestion{
		Question: *question,
		Timeout:  s.timeoutCalculator.CalcTimeout(ctx, game.CountCorrect),
	}
	// update game scenario
	game.Scenario = ScenarioSubmitAnswer
	// update game in storage
	err = s.gameStorage.PutGame(ctx, *game)
	if err != nil {
		return nil, fmt.Errorf("unable to update game storage due: %w", err)
	}
	// prepare output
	output := &NewQuestionOutput{
		GameID:   game.GameID,
		Scenario: game.Scenario,
		Problem:  game.CurrentQuestion.Problem,
		Choices:  game.CurrentQuestion.Choices,
		Timeout:  game.CurrentQuestion.Timeout,
	}
	return output, nil
}

func (s *service) SubmitAnswer(ctx context.Context, input SubmitAnswerInput) (*SubmitAnswerOutput, error) {
	// validate input
	err := input.Validate()
	if err != nil {
		return nil, err
	}
	// get game instance
	game, err := s.gameStorage.GetGame(ctx, input.GameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game instance due: %w", err)
	}
	if game == nil {
		return nil, ErrGameNotFound
	}
	// check if game scenario is SUBMIT_ANSWER
	if game.Scenario != ScenarioSubmitAnswer {
		return nil, ErrInvalidScenario
	}
	// check if it is already timeout or answer is incorrect
	duration := int(input.SentAt - input.StartAt)
	isTimeout := duration >= game.CurrentQuestion.Timeout
	isIncorrect := input.AnswerIdx != game.CurrentQuestion.CorrectIndex
	if isTimeout || isIncorrect {
		game.Scenario = ScenarioGameOver
	} else {
		game.Scenario = ScenarioNewQuestion
		game.CountCorrect += 1
		game.Score += s.addScore
	}
	// save game in database
	err = s.gameStorage.PutGame(ctx, *game)
	if err != nil {
		return nil, fmt.Errorf("unable to update game in storage due: %w", err)
	}
	// prepare output
	output := &SubmitAnswerOutput{
		GameID:       game.GameID,
		Scenario:     game.Scenario,
		AnswerIndex:  input.AnswerIdx,
		CorrectIndex: game.CurrentQuestion.CorrectIndex,
		Duration:     duration,
		Timeout:      game.CurrentQuestion.Timeout,
		Score:        game.Score,
	}
	return output, nil
}
