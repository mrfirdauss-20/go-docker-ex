package queststrg_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/queststrg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
)

const envKeySQLTestDSN = "SQL_TEST_DSN"

func TestGetRandomQuestion(t *testing.T) {
	// initialize sql client
	sqlClient, err := sqlx.Connect("mysql", os.Getenv(envKeySQLTestDSN))
	require.NoError(t, err)
	// reset question table for clean slate
	err = resetQuestionTable(sqlClient)
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
	err = insertQuestions(sqlClient, questions)
	require.NoError(t, err)
	// execute get random question
	storage, err := queststrg.New(queststrg.Config{SQLClient: sqlClient})
	require.NoError(t, err)
	restQuestion, err := storage.GetRandomQuestion(context.Background())
	require.NoError(t, err)
	// check if returned question is from inserted questions
	require.Contains(t, questions, *restQuestion)
}

func resetQuestionTable(sqlClient *sqlx.DB) error {
	query := `SET FOREIGN_KEY_CHECKS=0; TRUNCATE TABLE questions;`
	_, err := sqlClient.Exec(query)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}

func insertQuestions(sqlClient *sqlx.DB, questions []core.Question) error {
	var qRows []questionRow
	for _, question := range questions {
		qRows = append(qRows, newQuestionRow(question))
	}
	query := `
		INSERT INTO questions (problem, correct_index, answers) VALUES (
			:problem, :correct_index, :answers
		)
	`
	_, err := sqlClient.NamedExec(query, qRows)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}

type questionRow struct {
	Problem      string         `db:"problem"`
	CorrectIndex int            `db:"correct_index"`
	Answers      types.JSONText `db:"answers"`
}

func newQuestionRow(q core.Question) questionRow {
	b, _ := json.Marshal(q.Choices)
	return questionRow{
		Problem:      q.Problem,
		CorrectIndex: q.CorrectIndex,
		Answers:      types.JSONText(b),
	}
}
