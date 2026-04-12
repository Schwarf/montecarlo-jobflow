package main

import (
	"log"
	"net/http"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/api"
)

func main() {
	addr := ":8080"
	mux := api.NewMux()

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
