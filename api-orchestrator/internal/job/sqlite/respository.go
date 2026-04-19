package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, j job.Job) error {
	variablesJSON, err := json.Marshal(j.IntegrationVariables)
	if err != nil {
		return fmt.Errorf("marshal integration variables: %w", err)
	}

	const query = `
INSERT INTO jobs (
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
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = r.db.ExecContext(
		ctx,
		query,
		j.ID,
		j.Name,
		j.Integrand,
		string(variablesJSON),
		j.Evaluations,
		string(j.Status),
		j.ErrorMessage,
		j.ResultJSON,
		j.CreatedAt.Format(time.RFC3339),
		j.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("insert job: %w", err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (job.Job, error) {
	const query = `
SELECT
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

	var (
		j            job.Job
		variableJSON string
		status       string
		errorMessage sql.NullString
		resultJSON   sql.NullString
		createdAt    string
		updatedAt    string
	)

	err := r.db.QueryRowContext(ctx, query, id, id).Scan(
		&j.ID,
		&j.Name,
		&j.Integrand,
		&variableJSON,
		j.Evaluations,
		&status,
		&errorMessage,
		&resultJSON,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return job.Job{}, fmt.Errorf("job %q not found: %w", id, err)
		}
		return job.Job{}, fmt.Errorf("query job by id: %w", err)
	}

	if err := json.Unmarshal([]byte(variableJSON), &j.IntegrationVariables); err != nil {
		return job.Job{}, fmt.Errorf("unmarshal integration variables: %w", err)
	}

	j.Status = job.Status(status)

	if errorMessage.Valid {
		j.ErrorMessage = &errorMessage.String
	}

	if resultJSON.Valid {
		j.ResultJSON = &resultJSON.String
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt)

	if err != nil {
		return job.Job{}, fmt.Errorf("parse created at: %w", err)
	}
	j.CreatedAt = parsedCreatedAt

	parsedUpdatedAt, err := time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return job.Job{}, fmt.Errorf("parse updated at: %w", err)
	}
	j.UpdatedAt = parsedUpdatedAt

	return j, nil
}

var _ job.Repository = (*Repository)(nil)
