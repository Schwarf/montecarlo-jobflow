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
