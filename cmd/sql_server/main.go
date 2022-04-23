package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ghazlabs/hex-mathrush/internal/core"
	"github.com/ghazlabs/hex-mathrush/internal/driven/clock"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/mysql/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/toutcalc"
	"github.com/ghazlabs/hex-mathrush/internal/driver/rest"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

const (
	listenPort   = 9191
	apiKey       = "c4211664-47dc-4887-a2fe-9e694fbaf55a"
	envKeySQLDSN = "SQL_DSN"
)

func main() {
	// initialize auth
	auth, err := core.NewAuth(core.AuthConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("unable to initialize auth due: %v", err)
	}
	// initialize sql client
	sqlClient, err := sqlx.Connect("mysql", os.Getenv(envKeySQLDSN))
	if err != nil {
		log.Fatalf("unable to initialize sql client due: %v", err)
	}
	// initialize question storage
	questionStorage, err := queststrg.New(queststrg.Config{
		SQLClient: sqlClient,
	})
	if err != nil {
		log.Fatalf("unable to initialize question storage due: %v", err)
	}
	// initialize game storage
	gameStorage, err := gamestrg.New(gamestrg.Config{
		SQLClient: sqlClient,
	})
	if err != nil {
		log.Fatalf("unable to initialize game storage due: %v", err)
	}
	// initialize timeout calculator
	toutCalc, err := toutcalc.New(toutcalc.StandardConfig())
	if err != nil {
		log.Fatalf("unable to initialize timeout calculator due: %v", err)
	}
	// initialize service
	service, err := core.NewService(core.ServiceConfig{
		GameStorage:       gameStorage,
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
	log.Printf("server is listening on :%v", listenPort)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("unable to run server due: %v", err)
	}
}
