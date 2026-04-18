package sqlite

import (
	"context"
	"database/sql"

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
	panic("not implemented")
}

func (r *Repository) GetByID(ctx context.Context, id string) (job.Job, error) {
	panic("not implemented")
}

var _ job.Repository = (*Repository)(nil)
