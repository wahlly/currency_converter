package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	port := ":2030"
	log.Printf("starting server on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}