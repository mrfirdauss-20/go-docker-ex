package core

import "errors"

var (
	ErrGameNotFound    = errors.New("game is not found")
	ErrInvalidScenario = errors.New("invalid scenario")
	ErrInvalidAPIKey   = errors.New("invalid api key")
)
