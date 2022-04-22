package gamestrg_test

import (
	"context"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPutGameGetGame(t *testing.T) {
	// initialize sql client
	sqlClient, err := testutils.InitSQLClient()
	require.NoError(t, err)
	// reset tables for clean slate
	err = testutils.ResetTables(sqlClient)
	require.NoError(t, err)
	// insert sample question
	question := core.Question{
		Problem:      "1 + 1",
		Choices:      []string{"1", "2", "3"},
		CorrectIndex: 2,
	}
	_, err = testutils.InsertQuestion(sqlClient, question)
	require.NoError(t, err)
	// initialize game storage
	gameStorage, err := gamestrg.New(gamestrg.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	// define game objects
	testCases := []struct {
		Name string
		Game core.Game
	}{
		{
			Name: "Test Nil Question",
			Game: core.Game{
				GameID:          uuid.NewString(),
				PlayerName:      "Riandy",
				Scenario:        core.ScenarioNewQuestion,
				Score:           0,
				CountCorrect:    0,
				CurrentQuestion: nil,
			},
		},
		{
			Name: "Test Non-Nil Question",
			Game: core.Game{
				GameID:       uuid.NewString(),
				PlayerName:   "Riandy",
				Scenario:     core.ScenarioNewQuestion,
				Score:        0,
				CountCorrect: 0,
				CurrentQuestion: &core.TimedQuestion{
					Question: question,
					Timeout:  5,
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// put object
			ctx := context.Background()
			err := gameStorage.PutGame(ctx, testCase.Game)
			require.NoError(t, err)
			// get object
			game, err := gameStorage.GetGame(ctx, testCase.Game.GameID)
			require.NoError(t, err)
			// validate result
			require.Equal(t, testCase.Game, *game)
		})
	}
}
