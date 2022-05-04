package core

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	// define mock dependencies
	gameStorage := newMockGameStorage()
	questionStorage := newMockQuestionStorage(nil)
	timeoutCalculator := newMockTimeoutCalculator()
	clock := newMockClock()
	// define function for validating new service
	validateService := func(t *testing.T, s Service, cfg ServiceConfig) {
		assert.Equal(t, cfg.GameStorage, s.(*service).gameStorage)
		assert.Equal(t, cfg.QuestionStorage, s.(*service).questionStorage)
		assert.Equal(t, cfg.TimeoutCalculator, s.(*service).timeoutCalculator)
		assert.Equal(t, cfg.Clock, s.(*service).clock)
		assert.Equal(t, cfg.AddScore, s.(*service).addScore)
		assert.GreaterOrEqual(t, cfg.AddScore, 1)
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  ServiceConfig
		IsError bool
	}{
		{
			Name: "Test New Service Missing Game Storage",
			Config: ServiceConfig{
				GameStorage:       nil,
				QuestionStorage:   questionStorage,
				TimeoutCalculator: timeoutCalculator,
				Clock:             clock,
				AddScore:          1,
			},
			IsError: true,
		},
		{
			Name: "Test New Service Missing Question Storage",
			Config: ServiceConfig{
				GameStorage:       gameStorage,
				QuestionStorage:   nil,
				TimeoutCalculator: timeoutCalculator,
				Clock:             clock,
				AddScore:          1,
			},
			IsError: true,
		},
		{
			Name: "Test New Service Missing Timeout Calculator",
			Config: ServiceConfig{
				GameStorage:       gameStorage,
				QuestionStorage:   questionStorage,
				TimeoutCalculator: nil,
				Clock:             clock,
				AddScore:          1,
			},
			IsError: true,
		},
		{
			Name: "Test New Service Missing Clock",
			Config: ServiceConfig{
				GameStorage:       gameStorage,
				QuestionStorage:   questionStorage,
				TimeoutCalculator: timeoutCalculator,
				Clock:             nil,
				AddScore:          1,
			},
			IsError: true,
		},
		{
			Name: "Test New Service Invalid Score",
			Config: ServiceConfig{
				GameStorage:       gameStorage,
				QuestionStorage:   questionStorage,
				TimeoutCalculator: timeoutCalculator,
				Clock:             clock,
				AddScore:          0,
			},
			IsError: true,
		},
		{
			Name: "Test New Service Valid Config",
			Config: ServiceConfig{
				GameStorage:       gameStorage,
				QuestionStorage:   questionStorage,
				TimeoutCalculator: timeoutCalculator,
				Clock:             clock,
				AddScore:          1,
			},
			IsError: false,
		},
	}
	// execute test cases
	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			service, err := NewService(testCases[i].Config)
			assert.Equal(t, testCases[i].IsError, (err != nil), "unexpected error")
			if service == nil {
				return
			}
			validateService(t, service, testCases[i].Config)
		})
	}
}

func TestNewGame(t *testing.T) {
	// initialize new service
	service, err := newNewService()
	require.NoError(t, err)
	// test the new game function
	gameInput := NewGameInput{
		PlayerName: "Riandy",
	}
	output, err := service.NewGame(context.Background(), gameInput)
	require.NoError(t, err)
	// match the output
	assert.Equal(t, "NEW_QUESTION", output.Scenario, "mismatch scenario")
	assert.Equal(t, gameInput.PlayerName, output.PlayerName, "mismatch player name")

}

func TestNewQuestionValidScenario(t *testing.T) {
	// initialize new service
	s, err := newNewService()
	require.NoError(t, err)
	// initialize game
	game := Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "NEW_QUESTION",
		Score:        10,
		CountCorrect: 1,
		CurrentQuestion: &TimedQuestion{
			Question: Question{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			Timeout: 5,
		},
	}
	// insert new game to storage
	err = s.(*service).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)
	// test the new question function for game's scenario is NEW_QUESTION
	actualOutput, err := s.NewQuestion(context.Background(), NewQuestionInput{
		GameID: game.GameID,
	})
	require.NoError(t, err)

	// match the result
	assert.Equal(t, game.GameID, actualOutput.GameID, "mismatch game ID")
	assert.Equal(t, "SUBMIT_ANSWER", actualOutput.Scenario, "mismatch scenario")
	// based on the count correct of game, the current (for the new question)'s timeout will be 10
	assert.Equal(t, 10, actualOutput.Timeout, "mismatch timeout")
}

func TestNewQuestionInvalidScenario(t *testing.T) {
	// initialize new service
	s, err := newNewService()
	require.NoError(t, err)
	// initialize game
	game := Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "SUBMIT_ANSWER",
		Score:        10,
		CountCorrect: 1,
		CurrentQuestion: &TimedQuestion{
			Question: Question{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			Timeout: 5,
		},
	}
	// insert new game to storage
	err = s.(*service).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)
	// test the new question function for game's scenario is SUBMIT_ANSWER
	actualOutput, err := s.NewQuestion(context.Background(), NewQuestionInput{
		GameID: game.GameID,
	})
	// verify result
	assert.Error(t, err, "result is not error")
	assert.Nil(t, actualOutput, "the new question output is not nil")
	assert.Equal(t, ErrInvalidScenario, err, "mismatch error type")
}

func TestSubmitAnswerInvalidScenario(t *testing.T) {
	// initialize new service
	s, err := newNewService()
	require.NoError(t, err)
	// initialize game
	game := Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "NEW_QUESTION",
		Score:        10,
		CountCorrect: 1,
		CurrentQuestion: &TimedQuestion{
			Question: Question{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			Timeout: 5,
		},
	}
	// insert new game to storage
	err = s.(*service).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)
	// test the submit answer function for game's scenario is NEW_QUESTION
	nowTs := time.Now()
	input := SubmitAnswerInput{
		GameID:    game.GameID,
		AnswerIdx: 1,
		StartAt:   nowTs.Unix(),
		SentAt:    nowTs.Add(3 * time.Second).Unix(),
	}
	actualOutput, err := s.SubmitAnswer(context.Background(), input)

	// verify result
	assert.Error(t, err, "result is not error")
	assert.Nil(t, actualOutput, "the new question output is not nil")
	assert.Equal(t, ErrInvalidScenario, err, "mismatch error type")
}

func TestSubmitAnswerHasTimeout(t *testing.T) {
	// initialize new service
	s, err := newNewService()
	require.NoError(t, err)
	// initialize game
	game := Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "SUBMIT_ANSWER",
		Score:        10,
		CountCorrect: 1,
		CurrentQuestion: &TimedQuestion{
			Question: Question{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			Timeout: 5,
		},
	}
	// insert new game to storage
	err = s.(*service).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)
	// test the submit answer function for has timeout condition
	nowTs := time.Now()
	input := SubmitAnswerInput{
		GameID:    game.GameID,
		AnswerIdx: 1,
		StartAt:   nowTs.Unix(),
		SentAt:    nowTs.Add(7 * time.Second).Unix(),
	}
	actualOutput, err := s.SubmitAnswer(context.Background(), input)
	require.NoError(t, err)

	// verify result
	assert.Equal(t, game.GameID, actualOutput.GameID, "mismatch game ID")
	assert.Equal(t, "GAME_OVER", actualOutput.Scenario, "mismatch result scenario")
}

func TestSubmitAnswerValid(t *testing.T) {
	// initialize new service
	s, err := newNewService()
	require.NoError(t, err)
	// initialize game
	game := Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     "SUBMIT_ANSWER",
		Score:        10,
		CountCorrect: 1,
		CurrentQuestion: &TimedQuestion{
			Question: Question{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			Timeout: 5,
		},
	}
	// insert new game to storage
	err = s.(*service).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)
	// test the submit answer function for valid scenario
	nowTs := time.Now()
	input := SubmitAnswerInput{
		GameID:    game.GameID,
		AnswerIdx: 1,
		StartAt:   nowTs.Unix(),
		SentAt:    nowTs.Add(3 * time.Second).Unix(),
	}
	actualOutput, err := s.SubmitAnswer(context.Background(), input)
	require.NoError(t, err)

	// verify result
	assert.Equal(t, game.GameID, actualOutput.GameID, "mismatch game ID")
	assert.Equal(t, "NEW_QUESTION", actualOutput.Scenario, "mismatch game scenario")
}

type mockGameStorage struct {
	gameMap map[string]Game
}

func newMockGameStorage() *mockGameStorage {
	return &mockGameStorage{
		gameMap: map[string]Game{},
	}
}

func (gs *mockGameStorage) PutGame(ctx context.Context, g Game) error {
	gs.gameMap[g.GameID] = g
	return nil
}

func (gs *mockGameStorage) GetGame(ctx context.Context, gameID string) (*Game, error) {
	game := gs.gameMap[gameID]
	return &game, nil
}

type mockQuestionStorage struct {
	questions []Question
}

func newMockQuestionStorage(questions []Question) *mockQuestionStorage {
	return &mockQuestionStorage{
		questions: questions,
	}
}

func (qs *mockQuestionStorage) GetRandomQuestion(ctx context.Context) (*Question, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(qs.questions))

	return &qs.questions[idx], nil
}

type mockTimeoutCalculator struct {
	easyTimeout        int
	mediumTimeout      int
	hardTimeout        int
	mediumCorrectCount int
	hardCorrectCount   int
}

func newMockTimeoutCalculator() *mockTimeoutCalculator {
	return &mockTimeoutCalculator{
		easyTimeout:        10,
		mediumTimeout:      5,
		hardTimeout:        3,
		mediumCorrectCount: 10,
		hardCorrectCount:   25,
	}
}

func (tc *mockTimeoutCalculator) CalcTimeout(ctx context.Context, countCorrect int) int {
	if countCorrect > tc.hardCorrectCount {
		return tc.hardTimeout
	}
	if countCorrect > tc.mediumCorrectCount {
		return tc.mediumTimeout
	}
	return tc.easyTimeout
}

type mockClock struct {
}

func newMockClock() *mockClock {
	return &mockClock{}
}

func (c *mockClock) Now() time.Time {
	return time.Now()
}

func newNewService() (Service, error) {
	service, err := NewService(ServiceConfig{
		GameStorage: newMockGameStorage(),
		QuestionStorage: newMockQuestionStorage([]Question{
			{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			{
				Problem: "2 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 2,
			},
			{
				Problem: "3 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 3,
			},
		}),
		TimeoutCalculator: newMockTimeoutCalculator(),
		Clock:             newMockClock(),
		AddScore:          1,
	})
	return service, err
}
