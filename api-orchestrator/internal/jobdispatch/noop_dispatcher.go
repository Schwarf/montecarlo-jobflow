package jobdispatch

import "context"

// NoopDispatcher satisfies Dispatcher without sending the submission anywhere.
// It is useful while the API is wired for dispatch but no worker transport exists yet.
type NoopDispatcher struct{}

func (NoopDispatcher) Dispatch(ctx context.Context, submission WorkerSubmission) error {
	return nil
}
