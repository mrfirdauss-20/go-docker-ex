package queststrg_test

import (
	"context"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetRandomQuestion(t *testing.T) {
	// initialize sql client
	sqlClient, err := testutils.InitSQLClient()
	require.NoError(t, err)
	// reset question table for clean slate
	err = testutils.ResetTables(sqlClient)
	require.NoError(t, err)
	// insert several questions
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
	err = testutils.InsertQuestions(sqlClient, questions)
	require.NoError(t, err)
	// execute get random question
	storage, err := queststrg.New(queststrg.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	restQuestion, err := storage.GetRandomQuestion(context.Background())
	require.NoError(t, err)
	// check if returned question is from inserted questions
	require.Contains(t, questions, *restQuestion)
}
