package api

import "net/http"

func NewMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/api/v1/jobs", CreateJobHandler)

	return withCORS(mux)
}
