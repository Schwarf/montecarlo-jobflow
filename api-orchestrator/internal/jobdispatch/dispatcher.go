package jobdispatch

import "context"

type Dispatcher interface {
	Dispatch(ctx context.Context, submission WorkerSubmission) error
}
