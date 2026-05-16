package job

import (
	"strings"
	"testing"
)

func validComputationJob() Job {
	return Job{
		Integrand: "x + 1",
		IntegrationVariables: []VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
		},
		Evaluations: 1000,
	}
}

func TestValidateForComputationAcceptsValidJob(t *testing.T) {
	j := validComputationJob()

	if err := ValidateForComputation(j); err != nil {
		t.Fatalf("expected valid job, got error: %v", err)
	}
}

func TestValidateForComputationRejectsInvalidJobs(t *testing.T) {
	tests := []struct {
		name        string
		mutate      func(*Job)
		wantMessage string
	}{
		{
			name: "empty integrand",
			mutate: func(j *Job) {
				j.Integrand = ""
			},
			wantMessage: "integrand must not be empty",
		},
		{
			name: "no variables",
			mutate: func(j *Job) {
				j.IntegrationVariables = nil
			},
			wantMessage: "at least one variable is required",
		},
		{
			name: "zero evaluations",
			mutate: func(j *Job) {
				j.Evaluations = 0
			},
			wantMessage: "evaluations must be > 0",
		},
		{
			name: "negative evaluations",
			mutate: func(j *Job) {
				j.Evaluations = -1
			},
			wantMessage: "evaluations must be > 0",
		},
		{
			name: "empty variable name",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Name = ""
			},
			wantMessage: "variable name must not be empty",
		},
		{
			name: "variable starts with digit",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Name = "1x"
			},
			wantMessage: `invalid variable name: "1x"`,
		},
		{
			name: "variable starts with underscore",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Name = "_x"
			},
			wantMessage: `invalid variable name: "_x"`,
		},
		{
			name: "variable contains invalid character",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Name = "x-y"
			},
			wantMessage: `invalid variable name: "x-y"`,
		},
		{
			name: "variable contains space",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Name = "x y"
			},
			wantMessage: `invalid variable name: "x y"`,
		},
		{
			name: "duplicate variable name",
			mutate: func(j *Job) {
				j.IntegrationVariables = append(j.IntegrationVariables, VariableSpec{
					Name:  "x",
					Lower: "2",
					Upper: "3",
				})
			},
			wantMessage: `duplicate variable name: "x"`,
		},
		{
			name: "empty lower bound",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Lower = ""
			},
			wantMessage: `variable "x" has empty lower bound`,
		},
		{
			name: "empty upper bound",
			mutate: func(j *Job) {
				j.IntegrationVariables[0].Upper = ""
			},
			wantMessage: `variable "x" has empty upper bound`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := validComputationJob()
			tt.mutate(&j)

			err := ValidateForComputation(j)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if !strings.Contains(err.Error(), tt.wantMessage) {
				t.Fatalf("expected error containing %q, got %q", tt.wantMessage, err.Error())
			}
		})
	}
}

func TestValidateForComputationAcceptsValidVariableNames(t *testing.T) {
	tests := []string{
		"x",
		"x1",
		"x_1",
		"Alpha",
		"A2_b",
		"theta",
	}

	for _, variableName := range tests {
		t.Run(variableName, func(t *testing.T) {
			j := validComputationJob()
			j.IntegrationVariables[0].Name = variableName

			if err := ValidateForComputation(j); err != nil {
				t.Fatalf("expected valid variable name %q, got error: %v", variableName, err)
			}
		})
	}
}
