package expr

import (
	"reflect"
	"testing"
)

func TestComputationPlanBuilderNewTempVariable(t *testing.T) {
	var b ComputationPlanBuilder

	if got := b.NewTempVariable(); got != "h1" {
		t.Fatalf("expected h1, got %q", got)
	}
	if got := b.NewTempVariable(); got != "h2" {
		t.Fatalf("expected h2, got %q", got)
	}
}

func TestComputationPlanBuilderAssignToTempVariable(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	result := b.AssignToTempVariable(expr)

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

func TestComputationPlanBuilderMultipleAssignments(t *testing.T) {
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

	result1 := b.AssignToTempVariable(expr1)
	result2 := b.AssignToTempVariable(expr2)

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

func TestComputationPlanBuilderBuildSquare(t *testing.T) {
	var b ComputationPlanBuilder

	expr1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "2"},
	}

	expected := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	result, ok := b.BuildSquare(expr1)
	if !ok {
		t.Fatal("expected a square")
	}

	variable, ok := result.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", result)
	}

	if variable.Name != "h1" {
		t.Fatalf("expected returned variable name h1, got %q", variable.Name)
	}

	if len(b.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected first assignment name h1, got %q", b.Assignments[0].Name)
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expected) {
		t.Fatalf("expected multiplication expression to be stored, got %#v", b.Assignments[0].Expr)
	}
}

func TestComputationPlanBuilderBuild(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &BinaryExpression{
			Left:     &VariableExpression{Name: "x"},
			Operator: TokenPower,
			Right:    &NumberExpression{Value: "2"},
		},
		Operator: TokenPlus,
		Right:    &VariableExpression{Name: "y"},
	}

	expected := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenPlus,
		Right:    &VariableExpression{Name: "y"},
	}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression 'h1 + y'")
	}
	if len(b.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected assignment name h1, got %q", b.Assignments[0].Name)
	}

}

func TestComputationPlanBuilderBuildWithTwoPowers(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &BinaryExpression{
			Left:     &VariableExpression{Name: "x"},
			Operator: TokenPower,
			Right:    &NumberExpression{Value: "2"},
		},
		Operator: TokenPlus,
		Right: &BinaryExpression{
			Left:     &VariableExpression{Name: "z"},
			Operator: TokenPower,
			Right:    &NumberExpression{Value: "2"},
		},
	}

	expected := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenPlus,
		Right:    &VariableExpression{Name: "h2"},
	}

	result := b.Build(expr)

	expectedAssignment1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "z"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "z"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment x*x, got %#v", b.Assignments[0].Expr)
	}

	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment z*z, got %#v", b.Assignments[1].Expr)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression 'h1 + h2'")
	}
	if len(b.Assignments) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected assignment name h1, got %q", b.Assignments[0].Name)
	}

	if b.Assignments[1].Name != "h2" {
		t.Fatalf("expected assignment name h2, got %q", b.Assignments[1].Name)
	}
}

func TestComputationPlanBuilderAssignNonTrivialToTempVariable(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &FunctionCallExpression{Name: "sin", Arguments: []Expression{&VariableExpression{Name: "x"}}}

	result := b.AssignNonTrivialToTempVariable(expr)
	variable, ok := result.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", result)
	}

	if variable.Name != "h1" {
		t.Fatalf("expected returned variable name h1, got %q", variable.Name)
	}

	if len(b.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected assignment name h1, got %q", b.Assignments[0].Name)
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expr) {
		t.Fatalf("expected stored expression sin(x), got %#v", b.Assignments[0].Expr)
	}
}

func TestComputationPlanBuilderAssignNonTrivialToTempVariableTrivialCase(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &VariableExpression{Name: "x"}

	result := b.AssignNonTrivialToTempVariable(expr)

	variable, ok := result.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", result)
	}

	if variable.Name != "x" {
		t.Fatalf("expected variable name x, got %q", variable.Name)
	}

	if len(b.Assignments) != 0 {
		t.Fatalf("expected 0 assignments, got %d", len(b.Assignments))
	}
}

func TestComputationPlanBuilderBuildSinSquared(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &FunctionCallExpression{
			Name: "sin",
			Arguments: []Expression{
				&VariableExpression{Name: "x"},
			},
		},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "2"},
	}

	expected := &VariableExpression{Name: "h2"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h2, got %#v", result)
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

	expectedAssignment1 := &FunctionCallExpression{
		Name: "sin",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
		},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h1"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment sin(x), got %#v", b.Assignments[0].Expr)
	}

	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment h1*h1, got %#v", b.Assignments[1].Expr)
	}
}

func TestComputationPlanBuilderBuildLeavesCubeUnchanged(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "3"},
	}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expr) {
		t.Fatalf("expected expression to stay unchanged, got %#v", result)
	}

	if len(b.Assignments) != 0 {
		t.Fatalf("expected 0 assignments, got %d", len(b.Assignments))
	}
}

func TestComputationPlanBuilderBuildSimplifiesPowerOfOne(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "1"},
	}

	result := b.Build(expr)

	expected := &VariableExpression{Name: "x"}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression x, got %#v", result)
	}

	if len(b.Assignments) != 0 {
		t.Fatalf("expected 0 assignments, got %d", len(b.Assignments))
	}
}

func TestComputationPlanBuilderSimplifyPowerOfMinusOne(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right: &UnaryExpression{
			Operator: TokenMinus,
			Right:    &NumberExpression{Value: "1"},
		},
	}

	expected := &VariableExpression{Name: "h1"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h1, got %#v", result)
	}

	if len(b.Assignments) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(b.Assignments))
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected assignment name h1, got %q", b.Assignments[0].Name)
	}

	expectedAssignment := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    &VariableExpression{Name: "x"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment) {
		t.Fatalf("expected assignment 1/x, got %#v", b.Assignments[0].Expr)
	}
}

func TestComputationPlanBuilderSimplifyFunctionWithPowerOfMinusOne(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &FunctionCallExpression{
			Name: "sin",
			Arguments: []Expression{
				&VariableExpression{Name: "x"},
			},
		},
		Operator: TokenPower,
		Right: &UnaryExpression{
			Operator: TokenMinus,
			Right:    &NumberExpression{Value: "1"},
		},
	}

	expected := &VariableExpression{Name: "h2"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h2, got %#v", result)
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

	expectedAssignment1 := &FunctionCallExpression{
		Name: "sin",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
		},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    &VariableExpression{Name: "h1"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment sin(x), got %#v", b.Assignments[0].Expr)
	}

	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment 1/h1, got %#v", b.Assignments[1].Expr)
	}
}

func TestComputationPlanBuilderBuildInverseSquare(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right: &UnaryExpression{
			Operator: TokenMinus,
			Right:    &NumberExpression{Value: "2"},
		},
	}

	expected := &VariableExpression{Name: "h2"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h2, got %#v", result)
	}

	if len(b.Assignments) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(b.Assignments))
	}

	expectedAssignment1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    &VariableExpression{Name: "h1"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment x*x, got %#v", b.Assignments[0].Expr)
	}

	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment 1/h1, got %#v", b.Assignments[1].Expr)
	}
}

func TestComputationPlanBuilderBuildInverseSquareOfSin(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &FunctionCallExpression{
			Name: "sin",
			Arguments: []Expression{
				&VariableExpression{Name: "x"},
			},
		},
		Operator: TokenPower,
		Right: &UnaryExpression{
			Operator: TokenMinus,
			Right:    &NumberExpression{Value: "2"},
		},
	}

	expected := &VariableExpression{Name: "h3"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h3, got %#v", result)
	}

	if len(b.Assignments) != 3 {
		t.Fatalf("expected 3 assignments, got %d", len(b.Assignments))
	}

	expectedAssignment1 := &FunctionCallExpression{
		Name: "sin",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
		},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h1"},
	}

	expectedAssignment3 := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    &VariableExpression{Name: "h2"},
	}

	if b.Assignments[0].Name != "h1" {
		t.Fatalf("expected first assignment name h1, got %q", b.Assignments[0].Name)
	}
	if b.Assignments[1].Name != "h2" {
		t.Fatalf("expected second assignment name h2, got %q", b.Assignments[1].Name)
	}
	if b.Assignments[2].Name != "h3" {
		t.Fatalf("expected third assignment name h3, got %q", b.Assignments[2].Name)
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment sin(x), got %#v", b.Assignments[0].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment h1*h1, got %#v", b.Assignments[1].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[2].Expr, expectedAssignment3) {
		t.Fatalf("expected third assignment 1/h2, got %#v", b.Assignments[2].Expr)
	}
}

func TestComputationPlanBuilderBuildFourthPower(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "4"},
	}

	expected := &VariableExpression{Name: "h2"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h2, got %#v", result)
	}

	if len(b.Assignments) != 2 {
		t.Fatalf("expected 2 assignments, got %d", len(b.Assignments))
	}

	expectedAssignment1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h1"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment x*x, got %#v", b.Assignments[0].Expr)
	}

	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment h1*h1, got %#v", b.Assignments[1].Expr)
	}
}
