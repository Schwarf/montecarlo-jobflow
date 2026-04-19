package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
	storesqlite "github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/store/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := storesqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open in-memory sqlite db: %v", err)
	}

	if err := storesqlite.InitSchema(db); err != nil {
		_ = db.Close()
		t.Fatalf("init schema: %v", err)
	}

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

func testJob() job.Job {
	ts := time.Date(2026, 4, 19, 12, 30, 45, 0, time.UTC)

	return job.Job{
		ID:                   "job-123",
		Name:                 "test-job",
		Integrand:            "(1+x^2)^2",
		IntegrationVariables: nil,
		Evaluations:          1000000,
		Status:               job.StatusQueued,
		ErrorMessage:         nil,
		ResultJSON:           nil,
		CreatedAt:            ts,
		UpdatedAt:            ts,
	}
}

func TestNewRepository(t *testing.T) {
	db := newTestDB(t)

	repo := NewRepository(db)
	if repo == nil {
		t.Fatal("expected repository, got nil")
	}
	if repo.db != db {
		t.Fatal("repository does not store db pointer")
	}
}

func TestRepositoryCreateInsertsJob(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	j := testJob()

	if err := repo.Create(context.Background(), j); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	var (
		rowID         int64
		jobID         string
		name          string
		integrand     string
		variablesJSON string
		evaluations   int
		status        string
		errorMessage  sql.NullString
		resultJSON    sql.NullString
		createdAt     string
		updatedAt     string
	)

	const query = `
SELECT
    row_id,
    job_id,
    name,
    integrand,
    variables_json,
    evaluations,
    status,
    error_message,
    result_json,
    created_at,
    updated_at
FROM jobs
WHERE job_id = ?;`

	err := db.QueryRow(query, j.ID).Scan(
		&rowID,
		&jobID,
		&name,
		&integrand,
		&variablesJSON,
		&evaluations,
		&status,
		&errorMessage,
		&resultJSON,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		t.Fatalf("query inserted row: %v", err)
	}

	if rowID <= 0 {
		t.Fatalf("expected auto-generated positive row_id, got %d", rowID)
	}
	if jobID != j.ID {
		t.Fatalf("job_id mismatch: got %q want %q", jobID, j.ID)
	}
	if name != j.Name {
		t.Fatalf("name mismatch: got %q want %q", name, j.Name)
	}
	if integrand != j.Integrand {
		t.Fatalf("integrand mismatch: got %q want %q", integrand, j.Integrand)
	}
	if variablesJSON != "null" {
		t.Fatalf("variables_json mismatch: got %q want %q", variablesJSON, "null")
	}
	if evaluations != j.Evaluations {
		t.Fatalf("evaluations mismatch: got %d want %d", evaluations, j.Evaluations)
	}
	if status != string(j.Status) {
		t.Fatalf("status mismatch: got %q want %q", status, string(j.Status))
	}
	if errorMessage.Valid {
		t.Fatalf("expected NULL error_message, got %q", errorMessage.String)
	}
	if resultJSON.Valid {
		t.Fatalf("expected NULL result_json, got %q", resultJSON.String)
	}
	if createdAt != j.CreatedAt.Format(time.RFC3339) {
		t.Fatalf("created_at mismatch: got %q want %q", createdAt, j.CreatedAt.Format(time.RFC3339))
	}
	if updatedAt != j.UpdatedAt.Format(time.RFC3339) {
		t.Fatalf("updated_at mismatch: got %q want %q", updatedAt, j.UpdatedAt.Format(time.RFC3339))
	}
}

func TestRepositoryCreateDuplicateJobIDFails(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	j := testJob()

	if err := repo.Create(context.Background(), j); err != nil {
		t.Fatalf("first Create returned error: %v", err)
	}

	err := repo.Create(context.Background(), j)
	if err == nil {
		t.Fatal("expected duplicate insert to fail, got nil")
	}

	if !strings.Contains(err.Error(), "insert job") {
		t.Fatalf("expected wrapped insert error, got %v", err)
	}
}

// This test verifies the expected behavior with the current SQLite driver
// (modernc.org/sqlite): if the context is already canceled before ExecContext
// is called, Create should return a cancellation-related error. The exact error
// behavior is somewhat driver-dependent and may differ with another SQL driver.
func TestRepositoryCreateWithCanceledContextFails(t *testing.T) {
	db := newTestDB(t)
	repo := NewRepository(db)

	j := testJob()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := repo.Create(ctx, j)
	if err == nil {
		t.Fatal("expected error for canceled context, got nil")
	}

	if !errors.Is(err, context.Canceled) && !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected context canceled error, got %v", err)
	}
}
