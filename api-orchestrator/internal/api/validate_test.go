package api

import "testing"

func validCreateJobRequest() CreateJobRequest {
	return CreateJobRequest{
		Name:      "test-job",
		Integrand: "x + 1",
		IntegrationVariables: []VariableSpec{
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

func TestValidateBasicAcceptsValidVariableNames(t *testing.T) {
	validNames := []string{
		"x",
		"x1",
		"x_1",
		"Alpha",
		"A2_b",
	}

	for _, name := range validNames {
		req := validCreateJobRequest()
		req.IntegrationVariables[0].Name = name

		if err := req.ValidateBasic(); err != nil {
			t.Fatalf("expected variable name %q to be valid, got error: %v", name, err)
		}
	}
}

func TestValidateBasicRejectsInvalidVariableNames(t *testing.T) {
	invalidNames := []string{
		"",
		"_x",
		"1x",
		"x-y",
		"x y",
	}

	for _, name := range invalidNames {
		req := validCreateJobRequest()
		req.IntegrationVariables[0].Name = name

		err := req.ValidateBasic()
		if err == nil {
			t.Fatalf("expected variable name %q to be invalid, got nil error", name)
		}
	}
}

func TestValidateBasicRejectsDuplicateVariableNames(t *testing.T) {
	req := CreateJobRequest{
		Name:      "test-job",
		Integrand: "x + 1",
		IntegrationVariables: []VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
			{Name: "x", Lower: "2", Upper: "3"},
		},
		Evaluations: 1000,
	}

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for duplicate variable names, got nil")
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

func TestValidateBasicRejectsEmptyIntegrand(t *testing.T) {
	req := validCreateJobRequest()
	req.Integrand = ""

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for empty integrand, got nil")
	}
}

func TestValidateBasicRejectsMissingIntegrationVariables(t *testing.T) {
	req := validCreateJobRequest()
	req.IntegrationVariables = nil

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for missing integration variables, got nil")
	}
}

func TestValidateBasicRejectsZeroEvaluations(t *testing.T) {
	req := validCreateJobRequest()
	req.Evaluations = 0

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for zero evaluations, got nil")
	}
}

func TestValidateBasicRejectsNegativeEvaluations(t *testing.T) {
	req := validCreateJobRequest()
	req.Evaluations = -1

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for negative evaluations, got nil")
	}
}

func TestValidateBasicRejectsEmptyVariableName(t *testing.T) {
	req := validCreateJobRequest()
	req.IntegrationVariables[0].Name = ""

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for empty variable name, got nil")
	}
}

func TestValidateBasicRejectsEmptyLowerBound(t *testing.T) {
	req := validCreateJobRequest()
	req.IntegrationVariables[0].Lower = ""

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for empty lower bound, got nil")
	}
}

func TestValidateBasicRejectsEmptyUpperBound(t *testing.T) {
	req := validCreateJobRequest()
	req.IntegrationVariables[0].Upper = ""

	err := req.ValidateBasic()
	if err == nil {
		t.Fatal("expected error for empty upper bound, got nil")
	}
}

func TestValidateBasicAcceptsMultipleDistinctVariables(t *testing.T) {
	req := CreateJobRequest{
		Name:      "test-job",
		Integrand: "x + y",
		IntegrationVariables: []VariableSpec{
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
		IntegrationVariables: []VariableSpec{
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

// TODO more tests involving ValidateSemantics
