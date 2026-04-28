package job

import (
	"context"
	"errors"
)

var ErrJobNotFound = errors.New("job not found")

type Repository interface {
	Create(ctx context.Context, job Job) error
	GetByID(ctx context.Context, id string) (Job, error)

	MarkRunning(ctx context.Context, id string) error
	MarkCompleted(ctx context.Context, id string, resultJSON string) error
	MarkFailed(ctx context.Context, id string, errorMessage string) error
}
