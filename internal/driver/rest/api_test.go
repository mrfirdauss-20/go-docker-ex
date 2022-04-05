package rest

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/clock"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/memory/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/memory/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/toutcalc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeWeb(t *testing.T) {
	// initialize auth
	auth, err := core.NewAuth(core.AuthConfig{
		APIKey: "c4211664-47dc-4887-a2fe-9e694fbaf55a",
	})
	require.NoError(t, err)
	// initialize service
	gameStrg := gamestrg.New()
	questStrg, err := queststrg.New(queststrg.Config{
		Questions: []core.Question{
			{
				Problem: "1 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 1,
			},
			{
				Problem: "2 + 2",
				Choices: []string{
					"3",
					"4",
					"5",
				},
				CorrectIndex: 2,
			},
		},
	})
	require.NoError(t, err)
	timeoutCalculator, err := toutcalc.New(toutcalc.StandardConfig())
	require.NoError(t, err)
	clock := clock.New()
	service, err := core.NewService(core.ServiceConfig{
		GameStorage:       gameStrg,
		QuestionStorage:   questStrg,
		TimeoutCalculator: timeoutCalculator,
		Clock:             clock,
		AddScore:          1,
	})
	require.NoError(t, err)
	// initialize api
	api, err := NewAPI(APIConfig{
		Auth:    auth,
		Service: service,
	})
	require.NoError(t, err)
	// initialize route
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", api.serveWeb)
	// test the serveWeb function
	ts := httptest.NewServer(r)
	ts.URL = "http://localhost:9190"
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/", nil)
	// verify output
	assert.Equal(t, http.StatusOK, resp.StatusCode, "response status is not OK")
	byteFile, err := os.ReadFile("../../../cmd/mem_server/web/index.html")
	require.NoError(t, err)
	assert.Equal(t, string(byteFile), body, "mismatch response body")
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()

	return resp, string(respBody)
}
