package queststrg_test

import (
	"context"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/redis/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/redis/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetRandomQuestion(t *testing.T) {
	client, err := testutils.InitRedisClient()
	require.NoError(t, err)
	err = testutils.ResetRedis(client)
	require.NoError(t, err)
	questions := []core.Question{
		{
			Problem:      "1 + 2",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 3,
		},
		{
			Problem:      "3 - 2",
			Choices:      []string{"1", "2", "3"},
			CorrectIndex: 1,
		},
	}
	err = testutils.InsertQuestion(client, questions)
	require.NoError(t, err)

	//get random
	storage, err := queststrg.New(queststrg.Config{RedisClient: client})
	require.NoError(t, err)
	restQuestion, err := storage.GetRandomQuestion(context.Background())
	require.NoError(t, err)
	require.Contains(t, questions, *restQuestion)
}
