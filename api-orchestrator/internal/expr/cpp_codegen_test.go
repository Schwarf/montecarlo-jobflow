package expr

import "testing"

func TestCppCodeGeneratorGenerateFunctionWithoutVariables(t *testing.T) {
	generator := &CppCodeGenerator{}

	code, err := generator.GenerateFunction("evaluate", nil, "1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "double evaluate() {\n" +
		"    return 1.0;\n" +
		"}\n"

	if code != expected {
		t.Fatalf("unexpected generated code:\nexpected:\n%s\nactual:\n%s", expected, code)
	}
}

func TestCppCodeGeneratorGenerateFunctionWithVariables(t *testing.T) {
	generator := &CppCodeGenerator{}

	code, err := generator.GenerateFunction("evaluate", []string{"x", "y", "z"}, "x + y + z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "double evaluate(double x, double y, double z) {\n" +
		"    return x + y + z;\n" +
		"}\n"

	if code != expected {
		t.Fatalf("unexpected generated code:\nexpected:\n%s\nactual:\n%s", expected, code)
	}
}

func TestCppCodeGeneratorRejectsEmptyFunctionName(t *testing.T) {
	generator := &CppCodeGenerator{}

	_, err := generator.GenerateFunction("", []string{"x"}, "x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCppCodeGeneratorGenerateNumberExpression(t *testing.T) {
	generator := &CppCodeGenerator{}

	code, err := generator.GenerateExpression(&NumberExpression{Value: "1.23"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "1.23"
	if code != expected {
		t.Fatalf("expected %q, got %q", expected, code)
	}
}

func TestCppCodeGeneratorGenerateVariableExpression(t *testing.T) {
	generator := &CppCodeGenerator{}

	code, err := generator.GenerateExpression(&VariableExpression{Name: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "x"
	if code != expected {
		t.Fatalf("expected %q, got %q", expected, code)
	}
}

func TestCppCodeGeneratorGenerateBinaryExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name: "addition",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenPlus,
				Right:    &VariableExpression{Name: "y"},
			},
			expected: "(x + y)",
		},
		{
			name: "subtraction",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenMinus,
				Right:    &NumberExpression{Value: "1.0"},
			},
			expected: "(x - 1.0)",
		},
		{
			name: "multiplication",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "y"},
			},
			expected: "(x * y)",
		},
		{
			name: "division",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenDivide,
				Right:    &VariableExpression{Name: "y"},
			},
			expected: "(x / y)",
		},
		{
			name: "nested expression",
			expr: &BinaryExpression{
				Left: &BinaryExpression{
					Left:     &VariableExpression{Name: "x"},
					Operator: TokenPlus,
					Right:    &VariableExpression{Name: "y"},
				},
				Operator: TokenMultiply,
				Right:    &VariableExpression{Name: "z"},
			},
			expected: "((x + y) * z)",
		},
	}

	generator := &CppCodeGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := generator.GenerateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if code != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, code)
			}
		})
	}
}

func TestCppCodeGeneratorGenerateUnaryExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name: "unary plus",
			expr: &UnaryExpression{
				Operator: TokenPlus,
				Right:    &VariableExpression{Name: "x"},
			},
			expected: "(+x)",
		},
		{
			name: "unary minus",
			expr: &UnaryExpression{
				Operator: TokenMinus,
				Right:    &VariableExpression{Name: "x"},
			},
			expected: "(-x)",
		},
		{
			name: "negative number",
			expr: &UnaryExpression{
				Operator: TokenMinus,
				Right:    &NumberExpression{Value: "4"},
			},
			expected: "(-4)",
		},
		{
			name: "nested unary around binary",
			expr: &UnaryExpression{
				Operator: TokenMinus,
				Right: &BinaryExpression{
					Left:     &VariableExpression{Name: "x"},
					Operator: TokenPlus,
					Right:    &VariableExpression{Name: "y"},
				},
			},
			expected: "(-(x + y))",
		},
	}

	generator := &CppCodeGenerator{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := generator.GenerateExpression(tt.expr)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if code != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, code)
			}
		})
	}
}
