package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type VariableSpec struct {
	Name  string `json:"name"`
	Lower string `json:"lower"`
	Upper string `json:"upper"`
}

type CreateJobRequest struct {
	Name                 string         `json:"name"`
	Integrand            string         `json:"integrand"`
	IntegrationVariables []VariableSpec `json:"variables"`
	Evaluations          int            `json:"evaluations"`
}

type CreateJobResponse struct {
	JobID  string `json:"jobId"`
	Status string `json:"status"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
