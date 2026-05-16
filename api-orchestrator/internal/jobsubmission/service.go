package jobsubmission

import (
	"context"
	"fmt"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/jobdispatch"
)

type Repository interface {
	Create(ctx context.Context, j job.Job) error
	MarkFailed(ctx context.Context, id string, errorMessage string) error
}

type Preparer interface {
	Prepare(j job.Job) (jobdispatch.WorkerSubmission, error)
}

type Service struct {
	Repository Repository
	Preparer   Preparer
	Dispatcher jobdispatch.Dispatcher
}

func (s Service) Submit(ctx context.Context, j job.Job) (job.Job, error) {
	if s.Repository == nil {
		return job.Job{}, fmt.Errorf("repository must not be nil")
	}
	if s.Preparer == nil {
		return job.Job{}, fmt.Errorf("preparer must not be nil")
	}
	if s.Dispatcher == nil {
		return job.Job{}, fmt.Errorf("dispatcher must not be nil")
	}
	submission, err := s.Preparer.Prepare(j)
	if err != nil {
		return job.Job{}, fmt.Errorf("prepare job dispatch: %w", err)
	}

	if err := s.Repository.Create(ctx, j); err != nil {
		return job.Job{}, fmt.Errorf("create job: %w", err)
	}

	if err := s.Dispatcher.Dispatch(ctx, submission); err != nil {
		markErr := s.Repository.MarkFailed(ctx, j.ID, err.Error())
		if markErr != nil {
			return job.Job{}, fmt.Errorf(
				"dispatch job: %w; additionally failed to mark job as failed: %v",
				err,
				markErr,
			)
		}

		return job.Job{}, fmt.Errorf("dispatch job: %w", err)
	}

	return j, nil
}
