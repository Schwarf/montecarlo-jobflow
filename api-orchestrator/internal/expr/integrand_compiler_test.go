package expr

import "testing"

func TestIntegrandCompilerCompileToHeader(t *testing.T) {
	context := DefaultValidationContext()
	context.UserVariables["x"] = struct{}{}
	context.UserVariables["y"] = struct{}{}

	compiler := NewIntegrandCompiler()

	header, err := compiler.CompileToHeader(
		"generated_integrand",
		"(x+y)^2",
		[]string{"x", "y"},
		context,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "#pragma once\n\n" +
		"#include <array>\n" +
		"#include <cmath>\n\n" +
		"template <int dimension>\n" +
		"double generated_integrand(const std::array<double, dimension>& sample, void* param) {\n" +
		"    (void)param;\n\n" +
		"    const double x = sample[0];\n" +
		"    const double y = sample[1];\n\n" +
		"    const double h1 = (x + y);\n" +
		"    const double h2 = (h1 * h1);\n" +
		"    return h2;\n" +
		"}\n"

	if header != expected {
		t.Fatalf("expected:\n%s\ngot:\n%s", expected, header)
	}
}

func TestIntegrandCompilerCompileToHeaderRejectsInvalidExpression(t *testing.T) {
	context := DefaultValidationContext()
	context.UserVariables["x"] = struct{}{}

	compiler := NewIntegrandCompiler()

	_, err := compiler.CompileToHeader(
		"generated_integrand",
		"x + unknown",
		[]string{"x"},
		context,
	)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
