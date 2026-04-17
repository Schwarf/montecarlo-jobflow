package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "api-orchestrator",
	})
}

func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "invalid JSON body",
		})
		return
	}

	// Reject trailing garbage or multiple JSON objects.
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "request body must contain exactly one JSON object",
		})
		return
	}

	if err := req.ValidateBasic(); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := req.ValidateSemantics(); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	jobID := uuid.NewString()

	resp := CreateJobResponse{
		JobID:  jobID,
		Status: "queued",
	}

	writeJSON(w, http.StatusAccepted, resp)
}
