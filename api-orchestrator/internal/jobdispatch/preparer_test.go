package jobdispatch

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/expr"
	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

type fakeCompiler struct {
	functionName  string
	integrand     string
	variableNames []string
	context       expr.ValidationContext

	header string
	err    error
}

func (c *fakeCompiler) CompileToHeader(
	functionName string,
	integrand string,
	variableNames []string,
	context expr.ValidationContext,
) (string, error) {
	c.functionName = functionName
	c.integrand = integrand
	c.variableNames = append([]string(nil), variableNames...)
	c.context = context

	if c.err != nil {
		return "", c.err
	}

	return c.header, nil
}

func validJob() job.Job {
	now := time.Date(2026, 5, 2, 12, 0, 0, 0, time.UTC)

	return job.Job{
		ID:        "job-123",
		Name:      "test-job",
		Integrand: "x + y",
		IntegrationVariables: []job.VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
			{Name: "y", Lower: "2", Upper: "3"},
		},
		Evaluations: 1000,
		Status:      job.StatusQueued,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func TestPrepareCreatesWorkerSubmissionAndCompilesHeader(t *testing.T) {
	compiler := &fakeCompiler{
		header: "// generated header",
	}

	preparer := Preparer{
		Compiler: compiler,
	}

	submission, err := preparer.Prepare(validJob())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if compiler.functionName != IntegrandFunctionName {
		t.Fatalf("expected function name %q, got %q", IntegrandFunctionName, compiler.functionName)
	}

	if compiler.integrand != "x + y" {
		t.Fatalf("expected integrand %q, got %q", "x + y", compiler.integrand)
	}

	wantVariableNames := []string{"x", "y"}
	if !reflect.DeepEqual(compiler.variableNames, wantVariableNames) {
		t.Fatalf("expected variable names %v, got %v", wantVariableNames, compiler.variableNames)
	}

	if _, ok := compiler.context.UserVariables["x"]; !ok {
		t.Fatalf("expected validation context to contain variable %q", "x")
	}

	if _, ok := compiler.context.UserVariables["y"]; !ok {
		t.Fatalf("expected validation context to contain variable %q", "y")
	}

	if submission.JobID != "job-123" {
		t.Fatalf("expected job id %q, got %q", "job-123", submission.JobID)
	}

	if submission.FunctionName != IntegrandFunctionName {
		t.Fatalf("expected function name %q, got %q", IntegrandFunctionName, submission.FunctionName)
	}

	if submission.Integrand != "x + y" {
		t.Fatalf("expected integrand %q, got %q", "x + y", submission.Integrand)
	}

	if submission.Evaluations != 1000 {
		t.Fatalf("expected evaluations %d, got %d", 1000, submission.Evaluations)
	}

	if submission.Header != "// generated header" {
		t.Fatalf("expected generated header, got %q", submission.Header)
	}

	wantVariables := []job.VariableSpec{
		{Name: "x", Lower: "0", Upper: "1"},
		{Name: "y", Lower: "2", Upper: "3"},
	}

	if !reflect.DeepEqual(submission.Variables, wantVariables) {
		t.Fatalf("expected variables %v, got %v", wantVariables, submission.Variables)
	}
}

func TestPreparePreservesVariableOrder(t *testing.T) {
	compiler := &fakeCompiler{
		header: "// generated header",
	}

	preparer := Preparer{
		Compiler: compiler,
	}

	j := validJob()
	j.IntegrationVariables = []job.VariableSpec{
		{Name: "z", Lower: "0", Upper: "1"},
		{Name: "x", Lower: "0", Upper: "1"},
		{Name: "y", Lower: "0", Upper: "1"},
	}

	_, err := preparer.Prepare(j)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	want := []string{"z", "x", "y"}
	if !reflect.DeepEqual(compiler.variableNames, want) {
		t.Fatalf("expected variable names %v, got %v", want, compiler.variableNames)
	}
}

func TestPrepareReturnsCompilerError(t *testing.T) {
	compilerErr := errors.New("compiler failed")

	preparer := Preparer{
		Compiler: &fakeCompiler{
			err: compilerErr,
		},
	}

	_, err := preparer.Prepare(validJob())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, compilerErr) {
		t.Fatalf("expected error to wrap compiler error, got %v", err)
	}
}

func TestPrepareRejectsNilCompiler(t *testing.T) {
	preparer := Preparer{}

	_, err := preparer.Prepare(validJob())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "compiler must not be nil") {
		t.Fatalf("expected nil compiler error, got %v", err)
	}
}

func TestPrepareRejectsInvalidJob(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*job.Job)
		wantErr string
	}{
		{
			name: "empty job id",
			modify: func(j *job.Job) {
				j.ID = ""
			},
			wantErr: "job id must not be empty",
		},
		{
			name: "empty integrand",
			modify: func(j *job.Job) {
				j.Integrand = ""
			},
			wantErr: "integrand must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := validJob()
			tt.modify(&j)

			preparer := Preparer{
				Compiler: &fakeCompiler{
					header: "// generated header",
				},
			}

			_, err := preparer.Prepare(j)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}
