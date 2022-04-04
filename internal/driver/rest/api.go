package rest

import (
	"errors"
	"net/http"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"gopkg.in/validator.v2"
)

type API struct {
	auth    core.Auth
	service core.Service
}

type APIConfig struct {
	Auth    core.Auth    `validate:"nonnil"`
	Service core.Service `validate:"nonnil"`
}

func (c APIConfig) Validate() error {
	return validator.Validate(c)
}

func NewAPI(cfg APIConfig) (*API, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	a := &API{
		auth:    cfg.Auth,
		service: cfg.Service,
	}
	return a, nil
}

func (a *API) GetHandler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", a.serveWeb)
	r.Get("/assets/*", a.serveAssets)
	r.Route("/games", func(r chi.Router) {
		r.Use(a.validateAPIKey)
		r.Use(render.SetContentType(render.ContentTypeJSON))

		r.Post("/", a.serveNewGame)
		r.Route("/{game_id}", func(r chi.Router) {
			r.Put("/question", a.serveNewQuestion)
			r.Put("/answer", a.serveSubmitAnswer)
		})
	})

	return r
}

func (a *API) serveWeb(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./web"))
	fs.ServeHTTP(w, r)
}

func (a *API) serveAssets(w http.ResponseWriter, r *http.Request) {
	fs := http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/assets")))
	fs.ServeHTTP(w, r)
}

func (a *API) serveNewGame(w http.ResponseWriter, r *http.Request) {
	// parse response body
	var rb newGameReqBody
	err := render.Bind(r, &rb)
	if err != nil {
		render.Render(w, r, newErrorResp(newBadRequestError(err.Error())))
		return
	}
	// create new game
	output, err := a.service.NewGame(r.Context(), core.NewGameInput{
		PlayerName: rb.PlayerName,
	})
	if err != nil {
		render.Render(w, r, newErrorResp(err))
		return
	}
	// output success response
	render.Render(w, r, newSuccessResp(output))
}

func (a *API) serveNewQuestion(w http.ResponseWriter, r *http.Request) {
	// initiate new question
	gameID := chi.URLParam(r, "game_id")
	output, err := a.service.NewQuestion(r.Context(), core.NewQuestionInput{GameID: gameID})
	if err != nil {
		switch err {
		case core.ErrGameNotFound:
			err = newGameNotFoundError()
		case core.ErrInvalidScenario:
			err = newInvalidScenarioError()
		}
		render.Render(w, r, newErrorResp(err))
		return
	}
	// output success response
	render.Render(w, r, newSuccessResp(output))
}

func (a *API) serveSubmitAnswer(w http.ResponseWriter, r *http.Request) {
	// parse response body
	var rb submitAnswerReqBody
	err := render.Bind(r, &rb)
	if err != nil {
		render.Render(w, r, newErrorResp(newBadRequestError(err.Error())))
		return
	}
	// submit answer
	gameID := chi.URLParam(r, "game_id")
	output, err := a.service.SubmitAnswer(r.Context(), core.SubmitAnswerInput{
		GameID:    gameID,
		AnswerIdx: rb.AnswerIndex,
		StartAt:   rb.StartAt,
		SentAt:    rb.SentAt,
	})
	if err != nil {
		switch err {
		case core.ErrGameNotFound:
			err = newGameNotFoundError()
		case core.ErrInvalidScenario:
			err = newInvalidScenarioError()
		}
		render.Render(w, r, newErrorResp(err))
		return
	}
	// output success response
	render.Render(w, r, newSuccessResp(output))
}

func (a *API) validateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		err := a.auth.ValidateAPIKey(r.Context(), apiKey)
		if err != nil {
			if errors.Is(err, core.ErrInvalidAPIKey) {
				err = newInvalidAPIKeyError()
			}
			render.Render(w, r, newErrorResp(err))
			return
		}
		next.ServeHTTP(w, r)
	})
}
