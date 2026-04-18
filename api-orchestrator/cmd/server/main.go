package main

import (
	"log"
	"net/http"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/api"
	storesqlite "github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/store/sqlite"
)

func main() {
	db, err := storesqlite.Open("jobs.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	addr := ":8080"
	mux := api.NewMux()

	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
