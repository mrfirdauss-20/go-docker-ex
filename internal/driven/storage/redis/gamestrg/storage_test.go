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
	redisClient, err := testutils.InitRedisClient()
	require.NoError(t, err)
	err = testutils.ResetRedis(redisClient)
	require.NoError(t, err)

	question := core.Question{
		Problem:      "1 + 1",
		Choices:      []string{"1", "2", "3"},
		CorrectIndex: 2,
	}

	err = redisClient.Set(context.Background(), "question", question, 0).Err()
	require.NoError(t, err)

	gs, err := gamestrg.New(gamestrg.Config{RedisClient: redisClient})
	require.NoError(t, err)

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
			ctx := context.Background()
			err := gs.PutGame(ctx, tc.Game)
			require.NoError(t, err)

			game, err := gs.GetGame(ctx, tc.Game.GameID)
			require.NoError(t, err)
			require.Equal(t, tc.Game, game)
		})
	}
}
