package expr

import (
	"fmt"
	"strings"
)

type CppCodeGenerator struct{}

func cppBinaryOperator(operator TokenType) (string, error) {
	switch operator {
	case TokenPlus:
		return "+", nil
	case TokenMinus:
		return "-", nil
	case TokenMultiply:
		return "*", nil
	case TokenDivide:
		return "/", nil
	default:
		return "", fmt.Errorf("unsupported binary operator %v", operator)
	}
}

func cppUnaryOperator(operator TokenType) (string, error) {
	switch operator {
	case TokenPlus:
		return "+", nil
	case TokenMinus:
		return "-", nil
	default:
		return "", fmt.Errorf("unsupported unary operator %v", operator)
	}
}

func (g *CppCodeGenerator) GenerateFunction(functionName string, variableNames []string, returnExpression string) (string, error) {
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

func (g *CppCodeGenerator) GenerateExpression(expr Expression) (string, error) {
	switch e := expr.(type) {
	case *NumberExpression:
		return e.Value, nil
	case *VariableExpression:
		return e.Name, nil
	case *BinaryExpression:
		left, err := g.GenerateExpression(e.Left)
		if err != nil {
			return "", err
		}

		right, err := g.GenerateExpression(e.Right)
		if err != nil {
			return "", err
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
	default:
		return "", fmt.Errorf("unsupported expression type %T", e)
	}
}
