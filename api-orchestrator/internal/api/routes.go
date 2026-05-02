package api

import (
	"net/http"
)

func NewMux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", h.HealthHandler)
	mux.HandleFunc("/api/v1/jobs", h.CreateJobHandler)
	mux.HandleFunc("GET /api/v1/jobs/{jobId}", h.GetJobHandler)

	return withCORS(mux)
}
