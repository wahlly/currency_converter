package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wahlly/currency_converter/cmd/api/controllers"
	"github.com/wahlly/currency_converter/cmd/api/routes"
	"github.com/wahlly/currency_converter/cmd/api/services"
)

// const version = "1.0.0"

// type config struct {
// 	port	int
// 	env	string
// }

// type application struct {
// 	config config
// 	logger *slog.Logger
// }

var (
	port	int
	env	string
)

func main() {
	// var cfg config
	flag.IntVar(&port, "port", 2030, "server port")
	flag.StringVar(&env, "env", "development", "Environment")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	redisCient :=  redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	ctx := context.Background()
	if err := redisCient.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connection failed: %s", err)
	}

	rateService := services.NewRatesService(redisCient)
	services.HandleRateUpdates(ctx, redisCient)
	rc := &controllers.RatesController{RateService: rateService}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, rc)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	svr := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", svr.Addr, "env", env)
	err := svr.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}