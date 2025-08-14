package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/wahlly/currency_converter/controllers"
	"github.com/wahlly/currency_converter/routes"
	"github.com/wahlly/currency_converter/services"
)

func main() {
	cache := services.NewRatesCache(1*time.Minute)
	rc := &controllers.RatesController{Cache: cache}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, rc)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	port := ":2030"
	log.Printf("starting server on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}