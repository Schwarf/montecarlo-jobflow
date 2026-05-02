package api

import (
	"net/http"
)

func NewMux(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.HealthHandler)
	mux.HandleFunc("POST /api/v1/jobs", h.CreateJobHandler)
	mux.HandleFunc("GET /api/v1/jobs/{jobId}", h.GetJobHandler)

	return withCORS(mux)
}
