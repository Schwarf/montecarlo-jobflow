package api

import (
	"net/http"
)

func NewMux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HealthHandler)
	mux.HandleFunc("/api/v1/jobs", h.CreateJobHandler)

	return withCORS(mux)
}
