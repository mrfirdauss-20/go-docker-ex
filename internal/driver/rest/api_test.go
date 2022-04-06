package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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

// this test function still remains fail
func TestServeWeb(t *testing.T) {
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)

	// create request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	// test the serveWeb handler function
	api.serveWeb(w, req)
	res := w.Result()
	defer res.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	// verify result
	assert.Equal(t, http.StatusOK, res.StatusCode, "mismatch response code")
	// read homepage file
	file, err := os.ReadFile("../../../cmd/mem_server/web/index.html")
	require.NoError(t, err)
	// verify body
	assert.Equal(t, string(file), string(body))
}

// this test function still remains fail
func TestServeAssets(t *testing.T) {
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)

	// create request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/assets/", nil)
	w := httptest.NewRecorder()
	// test the serveWeb handler function
	api.serveAssets(w, req)
	res := w.Result()
	defer res.Body.Close()

	// read the response body
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	// verify result
	assert.Equal(t, http.StatusOK, res.StatusCode, "mismatch response code")

	expectedBody := `<pre>
	<a href="css/">css/</a>
	<a href="js/">js/</a>
	</pre>`
	assert.Equal(t, expectedBody, string(body))
}

func TestServeNewGameInvalid(t *testing.T) {
	// initialize new API
	api, err := newNewAPI()
	require.NoError(t, err)

	// create request
	reqBody := newGameReqBody{}
	reqBodyJson, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/games/", strings.NewReader(string(reqBodyJson)))
	req.Header = map[string][]string{
		"X-API-Key":    {API_KEY},
		"Content-Type": {"application/json"},
	}
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveNewGame handler function
	api.serveNewGame(w, req)

	// read the response body
	res := w.Result()
	defer res.Body.Close()
	// build expected response body
	expectedResp := respBody{
		OK:         false,
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
		// Message:    "missing `player_name`",
	}
	// build actual response body
	byteData, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	actualResp := respBody{}
	err = json.Unmarshal(byteData, &actualResp)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, expectedResp.StatusCode, res.StatusCode, "mismatch response code")
	assert.Equal(t, expectedResp.OK, actualResp.OK)
	assert.Equal(t, expectedResp.Err, actualResp.Err)
	// assert.Equal(t, expectedResp.Message, actualResp.Message)
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
	// require.NoError(t, err)
	req.Header = map[string][]string{
		"X-API-Key":    {API_KEY},
		"Content-Type": {"application/json"},
	}
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveNewGame handler function
	api.serveNewGame(w, req)

	// read the response body
	res := w.Result()
	defer res.Body.Close()
	// build expected response body
	expectedResp := newGameResp{
		Ok: true,
		Data: core.NewGameOutput{
			PlayerName: "Riandy",
			Scenario:   "NEW_QUESTION",
		},
	}
	// build actual response body
	byteData, _ := ioutil.ReadAll(res.Body)
	actualResp := newGameResp{}
	err = json.Unmarshal(byteData, &actualResp)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusOK, res.StatusCode, "mismatch response code")
	assert.Equal(t, expectedResp.Ok, actualResp.Ok)
	assert.Equal(t, expectedResp.Data.PlayerName, actualResp.Data.PlayerName, "mismatch response body player name")
	assert.Equal(t, expectedResp.Data.Scenario, actualResp.Data.Scenario, "mismatch response body player name")

	// make sure that storage save the game
	storedGameID := actualResp.Data.GameID
	storedGame, err := api.service.(*mockService).gameStorage.GetGame(context.Background(), storedGameID)
	assert.NoError(t, err)
	assert.Equal(t, storedGameID, storedGame.GameID, "mismatch game ID")
	assert.Equal(t, expectedResp.Data.PlayerName, storedGame.PlayerName, "mismatch response body player name")
	assert.Equal(t, expectedResp.Data.Scenario, storedGame.Scenario, "mismatch response body player name")
}

func TestServeNewQuestionValid(t *testing.T) {
	// initialize game
	game := core.Game{
		GameID:          uuid.NewString(),
		PlayerName:      "Riandy",
		Scenario:        "NEW_QUESTION",
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
	req.Header = map[string][]string{
		"X-API-Key":    {API_KEY},
		"Content-Type": {"application/json"},
	}
	log.Println(req)
	// create response recorder
	w := httptest.NewRecorder()
	// test the serveNewGame handler function
	api.serveNewQuestion(w, req)

	// read the response body
	res := w.Result()
	defer res.Body.Close()
	// build expected response ok status
	expectedRespOk := true
	// build actual response body
	byteData, _ := ioutil.ReadAll(res.Body)
	log.Println(string(byteData))
	actualResp := newQuestionResp{}
	err = json.Unmarshal(byteData, &actualResp)
	require.NoError(t, err)
	// verify response body
	assert.Equal(t, http.StatusOK, res.StatusCode, "mismatch response code")
	assert.Equal(t, expectedRespOk, actualResp.Ok)
	assert.Equal(t, game.GameID, actualResp.Data.GameID)
	assert.Equal(t, "SUBMIT_ANSWER", actualResp.Data.Scenario)

	// make sure that problem in response body is part of problem stored in queststrg
	questResponse := core.Question{
		Problem: actualResp.Data.Problem,
		Choices: actualResp.Data.Choices,
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
	assert.True(t, questInQuestions(questResponse, questions))
}

type newGameResp struct {
	Ok   bool               `json:"ok"`
	Data core.NewGameOutput `json:"data"`
	Time int64              `json:"ts"`
}

type newQuestionResp struct {
	Ok   bool                   `json:"ok"`
	Data core.NewQuestionOutput `json:"data"`
	Time int64                  `json:"ts"`
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
