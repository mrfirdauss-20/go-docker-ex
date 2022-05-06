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
	// initialize redis client
	client, err := testutils.InitRedisClient()
	require.NoError(t, err)

	// reset redis to make sure the data is clean
	err = testutils.ResetRedis(client)
	require.NoError(t, err)

	// insert questions to redis
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

	// initialize question storage
	storage, err := queststrg.New(queststrg.Config{RedisClient: client})
	require.NoError(t, err)

	// get random question
	resQuestion, err := storage.GetRandomQuestion(context.Background())
	require.NoError(t, err)

	// make sure the returned question is from the questions we have inserted
	require.Contains(t, questions, *resQuestion)
}
