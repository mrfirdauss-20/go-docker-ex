package testutils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"

	_ "github.com/go-sql-driver/mysql"
)

const envKeySQLTestDSN = "SQL_TEST_DSN"

func InitSQLClient() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", os.Getenv(envKeySQLTestDSN))
}

func ResetTables(sqlClient *sqlx.DB) error {
	query := `
		SET FOREIGN_KEY_CHECKS=0; 
		TRUNCATE TABLE questions; 
		TRUNCATE TABLE games; 
		SET FOREIGN_KEY_CHECKS=1;
	`
	_, err := sqlClient.Exec(query)
	if err != nil {
		return fmt.Errorf("unable to execute query due: %w", err)
	}
	return nil
}

func InsertQuestions(sqlClient *sqlx.DB, questions []core.Question) error {
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
