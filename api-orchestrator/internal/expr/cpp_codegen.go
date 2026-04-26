package expr

import (
	"fmt"
	"strings"
)

type CppCodeGenerator struct{}

func (g CppCodeGenerator) GenerateFunction(functionName string, variableNames []string, returnExpression string) (string, error) {
	if functionName == "" {
		return "", fmt.Errorf("function name must not be empty")
	}

	var builder strings.Builder

	builder.WriteString("double ")
	builder.WriteString(functionName)
	builder.WriteString("(")

	for i, variableName := range variableNames {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString("double ")
		builder.WriteString(variableName)
	}

	builder.WriteString(") {\n")
	builder.WriteString("    return ")
	builder.WriteString(returnExpression)
	builder.WriteString(";\n")
	builder.WriteString("}\n")

	return builder.String(), nil
}
