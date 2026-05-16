package api

import (
	"strings"
	"testing"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

func validCreateJobRequest() CreateJobRequest {
	return CreateJobRequest{
		Name:      "test-job",
		Integrand: "x + 1",
		IntegrationVariables: []job.VariableSpec{
			{
				Name:  "x",
				Lower: "0",
				Upper: "1",
			},
		},
		Evaluations: 1000,
	}
}

func TestValidateBasicAcceptsValidRequest(t *testing.T) {
	req := validCreateJobRequest()

	if err := req.ValidateBasic(); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestValidateBasicRejectsEmptyName(t *testing.T) {
	req := validCreateJobRequest()
	req.Name = ""

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestValidateBasicRejectsInvalidComputationJob(t *testing.T) {
	req := validCreateJobRequest()
	req.Integrand = ""

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for invalid computation job, got nil")
	}
	if !strings.Contains(err.Error(), "integrand must not be empty") {
		t.Fatalf("expected computation validation error, got %v", err)
	}
}

func TestValidateBasicAcceptsMultipleDistinctVariables(t *testing.T) {
	req := CreateJobRequest{
		Name:      "test-job",
		Integrand: "x + y",
		IntegrationVariables: []job.VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
			{Name: "y", Lower: "-1", Upper: "2"},
		},
		Evaluations: 1000,
	}

	if err := req.ValidateBasic(); err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}
}

func TestValidationComponentsValidIntegrand(t *testing.T) {
	const validIntegrand = "(1+x^2+y^2 + Pi*ln(1+z^2+2*x*y))^4"
	r := CreateJobRequest{
		Name:      "test-job",
		Integrand: validIntegrand,
		IntegrationVariables: []job.VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
			{Name: "y", Lower: "0", Upper: "1"},
			{Name: "z", Lower: "0", Upper: "1"},
		},
		Evaluations: 1000000,
	}
	err := r.ValidateBasic()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	err = r.ValidateSemantics()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
