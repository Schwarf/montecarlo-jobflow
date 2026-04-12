package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateJobRequest struct {
	Name        string `json:"name"`
	Integrand   string `json:"integrand"`
	Dimensions  int    `json:"dimensions"`
	Evaluations int    `json:"evaluations"`
}

type CreateJobResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
