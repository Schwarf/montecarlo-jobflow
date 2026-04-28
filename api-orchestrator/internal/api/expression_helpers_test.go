package api

import "testing"

func TestCreateJobRequestVariableNamesPreservesOrder(t *testing.T) {
	req := CreateJobRequest{
		IntegrationVariables: []VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
			{Name: "y", Lower: "0", Upper: "1"},
			{Name: "z", Lower: "0", Upper: "1"},
		},
	}

	names := req.VariableNames()

	expected := []string{"x", "y", "z"}
	if len(names) != len(expected) {
		t.Fatalf("expected %d names, got %d", len(expected), len(names))
	}

	for i := range expected {
		if names[i] != expected[i] {
			t.Fatalf("expected name %q at index %d, got %q", expected[i], i, names[i])
		}
	}
}

func TestCreateJobRequestExpressionValidationContextContainsVariables(t *testing.T) {
	req := CreateJobRequest{
		IntegrationVariables: []VariableSpec{
			{Name: "x", Lower: "0", Upper: "1"},
			{Name: "y", Lower: "0", Upper: "1"},
		},
	}

	context := req.ExpressionValidationContext()

	if _, ok := context.UserVariables["x"]; !ok {
		t.Fatal("expected variable x in validation context")
	}

	if _, ok := context.UserVariables["y"]; !ok {
		t.Fatal("expected variable y in validation context")
	}

	if _, ok := context.UserVariables["z"]; ok {
		t.Fatal("did not expect variable z in validation context")
	}
}
