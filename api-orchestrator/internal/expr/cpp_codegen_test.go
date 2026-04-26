package expr

import "testing"

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
			expected: "(-4.0)",
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

// TODO: Exact string testing ... rather brittle
func TestCppCodeGeneratorGenerateComputationPlanFunction(t *testing.T) {
	generator := &CppCodeGenerator{}

	assignments := []assignment{
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
				Operator: TokenPlus,
				Right:    &VariableExpression{Name: "y"},
			},
		},
	}

	result := &BinaryExpression{
		Left:     &VariableExpression{Name: "h2"},
		Operator: TokenMultiply,
		Right:    &VariableExpression{Name: "h2"},
	}

	code, err := generator.GenerateComputationPlanFunction(
		"evaluate",
		[]string{"x", "y"},
		assignments,
		result,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "template <int dimension>\n" +
		"double evaluate(const std::array<double, dimension>& sample, void* param) {\n" +
		"    (void)param;\n\n" +
		"    const double x = sample[0];\n" +
		"    const double y = sample[1];\n\n" +
		"    const double h1 = (x * x);\n" +
		"    const double h2 = (h1 + y);\n" +
		"    return (h2 * h2);\n" +
		"}\n"

	if code != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, code)
	}
}

// TODO: Exact string testing ... rather brittle
func TestCppCodegenPipelineRepeatedOptimizedSubexpressions(t *testing.T) {
	expr, err := parseForTest(t, "(1+(x+y)^2)^3/sin(x)^2 + (1+(x+y)^2)^2")
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	planBuilder := &ComputationPlanBuilder{}
	result := planBuilder.Build(expr)

	generator := &CppCodeGenerator{}

	code, err := generator.GenerateComputationPlanFunction(
		"evaluate",
		[]string{"x", "y"},
		planBuilder.Assignments,
		result,
	)
	if err != nil {
		t.Fatalf("unexpected codegen error: %v", err)
	}

	expected := "template <int dimension>\n" +
		"double evaluate(const std::array<double, dimension>& sample, void* param) {\n" +
		"    (void)param;\n\n" +
		"    const double x = sample[0];\n" +
		"    const double y = sample[1];\n\n" +
		"    const double h1 = (x + y);\n" +
		"    const double h2 = (h1 * h1);\n" +
		"    const double h3 = (1.0 + h2);\n" +
		"    const double h4 = (h3 * h3);\n" +
		"    const double h5 = (h4 * h3);\n" +
		"    const double h6 = std::sin(x);\n" +
		"    const double h7 = (h6 * h6);\n" +
		"    return ((h5 / h7) + h4);\n" +
		"}\n"

	if code != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, code)
	}
}

func TestCppCodeGeneratorGenerateGenericPowerExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected string
	}{
		{
			name: "symbolic exponent",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenPower,
				Right:    &VariableExpression{Name: "y"},
			},
			expected: "std::pow(x, y)",
		},
		{
			name: "non integer exponent",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenPower,
				Right:    &NumberExpression{Value: "2.3482"},
			},
			expected: "std::pow(x, 2.3482)",
		},
		{
			name: "fractional negative exponent",
			expr: &BinaryExpression{
				Left:     &VariableExpression{Name: "x"},
				Operator: TokenPower,
				Right: &BinaryExpression{
					Left: &UnaryExpression{
						Operator: TokenMinus,
						Right:    &NumberExpression{Value: "4"},
					},
					Operator: TokenDivide,
					Right:    &NumberExpression{Value: "3"},
				},
			},
			expected: "std::pow(x, ((-4.0) / 3.0))",
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

func TestCppCodeGeneratorGenerateSource(t *testing.T) {
	generator := &CppCodeGenerator{}

	assignments := []assignment{
		{
			Name: "h1",
			Expr: &FunctionCallExpression{
				Name: "sin",
				Arguments: []Expression{
					&VariableExpression{Name: "x"},
				},
			},
		},
	}

	result := &VariableExpression{Name: "h1"}

	code, err := generator.GenerateIntegrandHeader(
		"evaluate",
		[]string{"x"},
		assignments,
		result,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "#pragma once\n\n" +
		"#include <array>\n" +
		"#include <cmath>\n\n" +
		"template <int dimension>\n" +
		"double evaluate(const std::array<double, dimension>& sample, void* param) {\n" +
		"    (void)param;\n\n" +
		"    const double x = sample[0];\n\n" +
		"    const double h1 = std::sin(x);\n" +
		"    return h1;\n" +
		"}\n"

	if code != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, code)
	}
}
