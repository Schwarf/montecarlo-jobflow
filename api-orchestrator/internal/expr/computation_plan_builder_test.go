package expr

import "testing"

func TestComputationPlanBuilderNewTempVariable(t *testing.T) {
	var b ComputationPlanBuilder

	if got := b.NewTempVariable(); got != "h1" {
		t.Fatalf("expected h1, got %q", got)
	}
	if got := b.NewTempVariable(); got != "h2" {
		t.Fatalf("expected h2, got %q", got)
	}
}

func TestComputationPlanBuilderEmit(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	result := b.Emit(expr)

	if result.Name != "h1" {
		t.Fatalf("expected returned variable name h1, got %q", result.Name)
	}

	if len(b.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected assignment name h1, got %q", b.Assignments[0].Name)
	}

	if b.Assignments[0].Expr != expr {
		t.Fatal("expected emitted expression to be stored unchanged")
	}
}

func TestComputationPlanBuilderEmitMultipleAssignments(t *testing.T) {
	var b ComputationPlanBuilder

	expr1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	expr2 := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    &VariableExpression{Name: "h1"},
	}

	result1 := b.Emit(expr1)
	result2 := b.Emit(expr2)

	if result1.Name != "h1" {
		t.Fatalf("expected first returned variable name h1, got %q", result1.Name)
	}
	if result2.Name != "h2" {
		t.Fatalf("expected second returned variable name h2, got %q", result2.Name)
	}

	if len(b.Assignments) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected first assignment name h1, got %q", b.Assignments[0].Name)
	}
	if b.Assignments[1].Name != "h2" {
		t.Fatalf("expected second assignment name h2, got %q", b.Assignments[1].Name)
	}

	if b.Assignments[0].Expr != expr1 {
		t.Fatal("expected first emitted expression to be stored unchanged")
	}
	if b.Assignments[1].Expr != expr2 {
		t.Fatal("expected second emitted expression to be stored unchanged")
	}
}
