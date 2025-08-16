package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/wahlly/currency_converter/controllers"
	"github.com/wahlly/currency_converter/routes"
	"github.com/wahlly/currency_converter/services"
)

func main() {
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

	port := ":2030"
	log.Printf("starting server on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}