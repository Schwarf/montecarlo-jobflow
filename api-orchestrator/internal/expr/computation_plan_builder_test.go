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

func TestComputationPlanBuilderBuildCube(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "3"},
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

	expectedAssignment1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "x"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment x*x, got %#v", b.Assignments[0].Expr)
	}

	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment h1*x, got %#v", b.Assignments[1].Expr)
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

func TestComputationPlanBuilderBuildCubeOfCompoundBase(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &BinaryExpression{
			Left: &BinaryExpression{
				Left:     &NumberExpression{Value: "1"},
				Operator: TokenPlus,
				Right:    &VariableExpression{Name: "x"},
			},
			Operator: TokenPlus,
			Right:    &VariableExpression{Name: "y"},
		},
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "3"},
	}

	expected := &VariableExpression{Name: "h3"}

	result := b.Build(expr)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h3, got %#v", result)
	}

	if len(b.Assignments) != 3 {
		t.Fatalf("expected 3 assignments, got %d", len(b.Assignments))
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

	expectedAssignment1 := &BinaryExpression{
		Left: &BinaryExpression{
			Left:     &NumberExpression{Value: "1"},
			Operator: TokenPlus,
			Right:    &VariableExpression{Name: "x"},
		},
		Operator: TokenPlus,
		Right:    &VariableExpression{Name: "y"},
	}

	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h1"},
	}

	expectedAssignment3 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h2"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h1"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment ((1+x)+y), got %#v", b.Assignments[0].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment h1*h1, got %#v", b.Assignments[1].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[2].Expr, expectedAssignment3) {
		t.Fatalf("expected third assignment h2*h1, got %#v", b.Assignments[2].Expr)
	}
}

func TestComputationPlanBuilderBuildNegativeFifthPower(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "y"},
		Operator: TokenPower,
		Right: &UnaryExpression{
			Operator: TokenMinus,
			Right:    &NumberExpression{Value: "5"},
		},
	}

	result := b.Build(expr)

	expected := &VariableExpression{Name: "h4"}

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("expected expression h4, got %#v", result)
	}

	if len(b.Assignments) != 4 {
		t.Fatalf("expected 4 assignments, got %d", len(b.Assignments))
	}

	expectedAssignment1 := &BinaryExpression{
		Left:     &VariableExpression{Name: "y"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "y"},
	}
	expectedAssignment2 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h1"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h1"},
	}
	expectedAssignment3 := &BinaryExpression{
		Left:     &VariableExpression{Name: "h2"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "y"},
	}
	expectedAssignment4 := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    &VariableExpression{Name: "h3"},
	}

	if !reflect.DeepEqual(b.Assignments[0].Expr, expectedAssignment1) {
		t.Fatalf("expected first assignment y*y, got %#v", b.Assignments[0].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[1].Expr, expectedAssignment2) {
		t.Fatalf("expected second assignment h1*h1, got %#v", b.Assignments[1].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[2].Expr, expectedAssignment3) {
		t.Fatalf("expected third assignment h2*y, got %#v", b.Assignments[2].Expr)
	}
	if !reflect.DeepEqual(b.Assignments[3].Expr, expectedAssignment4) {
		t.Fatalf("expected fourth assignment 1/h3, got %#v", b.Assignments[3].Expr)
	}
}

func TestComputationPlanBuilderBuildComplexMixedExpression(t *testing.T) {
	var b ComputationPlanBuilder

	expr := &BinaryExpression{
		Left: &BinaryExpression{
			Left: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenPower,
				Right:    &NumberExpression{Value: "3"},
			},
			Operator: TokenPlus,
			Right: &BinaryExpression{
				Left:     &VariableExpression{Name: "y"},
				Operator: TokenPower,
				Right: &UnaryExpression{
					Operator: TokenMinus,
					Right:    &NumberExpression{Value: "5"},
				},
			},
		},
		Operator: TokenDivide,
		Right: &BinaryExpression{
			Left: &FunctionCallExpression{
				Name: "cos",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			Operator: TokenPower,
			Right:    &NumberExpression{Value: "3"},
		},
	}

	result := b.Build(expr)

	top, ok := result.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected top-level *BinaryExpression, got %T", result)
	}

	if top.Operator != TokenDivide {
		t.Fatalf("expected top-level operator divide, got %v", top.Operator)
	}

	denominator, ok := top.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected denominator *VariableExpression, got %T", top.Right)
	}
	if denominator.Name != "h9" {
		t.Fatalf("expected denominator h9, got %q", denominator.Name)
	}

	numerator, ok := top.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected numerator *BinaryExpression, got %T", top.Left)
	}
	if numerator.Operator != TokenPlus {
		t.Fatalf("expected numerator operator plus, got %v", numerator.Operator)
	}

	numeratorLeft, ok := numerator.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected numerator left *VariableExpression, got %T", numerator.Left)
	}
	if numeratorLeft.Name != "h2" {
		t.Fatalf("expected numerator left h2, got %q", numeratorLeft.Name)
	}

	numeratorRight, ok := numerator.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected numerator right *VariableExpression, got %T", numerator.Right)
	}
	if numeratorRight.Name != "h6" {
		t.Fatalf("expected numerator right h6, got %q", numeratorRight.Name)
	}

	if len(b.Assignments) != 9 {
		t.Fatalf("expected 9 assignments, got %d", len(b.Assignments))
	}

	expectedAssignments := []Assignment{
		{
			Name: "h1",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "x"},
			},
		},
		{
			Name: "h2",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "h1"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "x"},
			},
		},
		{
			Name: "h3",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "y"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "y"},
			},
		},
		{
			Name: "h4",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "h3"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "h3"},
			},
		},
		{
			Name: "h5",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "h4"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "y"},
			},
		},
		{
			Name: "h6",
			Expr: &BinaryExpression{
				Left:     &NumberExpression{Value: "1"},
				Operator: TokenDivide,
				Right:    &VariableExpression{Name: "h5"},
			},
		},
		{
			Name: "h7",
			Expr: &FunctionCallExpression{
				Name: "cos",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
		},
		{
			Name: "h8",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "h7"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "h7"},
			},
		},
		{
			Name: "h9",
			Expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "h8"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "h7"},
			},
		},
	}
	for i, expectedAssignment := range expectedAssignments {
		if b.Assignments[i].Name != expectedAssignment.Name {
			t.Fatalf("expected assignment %d name %q, got %q", i, expectedAssignment.Name, b.Assignments[i].Name)
		}
		if !reflect.DeepEqual(b.Assignments[i].Expr, expectedAssignment.Expr) {
			t.Fatalf("expected assignment %d expr %#v, got %#v", i, expectedAssignment.Expr, b.Assignments[i].Expr)
		}
	}
}
