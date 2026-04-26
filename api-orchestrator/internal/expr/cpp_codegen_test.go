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

func TestCppCodeGeneratorGenerateFunctionCallExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name: "sin",
			expr: &FunctionCallExpression{
				Name: "sin",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::sin(x)",
		},
		{
			name: "cos",
			expr: &FunctionCallExpression{
				Name: "cos",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::cos(x)",
		},
		{
			name: "tan",
			expr: &FunctionCallExpression{
				Name: "tan",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::tan(x)",
		},
		{
			name: "asin",
			expr: &FunctionCallExpression{
				Name: "asin",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::asin(x)",
		},
		{
			name: "acos",
			expr: &FunctionCallExpression{
				Name: "acos",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::acos(x)",
		},
		{
			name: "atan",
			expr: &FunctionCallExpression{
				Name: "atan",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::atan(x)",
		},
		{
			name: "sinh",
			expr: &FunctionCallExpression{
				Name: "sinh",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::sinh(x)",
		},
		{
			name: "cosh",
			expr: &FunctionCallExpression{
				Name: "cosh",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::cosh(x)",
		},
		{
			name: "tanh",
			expr: &FunctionCallExpression{
				Name: "tanh",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::tanh(x)",
		},
		{
			name: "asinh",
			expr: &FunctionCallExpression{
				Name: "asinh",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::asinh(x)",
		},
		{
			name: "acosh",
			expr: &FunctionCallExpression{
				Name: "acosh",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::acosh(x)",
		},
		{
			name: "atanh",
			expr: &FunctionCallExpression{
				Name: "atanh",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::atanh(x)",
		},
		{
			name: "natural logarithm",
			expr: &FunctionCallExpression{
				Name: "ln",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::log(x)",
		},
		{
			name: "logarithm to base 10",
			expr: &FunctionCallExpression{
				Name: "log10",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::log10(x)",
		},
		{
			name: "logarithm to base 2",
			expr: &FunctionCallExpression{
				Name: "log2",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::log2(x)",
		},
		{
			name: "exponential function",
			expr: &FunctionCallExpression{
				Name: "exp",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
			expected: "std::exp(x)",
		},
		{
			name: "function call with binary expression argument",
			expr: &FunctionCallExpression{
				Name: "exp",
				Arguments: []Expression{
					&BinaryExpression{
						Left:     &VariableExpression{Name: "x"},
						Operator: TokenPlus,
						Right:    &NumberExpression{Value: "1.0"},
					},
				},
			},
			expected: "std::exp((x + 1.0))",
		},
		{
			name: "nested function call",
			expr: &FunctionCallExpression{
				Name: "sin",
				Arguments: []Expression{
					&FunctionCallExpression{
						Name: "ln",
						Arguments: []Expression{
							&VariableExpression{Name: "x"},
						},
					},
				},
			},
			expected: "std::sin(std::log(x))",
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

func TestCppCodeGeneratorRejectsUnsupportedFunction(t *testing.T) {
	generator := &CppCodeGenerator{}

	_, err := generator.GenerateExpression(&FunctionCallExpression{
		Name: "foo",
		Arguments: []Expression{
			&VariableExpression{Name: "x"},
		},
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCppCodeGeneratorGenerateBuiltInConstants(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name:     "Pi",
			expr:     &VariableExpression{Name: "Pi"},
			expected: "3.141592653589793238462643383279502884",
		},
		{
			name:     "E",
			expr:     &VariableExpression{Name: "E"},
			expected: "2.718281828459045235360287471352662498",
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
