package expr

import "testing"

func TestValidateNumberExpression(t *testing.T) {
	expr := &NumberExpression{Value: "3.14"}

	errors := Validate(expr)
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateVariableExpression(t *testing.T) {
	expr := &VariableExpression{Name: "x"}

	errors := Validate(expr)
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateUnaryExpression(t *testing.T) {
	expr := &UnaryExpression{
		Operator: TokenMinus,
		Right:    &VariableExpression{Name: "x"},
	}

	errors := Validate(expr)
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateBinaryExpression(t *testing.T) {
	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "x"},
		Operator: TokenPlus,
		Right:    &NumberExpression{Value: "2"},
	}

	errors := Validate(expr)
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateKnownFunction(t *testing.T) {
	expr := &FunctionCallExpression{
		Name: "sin",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
		},
	}

	errors := Validate(expr)
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateUnknownFunction(t *testing.T) {
	expr := &FunctionCallExpression{
		Name: "foo",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
		},
	}

	errors := Validate(expr)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}

func TestValidateNestedExpression(t *testing.T) {
	expr := &BinaryExpression{
		Left: &FunctionCallExpression{
			Name: "sin",
			Arguments: []Expression{
				&VariableExpression{Name: "x"},
			},
		},
		Operator: TokenPlus,
		Right: &UnaryExpression{
			Operator: TokenMinus,
			Right: &BinaryExpression{
				Left:     &VariableExpression{Name: "y"},
				Operator: TokenMultiply,
				Right:    &NumberExpression{Value: "2"},
			},
		},
	}

	errors := Validate(expr)
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

type unknownExpression struct{}

func (u *unknownExpression) expressionNode() {}

func TestValidateUnknownExpression(t *testing.T) {
	expr := &unknownExpression{}

	errors := Validate(expr)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}

func TestValidateFunctionCallWithoutArguments(t *testing.T) {
	expr := &FunctionCallExpression{
		Name:      "sin",
		Arguments: nil,
	}

	errors := Validate(expr)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}
