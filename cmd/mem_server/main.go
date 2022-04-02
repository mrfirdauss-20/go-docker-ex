package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/clock"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/memory/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/memory/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/toutcalc"
	"github.com/ghazlabs/hex-mathrush/internal/driver/rest"
)

const (
	listenPort = 9190
	apiKey     = "c4211664-47dc-4887-a2fe-9e694fbaf55a"
)

func main() {
	// initialize auth
	auth, err := core.NewAuth(core.AuthConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("unable to initialize auth due: %v", err)
	}
	// load questions data
	data, err := ioutil.ReadFile("./questions.json")
	if err != nil {
		log.Fatalf("unable to open questions data file due: %v", err)
	}
	var questions []core.Question
	err = json.Unmarshal(data, &questions)
	if err != nil {
		log.Fatalf("unable to parse questions data file due: %v", err)
	}
	// initialize question storage
	questionStorage, err := queststrg.New(queststrg.Config{
		Questions: questions,
	})
	if err != nil {
		log.Fatalf("unable to initialize question storage due: %v", err)
	}
	// initialize timeout calculator
	toutCalc, err := toutcalc.New(toutcalc.StandardConfig())
	if err != nil {
		log.Fatalf("unable to initialize timeout calculator due: %v", err)
	}
	// initialize service
	service, err := core.NewService(core.ServiceConfig{
		GameStorage:       gamestrg.New(),
		QuestionStorage:   questionStorage,
		TimeoutCalculator: toutCalc,
		Clock:             clock.New(),
		AddScore:          1,
	})
	if err != nil {
		log.Fatalf("unable to initialize service due: %v", err)
	}
	// initialize api
	api, err := rest.NewAPI(rest.APIConfig{
		Auth:    auth,
		Service: service,
	})
	if err != nil {
		log.Fatalf("unable to initialize api due: %v", err)
	}
	// initialize server
	server := &http.Server{
		Addr:        fmt.Sprintf(":%v", listenPort),
		Handler:     api.GetHandler(),
		ReadTimeout: 3 * time.Second,
	}
	// initialize shutdown signal receiver
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-signalCh
		log.Printf("shutdown signal received: %+v, terminating app...", sig)
		server.Shutdown(context.Background())
	}()
	// run server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("unable to run server due: %v", err)
	}
}
