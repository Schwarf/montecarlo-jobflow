package expr

import "testing"

func testValidationContext(userVariables ...string) ValidationContext {
	context := DefaultValidationContext()
	context.UserVariables = make(map[string]struct{}, len(userVariables))
	for _, name := range userVariables {
		context.UserVariables[name] = struct{}{}
	}
	return context
}

func TestValidateNumberExpression(t *testing.T) {
	expr := &NumberExpression{Value: "3.14"}

	errors := Validate(expr, DefaultValidationContext())
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateVariableExpression(t *testing.T) {
	expr := &VariableExpression{Name: "x"}

	errors := Validate(expr, testValidationContext("x"))
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateUnaryExpression(t *testing.T) {
	expr := &UnaryExpression{
		Operator: TokenMinus,
		Right:    &VariableExpression{Name: "x"},
	}

	errors := Validate(expr, testValidationContext("x"))
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

	errors := Validate(expr, testValidationContext("x"))
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateKnownFunction(t *testing.T) {
	expr := &FunctionCallExpression{
		Name: "sin",
		Arguments: []Expression{
			&VariableExpression{Name: "y1"},
		},
	}

	errors := Validate(expr, testValidationContext("y1"))
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

	errors := Validate(expr, testValidationContext("x"))
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

	errors := Validate(expr, testValidationContext("x", "y"))
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

type unknownExpression struct{}

func (u *unknownExpression) expressionNode() {}

func TestValidateUnknownExpression(t *testing.T) {
	expr := &unknownExpression{}

	errors := Validate(expr, DefaultValidationContext())
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}

func TestValidateFunctionCallWithoutArguments(t *testing.T) {
	expr := &FunctionCallExpression{
		Name:      "sin",
		Arguments: nil,
	}

	errors := Validate(expr, DefaultValidationContext())
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}

func TestValidateFunctionCallWithTooManyArguments(t *testing.T) {
	expr := &FunctionCallExpression{
		Name: "cosh",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
			&VariableExpression{Name: "y"},
		},
	}

	errors := Validate(expr, testValidationContext("x", "y"))
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}

	if errors[0].Message == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestValidateBuiltInConstantPi(t *testing.T) {
	expr := &VariableExpression{Name: "Pi"}

	errors := Validate(expr, DefaultValidationContext())
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateBuiltInConstantE(t *testing.T) {
	expr := &VariableExpression{Name: "E"}

	errors := Validate(expr, DefaultValidationContext())
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateUnknownVariableExpression(t *testing.T) {
	expr := &VariableExpression{Name: "x"}

	errors := Validate(expr, testValidationContext("y"))
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}
}

func TestValidateBuiltInConstantAndUserVariable(t *testing.T) {
	expr := &BinaryExpression{
		Left:     &VariableExpression{Name: "Pi"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "r"},
	}

	errors := Validate(expr, testValidationContext("r"))
	if len(errors) != 0 {
		t.Fatalf("expected no errors, got %d", len(errors))
	}
}

func TestValidateUnknownVariableExpressionMessage(t *testing.T) {
	expr := &VariableExpression{Name: "foo"}

	errors := Validate(expr, DefaultValidationContext())
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errors))
	}

	expected := `unknown identifier "foo"`
	if errors[0].Message != expected {
		t.Fatalf("expected error %q, got %q", expected, errors[0].Message)
	}
}
