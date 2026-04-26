package expr

import (
	"fmt"
	"strings"
)

type CppCodeGenerator struct{}

func (g *CppCodeGenerator) GenerateSource(
	functionName string,
	variableNames []string,
	assignments []assignment,
	result Expression,
) (string, error) {
	functionCode, err := g.GenerateComputationPlanFunction(functionName, variableNames, assignments, result)
	if err != nil {
		return "", err
	}
	var builder strings.Builder

	builder.WriteString("#include <cmath>\n\n")
	builder.WriteString(functionCode)
	return builder.String(), nil
}

func (g *CppCodeGenerator) GenerateComputationPlanFunction(
	functionName string,
	variableNames []string,
	assignments []assignment,
	result Expression,
) (string, error) {
	if functionName == "" {
		return "", fmt.Errorf("functionName must not be empty")
	}

	var builder strings.Builder
	// function declaration
	builder.WriteString("double ")
	builder.WriteString(functionName)

	// function arguments
	builder.WriteString("(")

	for i, variableName := range variableNames {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("double ")
		builder.WriteString(variableName)
	}

	// function body
	builder.WriteString(") {\n")
	for _, assign := range assignments {
		rhs, err := g.GenerateExpression(assign.Expr)
		if err != nil {
			return "", err
		}

		builder.WriteString("    const double ")
		builder.WriteString(assign.Name)
		builder.WriteString(" = ")
		builder.WriteString(rhs)
		builder.WriteString(";\n")
	}

	returnExpr, err := g.GenerateExpression(result)
	if err != nil {
		return "", err
	}

	builder.WriteString("    return ")
	builder.WriteString(returnExpr)
	builder.WriteString(";\n")
	builder.WriteString("}\n")

	return builder.String(), nil
}

func (g *CppCodeGenerator) GenerateExpression(expr Expression) (string, error) {
	switch e := expr.(type) {
	case *NumberExpression:
		return cppNumberLiteral(e.Value), nil
	case *VariableExpression:
		return cppIdentifier(e.Name), nil
	case *BinaryExpression:
		left, err := g.GenerateExpression(e.Left)
		if err != nil {
			return "", err
		}

		right, err := g.GenerateExpression(e.Right)
		if err != nil {
			return "", err
		}

		if e.Operator == TokenPower {
			return "std::pow(" + left + ", " + right + ")", nil
		}

		operator, err := cppBinaryOperator(e.Operator)
		if err != nil {
			return "", err
		}
		return "(" + left + " " + operator + " " + right + ")", nil

	case *UnaryExpression:
		right, err := g.GenerateExpression(e.Right)
		if err != nil {
			return "", err
		}

		operator, err := cppUnaryOperator(e.Operator)
		if err != nil {
			return "", err
		}
		return "(" + operator + right + ")", nil

	case *FunctionCallExpression:
		functionName, err := cppFunctionName(e.Name)
		if err != nil {
			return "", err
		}

		args := make([]string, 0, len(e.Arguments))
		for _, arg := range e.Arguments {
			code, err := g.GenerateExpression(arg)
			if err != nil {
				return "", err
			}
			args = append(args, code)
		}
		return functionName + "(" + strings.Join(args, ",") + ")", nil
	default:
		return "", fmt.Errorf("unsupported expression type %T", e)
	}
}
