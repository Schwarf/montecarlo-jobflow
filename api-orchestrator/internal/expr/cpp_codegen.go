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

func cppFunctionName(name string) (string, error) {
	switch name {
	case "sin":
		return "std::sin", nil
	case "cos":
		return "std::cos", nil
	case "tan":
		return "std::tan", nil
	case "asin":
		return "std::asin", nil
	case "acos":
		return "std::acos", nil
	case "atan":
		return "std::atan", nil
	case "sinh":
		return "std::sinh", nil
	case "cosh":
		return "std::cosh", nil
	case "tanh":
		return "std::tanh", nil
	case "asinh":
		return "std::asinh", nil
	case "acosh":
		return "std::acosh", nil
	case "atanh":
		return "std::atanh", nil
	case "ln":
		return "std::log", nil
	case "log10":
		return "std::log10", nil
	case "log2":
		return "std::log2", nil
	case "exp":
		return "std::exp", nil
	default:
		return "", fmt.Errorf("unsupported function %q", name)
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
