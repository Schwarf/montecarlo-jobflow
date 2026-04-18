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
    id,
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
	panic("not implemented")
}

var _ job.Repository = (*Repository)(nil)
