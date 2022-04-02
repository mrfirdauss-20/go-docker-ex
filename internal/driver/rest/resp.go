package rest

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type respBody struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	Err        string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
	Timestamp  int64       `json:"ts"`
}

func (rb *respBody) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rb.StatusCode)
	rb.Timestamp = time.Now().Unix()
	return nil
}

func newSuccessResp(data interface{}) *respBody {
	return &respBody{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
	}
}

func newErrorResp(err error) *respBody {
	var restErr *apiError
	if !errors.As(err, &restErr) {
		restErr = newInternalServerError(err.Error())
	}
	return &respBody{
		StatusCode: restErr.StatusCode,
		OK:         false,
		Err:        restErr.Err,
		Message:    restErr.Message,
	}
}
