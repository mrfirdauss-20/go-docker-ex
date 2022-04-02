package rest

import (
	"fmt"
	"net/http"

	"gopkg.in/validator.v2"
)

type newGameReqBody struct {
	PlayerName string `json:"player_name" validate:"min=4"`
}

func (rb *newGameReqBody) Bind(r *http.Request) error {
	return validator.Validate(rb)
}

type submitAnswerReqBody struct {
	AnswerIndex int   `json:"answer_idx" validate:"min=0"`
	StartAt     int64 `json:"start_at" validate:"min=0"`
	SentAt      int64 `json:"sent_at" validate:"min=0"`
}

func (rb *submitAnswerReqBody) Bind(r *http.Request) error {
	if rb.SentAt < rb.StartAt {
		return fmt.Errorf("invalid value of `sent_at`")
	}
	return validator.Validate(rb)
}
