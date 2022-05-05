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
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/redis/gamestrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/storage/redis/queststrg"
	"github.com/ghazlabs/hex-mathrush/internal/driven/toutcalc"
	"github.com/ghazlabs/hex-mathrush/internal/driver/rest"
	"github.com/go-redis/redis/v8"
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

	redisClient := redis.NewClient(&redis.Options{Addr: "redis:6379"})
	if err != nil {
		log.Fatalf("unable to initialize redis client due: %v", err)
	}

	questionStorage, err := queststrg.New(queststrg.Config{RedisClient: redisClient})
	if err != nil {
		log.Fatalf("unable to initialize question storage due: %v", err)
	}

	gameStorage, err := gamestrg.New(gamestrg.Config{RedisClient: redisClient})
	if err != nil {
		log.Fatalf("unable to initialize game storage due: %v", err)
	}

	toutCalc, err := toutcalc.New(toutcalc.StandardConfig())
	if err != nil {
		log.Fatalf("unable to initialize timeout calculator due: %v", err)
	}

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
