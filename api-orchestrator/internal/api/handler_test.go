package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

type fakeRepository struct {
	createdJob *job.Job
	createErr  error

	jobByID job.Job
	getErr  error
}

func (r *fakeRepository) Create(ctx context.Context, j job.Job) error {
	if r.createErr != nil {
		return r.createErr
	}

	r.createdJob = &j
	return nil
}

func (r *fakeRepository) GetByID(ctx context.Context, id string) (job.Job, error) {
	if r.getErr != nil {
		return job.Job{}, r.getErr
	}

	if r.jobByID.ID == id {
		return r.jobByID, nil
	}

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

func TestGetJobHandlerReturnsExistingJob(t *testing.T) {
	const jobID = "job-123"

	resultJSON := `{"estimate":1.234,"error":0.001}`
	errorMessage := "some previous warning"

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
			Evaluations:  1000,
			Status:       job.StatusCompleted,
			ErrorMessage: &errorMessage,
			ResultJSON:   &resultJSON,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		},
	}

	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+jobID, nil)
	req.SetPathValue("jobId", jobID)

	rec := httptest.NewRecorder()

	h.GetJobHandler(rec, req)

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
	if resp.Integrand != "x + 1" {
		t.Fatalf("expected integrand %q, got %q", "x + 1", resp.Integrand)
	}
	if resp.Evaluations != 1000 {
		t.Fatalf("expected evaluations %d, got %d", 1000, resp.Evaluations)
	}
	if resp.Status != string(job.StatusCompleted) {
		t.Fatalf("expected status %q, got %q", job.StatusCompleted, resp.Status)
	}
	if resp.ResultJSON == nil || *resp.ResultJSON != resultJSON {
		t.Fatalf("expected resultJson %q, got %v", resultJSON, resp.ResultJSON)
	}
	if resp.ErrorMessage == nil || *resp.ErrorMessage != errorMessage {
		t.Fatalf("expected errorMessage %q, got %v", errorMessage, resp.ErrorMessage)
	}
	if resp.CreatedAt != createdAt.Format(time.RFC3339) {
		t.Fatalf("expected createdAt %q, got %q", createdAt.Format(time.RFC3339), resp.CreatedAt)
	}
	if resp.UpdatedAt != updatedAt.Format(time.RFC3339) {
		t.Fatalf("expected updatedAt %q, got %q", updatedAt.Format(time.RFC3339), resp.UpdatedAt)
	}

	if len(resp.IntegrationVariables) != 1 {
		t.Fatalf("expected 1 integration variable, got %d", len(resp.IntegrationVariables))
	}

	gotVar := resp.IntegrationVariables[0]
	if gotVar.Name != "x" || gotVar.Lower != "0" || gotVar.Upper != "1" {
		t.Fatalf("unexpected integration variable: %+v", gotVar)
	}
}

func TestGetJobHandlerReturnsNotFoundForUnknownJob(t *testing.T) {
	repo := &fakeRepository{
		getErr: job.ErrJobNotFound,
	}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/unknown-job-id", nil)
	req.SetPathValue("jobId", "unknown-job-id")

	rec := httptest.NewRecorder()

	h.GetJobHandler(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "job not found" {
		t.Fatalf("expected error %q, got %q", "job not found", resp.Error)
	}
}

func TestCreateJobHandlerReturnsInternalServerErrorWhenCreateFails(t *testing.T) {
	repo := &fakeRepository{
		createErr: fmt.Errorf("database unavailable"),
	}
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

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "failed to persist job" {
		t.Fatalf("expected error %q, got %q", "failed to persist job", resp.Error)
	}
}

func TestCreateJobHandlerRejectsInvalidJSON(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

	body := `{"name": "test-job"`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateJobHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "invalid JSON body" {
		t.Fatalf("expected error %q, got %q", "invalid JSON body", resp.Error)
	}

	if repo.createdJob != nil {
		t.Fatal("expected no job to be persisted")
	}
}

func TestCreateJobHandlerRejectsUnknownJSONField(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

	body := `{
		"name": "test-job",
		"integrand": "x + 1",
		"variables": [
			{"name": "x", "lower": "0", "upper": "1"}
		],
		"evaluations": 1000,
		"unexpected": true
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateJobHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "invalid JSON body" {
		t.Fatalf("expected error %q, got %q", "invalid JSON body", resp.Error)
	}

	if repo.createdJob != nil {
		t.Fatal("expected no job to be persisted")
	}
}

func TestCreateJobHandlerRejectsMultipleJSONObjects(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

	body := `{
		"name": "test-job",
		"integrand": "x + 1",
		"variables": [
			{"name": "x", "lower": "0", "upper": "1"}
		],
		"evaluations": 1000
	}
	{
		"name": "second-job",
		"integrand": "x + 2",
		"variables": [
			{"name": "x", "lower": "0", "upper": "1"}
		],
		"evaluations": 1000
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateJobHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	expectedError := "request body must contain exactly one JSON object"
	if resp.Error != expectedError {
		t.Fatalf("expected error %q, got %q", expectedError, resp.Error)
	}

	if repo.createdJob != nil {
		t.Fatal("expected no job to be persisted")
	}
}

func TestCreateJobHandlerRejectsBasicValidationError(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

	body := `{
		"name": "",
		"integrand": "x + 1",
		"variables": [
			{"name": "x", "lower": "0", "upper": "1"}
		],
		"evaluations": 1000
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateJobHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "name must not be empty" {
		t.Fatalf("expected error %q, got %q", "name must not be empty", resp.Error)
	}

	if repo.createdJob != nil {
		t.Fatal("expected no job to be persisted")
	}
}

func TestCreateJobHandlerRejectsSemanticValidationError(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

	body := `{
		"name": "test-job",
		"integrand": "x + y",
		"variables": [
			{"name": "x", "lower": "0", "upper": "1"}
		],
		"evaluations": 1000
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/jobs", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.CreateJobHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error == "" {
		t.Fatal("expected non-empty validation error")
	}

	if repo.createdJob != nil {
		t.Fatal("expected no job to be persisted")
	}
}

func TestGetJobHandlerReturnsInternalServerErrorForRepositoryError(t *testing.T) {
	repo := &fakeRepository{
		getErr: fmt.Errorf("database unavailable"),
	}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/job-123", nil)
	req.SetPathValue("jobId", "job-123")

	rec := httptest.NewRecorder()

	h.GetJobHandler(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "failed to retrieve job" {
		t.Fatalf("expected error %q, got %q", "failed to retrieve job", resp.Error)
	}
}

func TestGetJobHandlerRejectsMissingJobID(t *testing.T) {
	repo := &fakeRepository{}
	h := NewHandler(repo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/", nil)
	rec := httptest.NewRecorder()

	h.GetJobHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response body: %v", err)
	}

	if resp.Error != "missing job id" {
		t.Fatalf("expected error %q, got %q", "missing job id", resp.Error)
	}
}
