package expr

import "testing"

func TestCppCodeGeneratorGenerateFunctionWithoutVariables(t *testing.T) {
	generator := CppCodeGenerator{}

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
	generator := CppCodeGenerator{}

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
	generator := CppCodeGenerator{}

	_, err := generator.GenerateFunction("", []string{"x"}, "x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
