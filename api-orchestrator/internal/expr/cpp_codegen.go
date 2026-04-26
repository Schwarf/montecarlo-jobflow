package expr

import (
	"fmt"
	"strings"
)

type CppCodeGenerator struct{}

func (g *CppCodeGenerator) GenerateExpression(expr Expression) (string, error) {
	switch e := expr.(type) {
	case *NumberExpression:
		return e.Value, nil
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
