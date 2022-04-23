package queststrg

import (
	"fmt"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/jmoiron/sqlx/types"
)

type questionRow struct {
	Problem      string         `db:"problem"`
	CorrectIndex int            `db:"correct_index"`
	Answers      types.JSONText `db:"answers"`
}

func (r questionRow) toQuestion() (*core.Question, error) {
	var choices []string
	err := r.Answers.Unmarshal(&choices)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal choices due: %w", err)
	}
	q := &core.Question{
		Problem:      r.Problem,
		CorrectIndex: r.CorrectIndex,
		Choices:      choices,
	}
	return q, nil
}
