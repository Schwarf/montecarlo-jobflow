package api

import (
	"encoding/json"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	defer r.Body.Close()

	var req CreateJobRequest

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "api-orchestrator",
		})
		return
	}

	if err := req.ValidateBasic(); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	resp := CreateJobResponse{
		JobID:  "job-1",
		Status: "queued",
	}
	writeJSON(w, http.StatusAccepted, resp)
}

func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	var req CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid JSON body",
		})
		return
	}

	resp := CreateJobResponse{
		JobID:  "job-1",
		Status: "queued",
	}

	writeJSON(w, http.StatusAccepted, resp)
}
