package jobsubmission

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/jobdispatch"
)

type fakeRepository struct {
	createErr     error
	markFailedErr error

	createdJob *job.Job

	markFailedJobID        string
	markFailedErrorMessage string
	markFailedCalled       bool

	callLog *[]string
}

func (r *fakeRepository) Create(ctx context.Context, j job.Job) error {
	if r.callLog != nil {
		*r.callLog = append(*r.callLog, "create")
	}
	if r.createErr != nil {
		return r.createErr
	}

	r.createdJob = &j
	return nil
}

func (r *fakeRepository) MarkFailed(ctx context.Context, id string, errorMessage string) error {
	if r.callLog != nil {
		*r.callLog = append(*r.callLog, "mark_failed")
	}

	r.markFailedCalled = true
	r.markFailedJobID = id
	r.markFailedErrorMessage = errorMessage

	if r.markFailedErr != nil {
		return r.markFailedErr
	}

	return nil
}

type fakePreparer struct {
	submission jobdispatch.WorkerSubmission
	err        error

	preparedJob *job.Job

	callLog *[]string
}

func (p *fakePreparer) Prepare(j job.Job) (jobdispatch.WorkerSubmission, error) {
	if p.callLog != nil {
		*p.callLog = append(*p.callLog, "prepare")
	}

	p.preparedJob = &j

	if p.err != nil {
		return jobdispatch.WorkerSubmission{}, p.err
	}

	return p.submission, nil
}

type fakeDispatcher struct {
	err error

	dispatchedSubmission *jobdispatch.WorkerSubmission

	callLog *[]string
}

func (d *fakeDispatcher) Dispatch(ctx context.Context, submission jobdispatch.WorkerSubmission) error {
	if d.callLog != nil {
		*d.callLog = append(*d.callLog, "dispatch")
	}

	d.dispatchedSubmission = &submission

	if d.err != nil {
		return d.err
	}

	return nil
}

func testJob() job.Job {
	ts := time.Date(2026, 5, 16, 12, 0, 0, 0, time.UTC)

	return job.Job{
		ID:        "job-123",
		Name:      "test-job",
		Integrand: "x + 1",
		IntegrationVariables: []job.VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
		},
		Evaluations: 1000,
		Status:      job.StatusQueued,
		CreatedAt:   ts,
		UpdatedAt:   ts,
	}
}

func testSubmission() jobdispatch.WorkerSubmission {
	return jobdispatch.WorkerSubmission{
		JobID:        "job-123",
		FunctionName: jobdispatch.IntegrandFunctionName,
		Integrand:    "x + 1",
		Variables: []job.VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
		},
		Evaluations: 1000,
		Header:      "double integrand(double x) { return x + 1; }",
	}
}

func TestServiceSubmitPreparesStoresAndDispatchesJob(t *testing.T) {
	var callLog []string
	j := testJob()
	submission := testSubmission()

	repo := &fakeRepository{callLog: &callLog}
	preparer := &fakePreparer{
		submission: submission,
		callLog:    &callLog,
	}
	dispatcher := &fakeDispatcher{callLog: &callLog}

	service := NewService(repo, preparer, dispatcher)

	got, err := service.Submit(context.Background(), j)
	if err != nil {
		t.Fatalf("Submit returned error: %v", err)
	}

	if !reflect.DeepEqual(got, j) {
		t.Fatalf("returned job mismatch: got %+v want %+v", got, j)
	}

	if !reflect.DeepEqual(callLog, []string{"prepare", "create", "dispatch"}) {
		t.Fatalf("call order mismatch: got %v", callLog)
	}

	if preparer.preparedJob == nil || !reflect.DeepEqual(*preparer.preparedJob, j) {
		t.Fatalf("prepared job mismatch: got %+v want %+v", preparer.preparedJob, j)
	}

	if repo.createdJob == nil || !reflect.DeepEqual(*repo.createdJob, j) {
		t.Fatalf("created job mismatch: got %+v want %+v", repo.createdJob, j)
	}

	if dispatcher.dispatchedSubmission == nil || !reflect.DeepEqual(*dispatcher.dispatchedSubmission, submission) {
		t.Fatalf("dispatched submission mismatch: got %+v want %+v", dispatcher.dispatchedSubmission, submission)
	}

	if repo.markFailedCalled {
		t.Fatal("expected MarkFailed not to be called")
	}
}

func TestServiceSubmitReturnsPrepareErrorAndDoesNotStoreJob(t *testing.T) {
	var callLog []string
	prepareErr := errors.New("invalid integrand")

	repo := &fakeRepository{callLog: &callLog}
	preparer := &fakePreparer{
		err:     prepareErr,
		callLog: &callLog,
	}
	dispatcher := &fakeDispatcher{callLog: &callLog}

	service := NewService(repo, preparer, dispatcher)

	_, err := service.Submit(context.Background(), testJob())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, prepareErr) {
		t.Fatalf("expected wrapped prepare error, got %v", err)
	}
	if !strings.Contains(err.Error(), "prepare job dispatch") {
		t.Fatalf("expected prepare context in error, got %v", err)
	}

	if !reflect.DeepEqual(callLog, []string{"prepare"}) {
		t.Fatalf("call order mismatch: got %v", callLog)
	}

	if repo.createdJob != nil {
		t.Fatalf("expected no created job, got %+v", *repo.createdJob)
	}
	if dispatcher.dispatchedSubmission != nil {
		t.Fatalf("expected no dispatched submission, got %+v", *dispatcher.dispatchedSubmission)
	}
}

func TestServiceSubmitReturnsCreateErrorAndDoesNotDispatch(t *testing.T) {
	var callLog []string
	createErr := errors.New("database unavailable")

	repo := &fakeRepository{
		createErr: createErr,
		callLog:   &callLog,
	}
	preparer := &fakePreparer{
		submission: testSubmission(),
		callLog:    &callLog,
	}
	dispatcher := &fakeDispatcher{callLog: &callLog}

	service := NewService(repo, preparer, dispatcher)

	_, err := service.Submit(context.Background(), testJob())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, createErr) {
		t.Fatalf("expected wrapped create error, got %v", err)
	}
	if !strings.Contains(err.Error(), "create job") {
		t.Fatalf("expected create context in error, got %v", err)
	}

	if !reflect.DeepEqual(callLog, []string{"prepare", "create"}) {
		t.Fatalf("call order mismatch: got %v", callLog)
	}

	if dispatcher.dispatchedSubmission != nil {
		t.Fatalf("expected no dispatched submission, got %+v", *dispatcher.dispatchedSubmission)
	}
	if repo.markFailedCalled {
		t.Fatal("expected MarkFailed not to be called")
	}
}

func TestServiceSubmitMarksJobFailedWhenDispatchFails(t *testing.T) {
	var callLog []string
	dispatchErr := errors.New("queue unavailable")
	j := testJob()

	repo := &fakeRepository{callLog: &callLog}
	preparer := &fakePreparer{
		submission: testSubmission(),
		callLog:    &callLog,
	}
	dispatcher := &fakeDispatcher{
		err:     dispatchErr,
		callLog: &callLog,
	}

	service := NewService(repo, preparer, dispatcher)

	_, err := service.Submit(context.Background(), j)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dispatchErr) {
		t.Fatalf("expected wrapped dispatch error, got %v", err)
	}
	if !strings.Contains(err.Error(), "dispatch job") {
		t.Fatalf("expected dispatch context in error, got %v", err)
	}

	if !reflect.DeepEqual(callLog, []string{"prepare", "create", "dispatch", "mark_failed"}) {
		t.Fatalf("call order mismatch: got %v", callLog)
	}

	if !repo.markFailedCalled {
		t.Fatal("expected MarkFailed to be called")
	}
	if repo.markFailedJobID != j.ID {
		t.Fatalf("MarkFailed job id mismatch: got %q want %q", repo.markFailedJobID, j.ID)
	}
	if repo.markFailedErrorMessage != dispatchErr.Error() {
		t.Fatalf("MarkFailed error message mismatch: got %q want %q", repo.markFailedErrorMessage, dispatchErr.Error())
	}
}

func TestServiceSubmitReturnsCombinedErrorWhenDispatchAndMarkFailedFail(t *testing.T) {
	var callLog []string
	dispatchErr := errors.New("queue unavailable")
	markFailedErr := errors.New("database unavailable")

	repo := &fakeRepository{
		markFailedErr: markFailedErr,
		callLog:       &callLog,
	}
	preparer := &fakePreparer{
		submission: testSubmission(),
		callLog:    &callLog,
	}
	dispatcher := &fakeDispatcher{
		err:     dispatchErr,
		callLog: &callLog,
	}

	service := NewService(repo, preparer, dispatcher)

	_, err := service.Submit(context.Background(), testJob())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, dispatchErr) {
		t.Fatalf("expected wrapped dispatch error, got %v", err)
	}
	if !strings.Contains(err.Error(), "additionally failed to mark job as failed") {
		t.Fatalf("expected combined error context, got %v", err)
	}
	if !strings.Contains(err.Error(), markFailedErr.Error()) {
		t.Fatalf("expected mark failed error in combined error, got %v", err)
	}

	if !reflect.DeepEqual(callLog, []string{"prepare", "create", "dispatch", "mark_failed"}) {
		t.Fatalf("call order mismatch: got %v", callLog)
	}
}
