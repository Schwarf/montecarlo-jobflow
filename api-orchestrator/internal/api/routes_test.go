package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

func TestNewMuxRoutesHealthRequest(t *testing.T) {
	h := NewHandler(nil)
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", resp["status"])
	}

	if resp["service"] != "api-orchestrator" {
		t.Fatalf("expected service %q, got %q", "api-orchestrator", resp["service"])
	}
}

func TestNewMuxRoutesCreateJobRequest(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)
	mux := NewMux(h)

	body := `{
		"name": "test-job",
		"integrand": "x + 1",
		"variables": [
			{"name": "x", "lower": "0", "upper": "1"}
		],
		"evaluations": 1000
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}

	var resp CreateJobResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.JobID == "" {
		t.Fatal("expected non-empty jobId")
	}

	if resp.Status != string(job.StatusQueued) {
		t.Fatalf("expected status %q, got %q", job.StatusQueued, resp.Status)
	}

	if repo.createdJob == nil {
		t.Fatal("expected job to be persisted")
	}

	if repo.createdJob.ID != resp.JobID {
		t.Fatalf("persisted job ID mismatch: got %q want %q", repo.createdJob.ID, resp.JobID)
	}
}

func TestNewMuxRoutesGetJobRequest(t *testing.T) {
	const jobID = "job-123"

	createdAt := time.Date(2026, 4, 19, 12, 30, 45, 0, time.UTC)
	updatedAt := time.Date(2026, 4, 19, 12, 45, 0, 0, time.UTC)

	repo := &fakeRepository{
		jobByID: job.Job{
			ID:        jobID,
			Name:      "test-job",
			Integrand: "x + 1",
			IntegrationVariables: []job.VariableSpec{
				{Name: "x", Lower: "0", Upper: "1"},
			},
			Evaluations: 1000,
			Status:      job.StatusQueued,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		},
	}

	h := NewHandler(repo)
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+jobID, nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp GetJobResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.JobID != jobID {
		t.Fatalf("expected jobId %q, got %q", jobID, resp.JobID)
	}

	if resp.Name != "test-job" {
		t.Fatalf("expected name %q, got %q", "test-job", resp.Name)
	}

	if resp.Status != string(job.StatusQueued) {
		t.Fatalf("expected status %q, got %q", job.StatusQueued, resp.Status)
	}

	if len(resp.IntegrationVariables) != 1 {
		t.Fatalf("expected 1 integration variable, got %d", len(resp.IntegrationVariables))
	}

	gotVar := resp.IntegrationVariables[0]
	if gotVar.Name != "x" || gotVar.Lower != "0" || gotVar.Upper != "1" {
		t.Fatalf("unexpected integration variable: %+v", gotVar)
	}
}

func TestNewMuxRejectsWrongMethodForHealthRoute(t *testing.T) {
	h := NewHandler(nil)
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestNewMuxRejectsWrongMethodForCreateJobRoute(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	if repo.createdJob != nil {
		t.Fatal("expected no job to be persisted")
	}
}

func TestNewMuxRejectsWrongMethodForGetJobRoute(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs/job-123", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
}

func TestNewMuxReturnsNotFoundForUnknownRoute(t *testing.T) {
	h := NewHandler(nil)
	mux := NewMux(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/unknown", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
