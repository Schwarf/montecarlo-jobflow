package job

import (
	"context"
	"errors"
)

var ErrJobNotFound = errors.New("job not found")

type Repository interface {
	Create(ctx context.Context, job Job) error
	GetByID(ctx context.Context, id string) (Job, error)
}
