package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
	"github.com/google/uuid"
)

type Handler struct {
	repo job.Repository
}

func NewHandler(repo job.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
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
	now := time.Now().UTC()

	vars := make([]job.VariableSpec, 0, len(req.IntegrationVariables))
	for _, v := range req.IntegrationVariables {
		vars = append(vars, job.VariableSpec{
			Name:  v.Name,
			Lower: v.Lower,
			Upper: v.Upper,
		})
	}

	j := job.Job{
		ID:                   jobID,
		Name:                 req.Name,
		Integrand:            req.Integrand,
		IntegrationVariables: vars,
		Evaluations:          req.Evaluations,
		Status:               job.StatusQueued,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	if err := h.repo.Create(r.Context(), j); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "failed to persist job",
		})
		return
	}

	resp := CreateJobResponse{
		JobID:  jobID,
		Status: string(job.StatusQueued),
	}

	writeJSON(w, http.StatusAccepted, resp)
}

func (h *Handler) GetJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{
			Error: "method not allowed",
		})
		return
	}

	jobID := r.PathValue("jobId")
	if jobID == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Error: "missing job id",
		})
		return
	}

	j, err := h.repo.GetByID(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, job.ErrJobNotFound) {
			writeJSON(w, http.StatusNotFound, ErrorResponse{
				Error: "job not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, ErrorResponse{
			Error: "failed to retrieve job",
		})
		return
	}

	integrationVars := make([]VariableSpec, 0, len(j.IntegrationVariables))
	for _, v := range j.IntegrationVariables {
		integrationVars = append(integrationVars, VariableSpec{
			Name:  v.Name,
			Lower: v.Lower,
			Upper: v.Upper,
		})
	}

	writeJSON(w, http.StatusOK, GetJobResponse{
		JobID:                j.ID,
		Name:                 j.Name,
		Integrand:            j.Integrand,
		IntegrationVariables: integrationVars,
		Evaluations:          j.Evaluations,
		Status:               string(j.Status),
		ErrorMessage:         j.ErrorMessage,
		ResultJSON:           j.ResultJSON,
		CreatedAt:            j.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            j.UpdatedAt.Format(time.RFC3339),
	})
}
