package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/clock"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/memory/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/memory/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/toutcalc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const API_KEY = "c4211664-47dc-4887-a2fe-9e694fbaf55a"

var questions = []core.Question{
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
}

func TestNewAPI(t *testing.T) {
	// define mock dependencies
	auth := newMockAuth(API_KEY)
	service, err := newMockService()
	require.NoError(t, err)
	// define function for validating new API
	validateAPI := func(t *testing.T, a API, cfg APIConfig) {
		assert.Equal(t, cfg.Auth, a.auth)
		assert.Equal(t, cfg.Service, a.service)
	}
	// define test cases
	testCases := []struct {
		Name    string
		Config  APIConfig
		IsError bool
	}{
		{
			Name: "Test New API Missing Auth",
			Config: APIConfig{
				Auth:    nil,
				Service: service,
			},
			IsError: true,
		},
		{
			Name: "Test New API Missing Service",
			Config: APIConfig{
				Auth:    auth,
				Service: nil,
			},
			IsError: true,
		},
		{
			Name: "Test New API Valid",
			Config: APIConfig{
				Auth:    auth,
				Service: service,
			},
			IsError: false,
		},
	}
	// execute test cases
	for i := range testCases {
		t.Run(testCases[i].Name, func(t *testing.T) {
			api, err := NewAPI(testCases[i].Config)
			assert.Equal(t, testCases[i].IsError, (err != nil), "unexpected error")
			if api == nil {
				return
			}
			validateAPI(t, *api, testCases[i].Config)
		})
	}
}

func TestServeNewGameInvalid(t *testing.T) {
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)

	// create request
	reqBody := newGameReqBody{}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/games/", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveNewGame handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build expected response body
	expectedResp := respBody{
		OK:         false,
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
	}
	// build actual response body
	actualResp := respBody{}
	err = json.NewDecoder(resp.Body).Decode(&actualResp)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, expectedResp.StatusCode, resp.StatusCode, "mismatch response code")
	assert.Equal(t, expectedResp.OK, actualResp.OK, "mismatch response body ok")
	assert.Equal(t, expectedResp.Err, actualResp.Err, "mismatch response body error type")
}

func TestServeNewGameValid(t *testing.T) {
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)

	// create request
	reqBody := newGameReqBody{
		PlayerName: "Riandy",
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/games/", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveNewGame handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build expected response body
	expectedResp := newGameRespBody{
		Ok: true,
		Data: core.NewGameOutput{
			PlayerName: "Riandy",
			Scenario:   core.ScenarioNewQuestion,
		},
	}
	// build actual response body
	actualResp := newGameRespBody{}
	err = json.NewDecoder(resp.Body).Decode(&actualResp)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "mismatch response code")
	assert.Equal(t, expectedResp.Ok, actualResp.Ok, "mismatch response body ok")
	assert.Equal(t, expectedResp.Data.PlayerName, actualResp.Data.PlayerName, "mismatch response body player name")
	assert.Equal(t, expectedResp.Data.Scenario, actualResp.Data.Scenario, "mismatch response body game scenario")

	// make sure that storage save the game
	storedGameID := actualResp.Data.GameID
	storedGame, err := api.service.(*mockService).gameStorage.GetGame(context.Background(), storedGameID)
	assert.NoError(t, err)
	assert.Equal(t, storedGameID, storedGame.GameID, "mismatch game ID")
	assert.Equal(t, expectedResp.Data.PlayerName, storedGame.PlayerName, "mismatch stored game player name")
	assert.Equal(t, expectedResp.Data.Scenario, storedGame.Scenario, "mismatch stored game scenario")
}

func TestServeNewQuestionValid(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:          uuid.NewString(),
		PlayerName:      "Riandy",
		Scenario:        core.ScenarioNewQuestion,
		Score:           0,
		CountCorrect:    0,
		CurrentQuestion: &core.TimedQuestion{},
	}
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/question", nil)
	// set headers, here we set them using `Set()` because it will normalize the header key value
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the handler, here test the exported http handler because if we test the handler function
	// directly, chi context would be missing which resulted in `chi.URLParams(r, key)` returns empty.
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody newQuestionRespBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "mismatch response code")
	assert.Equal(t, true, respBody.Ok, "mismatch response body ok")
	assert.Equal(t, game.GameID, respBody.Data.GameID, "mismatch response body game ID")
	assert.Equal(t, core.ScenarioSubmitAnswer, respBody.Data.Scenario, "mismatch response body game scenario")

	// make sure that problem and choices in response body is part of question stored in queststrg
	question := core.Question{
		Problem: respBody.Data.Problem,
		Choices: respBody.Data.Choices,
	}
	questInQuestions := func(questResponse core.Question, questions []core.Question) bool {
		for i := range questions {
			if questResponse.Problem == questions[i].Problem {
				for j := range questResponse.Choices {
					if questResponse.Choices[j] != questions[i].Choices[j] {
						return false
					}
				}
				return true
			}
		}
		return false
	}
	assert.True(t, questInQuestions(question, questions), "current question is out of storage")
}

func TestServeNewQuestionInvalid(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:          uuid.NewString(),
		PlayerName:      "Riandy",
		Scenario:        core.ScenarioGameOver,
		Score:           0,
		CountCorrect:    0,
		CurrentQuestion: &core.TimedQuestion{},
	}
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/question", nil)
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveNewQuestion handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody respBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "mismatch response code")
	assert.Equal(t, false, respBody.OK, "mismatch response body ok")
	assert.Equal(t, "ERR_INVALID_SCENARIO", respBody.Err, "mismatch response body error type")
}

func TestServeSubmitAnswerCorrectAnswer(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     core.ScenarioSubmitAnswer,
		Score:        0,
		CountCorrect: 0,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
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
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	reqBody := submitAnswerReqBody{
		AnswerIndex: 1,
		StartAt:     time.Now().Unix(),
		SentAt:      time.Now().Add(3 * time.Second).Unix(),
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/answer", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveSubmitAnswer handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody submitAnswerRespBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	newScore := game.Score + api.service.(*mockService).addScore
	assert.Equal(t, http.StatusOK, resp.StatusCode, "mismatch response code")
	assert.Equal(t, true, respBody.Ok, "mismatch response body ok")
	assert.Equal(t, game.GameID, respBody.Data.GameID, "mismatch response body game ID")
	assert.Equal(t, core.ScenarioNewQuestion, respBody.Data.Scenario, "mismatch response body game scenario")
	assert.Equal(t, reqBody.AnswerIndex, respBody.Data.AnswerIndex, "mismatch response body answer index")
	assert.Equal(t, game.CurrentQuestion.CorrectIndex, respBody.Data.CorrectIndex, "mismatch response body correct index")
	assert.Equal(t, int(reqBody.SentAt-reqBody.StartAt), respBody.Data.Duration, "mismatch response body duration")
	assert.Equal(t, game.CurrentQuestion.Timeout, respBody.Data.Timeout, "mismatch response body timeout")
	assert.Equal(t, newScore, respBody.Data.Score, "mismatch response body score")
	// make sure that game is updated in gamestrg
	updatedGame, err := api.service.(*mockService).gameStorage.GetGame(context.Background(), game.GameID)
	assert.Equal(t, core.ScenarioNewQuestion, updatedGame.Scenario, "mismatch updated game scenario")
	assert.Equal(t, (game.CountCorrect + 1), updatedGame.CountCorrect, "mismatch updated game count correct")
	assert.Equal(t, newScore, updatedGame.Score, "mismatch updated game score")
}

func TestServeSubmitAnswerIncorrectAnswer(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     core.ScenarioSubmitAnswer,
		Score:        0,
		CountCorrect: 0,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
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
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	reqBody := submitAnswerReqBody{
		AnswerIndex: 2,
		StartAt:     time.Now().Unix(),
		SentAt:      time.Now().Add(3 * time.Second).Unix(),
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/answer", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveSubmitAnswer handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody submitAnswerRespBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "mismatch response code")
	assert.Equal(t, true, respBody.Ok, "mismatch response body ok")
	assert.Equal(t, game.GameID, respBody.Data.GameID, "mismatch response body game ID")
	assert.Equal(t, core.ScenarioGameOver, respBody.Data.Scenario, "mismatch response body game scenario")
	assert.Equal(t, reqBody.AnswerIndex, respBody.Data.AnswerIndex, "mismatch response body answer index")
	assert.Equal(t, game.CurrentQuestion.CorrectIndex, respBody.Data.CorrectIndex, "mismatch response body correct index")
	assert.Equal(t, int(reqBody.SentAt-reqBody.StartAt), respBody.Data.Duration, "mismatch response body duration")
	assert.Equal(t, game.CurrentQuestion.Timeout, respBody.Data.Timeout, "mismatch response body timeout")
	assert.Equal(t, game.Score, respBody.Data.Score, "mismatch response body score")
	// make sure that game is updated in gamestrg
	updatedGame, err := api.service.(*mockService).gameStorage.GetGame(context.Background(), game.GameID)
	assert.Equal(t, core.ScenarioGameOver, updatedGame.Scenario, "mismatch updated game scenario")
}

func TestServeSubmitAnswerHasTimeout(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     core.ScenarioSubmitAnswer,
		Score:        0,
		CountCorrect: 0,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
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
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	reqBody := submitAnswerReqBody{
		AnswerIndex: 1,
		StartAt:     time.Now().Unix(),
		SentAt:      time.Now().Add(7 * time.Second).Unix(),
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/answer", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveSubmitAnswer handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody submitAnswerRespBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusOK, resp.StatusCode, "mismatch response code")
	assert.Equal(t, true, respBody.Ok, "mismatch response body ok")
	assert.Equal(t, game.GameID, respBody.Data.GameID, "mismatch response body game ID")
	assert.Equal(t, core.ScenarioGameOver, respBody.Data.Scenario, "mismatch response body game scenario")
	assert.Equal(t, reqBody.AnswerIndex, respBody.Data.AnswerIndex, "mismatch response body answer index")
	assert.Equal(t, game.CurrentQuestion.CorrectIndex, respBody.Data.CorrectIndex, "mismatch response body correct index")
	assert.Equal(t, int(reqBody.SentAt-reqBody.StartAt), respBody.Data.Duration, "mismatch response body duration")
	assert.Equal(t, game.CurrentQuestion.Timeout, respBody.Data.Timeout, "mismatch response body timeout")
	assert.Equal(t, game.Score, respBody.Data.Score, "mismatch response body score")
	// make sure that game is updated in gamestrg
	updatedGame, err := api.service.(*mockService).gameStorage.GetGame(context.Background(), game.GameID)
	assert.Equal(t, core.ScenarioGameOver, updatedGame.Scenario, "mismatch updated game scenario")
}

func TestServeSubmitAnswerInvalidScenario(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     core.ScenarioNewQuestion,
		Score:        0,
		CountCorrect: 0,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
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
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	reqBody := submitAnswerReqBody{
		AnswerIndex: 1,
		StartAt:     time.Now().Unix(),
		SentAt:      time.Now().Add(4 * time.Second).Unix(),
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/answer", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveSubmitAnswer handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody respBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusConflict, resp.StatusCode, "mismatch response code")
	assert.Equal(t, false, respBody.OK, "mismatch response body ok")
	assert.Equal(t, "ERR_INVALID_SCENARIO", respBody.Err, "mismatch response body error type")
}

func TestServeSubmitAnswerInvalidRequest(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:       uuid.NewString(),
		PlayerName:   "Riandy",
		Scenario:     core.ScenarioNewQuestion,
		Score:        0,
		CountCorrect: 0,
		CurrentQuestion: &core.TimedQuestion{
			Question: core.Question{
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
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)
	// store game to service
	err = api.service.(*mockService).gameStorage.PutGame(context.Background(), game)
	require.NoError(t, err)

	// create request
	reqBody := submitAnswerReqBody{
		AnswerIndex: 1,
		StartAt:     time.Now().Unix(),
	}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/games/"+game.GameID+"/answer", strings.NewReader(string(reqBodyJson)))
	req.Header.Set("X-API-Key", API_KEY)
	req.Header.Set("Content-Type", "application/json")
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveSubmitAnswer handler function
	api.GetHandler().ServeHTTP(w, req)

	// read the response body
	resp := w.Result()
	defer resp.Body.Close()
	// build actual response body
	var respBody respBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "mismatch response code")
	assert.Equal(t, false, respBody.OK, "mismatch response body ok")
	assert.Equal(t, "ERR_BAD_REQUEST", respBody.Err, "mismatch response body error type")
}

type newGameRespBody struct {
	Ok   bool               `json:"ok"`
	Data core.NewGameOutput `json:"data"`
	Time int64              `json:"ts"`
}

type newQuestionRespBody struct {
	Ok   bool                   `json:"ok"`
	Data core.NewQuestionOutput `json:"data"`
	Time int64                  `json:"ts"`
}

type submitAnswerRespBody struct {
	Ok   bool                    `json:"ok"`
	Data core.SubmitAnswerOutput `json:"data"`
	Time int64                   `json:"ts"`
}

// mock Auth
type mockAuth struct {
	apiKey string
}

func newMockAuth(key string) *mockAuth {
	return &mockAuth{
		apiKey: key,
	}
}

func (a *mockAuth) ValidateAPIKey(ctx context.Context, apiKey string) error {
	if a.apiKey != apiKey {
		return core.ErrInvalidAPIKey
	}
	return nil
}

// mock Service
type mockService struct {
	gameStorage       gamestrg.Storage
	questionStorage   queststrg.Storage
	timeoutCalculator toutcalc.Calculator
	clock             clock.Clock
	addScore          int
}

func newMockService() (*mockService, error) {
	gameStrg := gamestrg.New()
	questStrg, err := queststrg.New(queststrg.Config{
		Questions: questions,
	})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	toutCalc, err := toutcalc.New(toutcalc.StandardConfig())
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	clock := clock.New()
	s := mockService{
		gameStorage:       *gameStrg,
		questionStorage:   *questStrg,
		timeoutCalculator: *toutCalc,
		clock:             *clock,
		addScore:          1,
	}
	return &s, nil
}

func (s *mockService) NewGame(ctx context.Context, input core.NewGameInput) (*core.NewGameOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	game := core.Game{
		GameID:          uuid.NewString(),
		PlayerName:      input.PlayerName,
		Scenario:        "NEW_QUESTION",
		Score:           0,
		CountCorrect:    0,
		CurrentQuestion: nil,
	}
	err = s.gameStorage.PutGame(ctx, game)
	if err != nil {
		return nil, fmt.Errorf("unable to create game in storage due: %w", err)
	}

	output := &core.NewGameOutput{
		GameID:     game.GameID,
		PlayerName: game.PlayerName,
		Scenario:   game.Scenario,
	}
	return output, nil
}

func (s *mockService) NewQuestion(ctx context.Context, input core.NewQuestionInput) (*core.NewQuestionOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	game, err := s.gameStorage.GetGame(ctx, input.GameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game instance due: %w", err)
	}
	if game == nil {
		return nil, core.ErrGameNotFound
	}
	if game.Scenario != "NEW_QUESTION" {
		return nil, core.ErrInvalidScenario
	}

	question, err := s.questionStorage.GetRandomQuestion(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get random question due: %w", err)
	}

	game.CurrentQuestion = &core.TimedQuestion{
		Question: *question,
		Timeout:  s.timeoutCalculator.CalcTimeout(ctx, game.CountCorrect),
	}
	game.Scenario = "SUBMIT_ANSWER"
	err = s.gameStorage.PutGame(ctx, *game)
	if err != nil {
		return nil, fmt.Errorf("unable to update game storage due: %w", err)
	}

	output := &core.NewQuestionOutput{
		GameID:   game.GameID,
		Scenario: game.Scenario,
		Problem:  game.CurrentQuestion.Problem,
		Choices:  game.CurrentQuestion.Choices,
		Timeout:  game.CurrentQuestion.Timeout,
	}
	return output, nil
}

func (s *mockService) SubmitAnswer(ctx context.Context, input core.SubmitAnswerInput) (*core.SubmitAnswerOutput, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	game, err := s.gameStorage.GetGame(ctx, input.GameID)
	if err != nil {
		return nil, fmt.Errorf("unable to get game instance due: %w", err)
	}
	if game == nil {
		return nil, core.ErrGameNotFound
	}
	if game.Scenario != "SUBMIT_ANSWER" {
		return nil, core.ErrInvalidScenario
	}
	duration := int(input.SentAt - input.StartAt)
	isTimeout := duration >= game.CurrentQuestion.Timeout
	isIncorrect := input.AnswerIdx != game.CurrentQuestion.CorrectIndex
	if isTimeout || isIncorrect {
		game.Scenario = "GAME_OVER"
	} else {
		game.Scenario = "NEW_QUESTION"
		game.CountCorrect += 1
		game.Score += s.addScore
	}

	err = s.gameStorage.PutGame(ctx, *game)
	if err != nil {
		return nil, fmt.Errorf("unable to update game in storage due: %w", err)
	}

	output := &core.SubmitAnswerOutput{
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

func newNewAPI() (*API, error) {
	auth := newMockAuth(API_KEY)
	service, err := newMockService()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	api, err := NewAPI(APIConfig{
		Auth:    auth,
		Service: service,
	})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return api, nil
}
