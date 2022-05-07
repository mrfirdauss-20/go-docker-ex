package gamestrg_test

import (
	"context"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/redis/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/redis/testutils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPutGameGetGame(t *testing.T) {
	// initialize redis client
	redisClient, err := testutils.InitRedisClient()
	require.NoError(t, err)

	// initialize game storage
	gs, err := gamestrg.New(gamestrg.Config{RedisClient: redisClient})
	require.NoError(t, err)

	// define question
	question := core.Question{
		Problem:      "1 + 1",
		Choices:      []string{"1", "2", "3"},
		CorrectIndex: 2,
	}

	// define test cases
	testCases := []struct {
		Name string
		Game core.Game
	}{
		{
			Name: "Test Nil Question",
			Game: core.Game{
				GameID:          uuid.NewString(),
				PlayerName:      "Risqi",
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
				PlayerName:   "Risqi",
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

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// put game
			ctx := context.Background()
			err := gs.PutGame(ctx, tc.Game)
			require.NoError(t, err)
			// get game
			game, err := gs.GetGame(ctx, tc.Game.GameID)
			require.NoError(t, err)
			// validate game data
			require.Equal(t, tc.Game, *game)
		})
	}
}
