package rest

import (
	"errors"
	"fmt"
	"net/http"
)

type apiError struct {
	StatusCode int
	Err        string
	Message    string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%v - %v - %v", e.StatusCode, e.Err, e.Message)
}

func (e *apiError) Is(target error) bool {
	var restErr *apiError
	if !errors.As(target, &restErr) {
		return false
	}
	return *e == *restErr
}

func newInternalServerError(msg string) *apiError {
	return &apiError{
		StatusCode: http.StatusInternalServerError,
		Err:        "ERR_INTERNAL_ERROR",
		Message:    msg,
	}
}

func newBadRequestError(msg string) *apiError {
	return &apiError{
		StatusCode: http.StatusBadRequest,
		Err:        "ERR_BAD_REQUEST",
		Message:    msg,
	}
}

func newGameNotFoundError() *apiError {
	return &apiError{
		StatusCode: http.StatusNotFound,
		Err:        "ERR_GAME_NOT_FOUND",
		Message:    "game is not found",
	}
}

func newInvalidScenarioError() *apiError {
	return &apiError{
		StatusCode: http.StatusConflict,
		Err:        "ERR_INVALID_SCENARIO",
		Message:    "invalid scenario for the action",
	}
}

func newInvalidAPIKeyError() *apiError {
	return &apiError{
		StatusCode: http.StatusUnauthorized,
		Err:        "ERR_INVALID_API_KEY",
		Message:    "invalid api key",
	}
}
