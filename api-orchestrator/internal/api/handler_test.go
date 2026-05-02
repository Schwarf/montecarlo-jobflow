package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

type fakeRepository struct {
	createdJob *job.Job
	createErr  error
}

func (r *fakeRepository) Create(ctx context.Context, j job.Job) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.createdJob = &j
	return nil
}

func (r *fakeRepository) GetByID(ctx context.Context, id string) (job.Job, error) {
	return job.Job{}, job.ErrJobNotFound
}

func (r *fakeRepository) MarkRunning(ctx context.Context, id string) error {
	return nil
}

func (r *fakeRepository) MarkCompleted(ctx context.Context, id string, resultJSON string) error {
	return nil
}

func (r *fakeRepository) MarkFailed(ctx context.Context, id string, errorMessage string) error {
	return nil
}

func TestHealthHandlerReturnsOK(t *testing.T) {
	h := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil) // fake HTTP request
	rec := httptest.NewRecorder()                              // fake HTTP response

	h.HealthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status %q, got %q", "ok", body["status"])
	}

	if body["service"] != "api-orchestrator" {
		t.Fatalf("expected service %q, got %q", "api-orchestrator", body["service"])
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type %q, got %q", "application/json", contentType)
	}
}

func TestCreateJobHandlerAcceptsValidJob(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

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

	h.CreateJobHandler(rec, req)

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

	if repo.createdJob.Name != "test-job" {
		t.Fatalf("expected job name %q, got %q", "test-job", repo.createdJob.Name)
	}

	if repo.createdJob.Integrand != "x + 1" {
		t.Fatalf("expected integrand %q, got %q", "x + 1", repo.createdJob.Integrand)
	}

	if repo.createdJob.Evaluations != 1000 {
		t.Fatalf("expected evaluations %d, got %d", 1000, repo.createdJob.Evaluations)
	}

	if repo.createdJob.Status != job.StatusQueued {
		t.Fatalf("expected persisted status %q, got %q", job.StatusQueued, repo.createdJob.Status)
	}

	if len(repo.createdJob.IntegrationVariables) != 1 {
		t.Fatalf("expected 1 integration variable, got %d", len(repo.createdJob.IntegrationVariables))
	}

	gotVar := repo.createdJob.IntegrationVariables[0]
	if gotVar.Name != "x" || gotVar.Lower != "0" || gotVar.Upper != "1" {
		t.Fatalf("unexpected integration variable: %+v", gotVar)
	}
}
