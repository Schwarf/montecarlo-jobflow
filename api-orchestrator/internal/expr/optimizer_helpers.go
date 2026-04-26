package expr

import "strconv"

func IntegerLiteralValue(expr Expression) (int, bool) {
	switch e := expr.(type) {
	case *NumberExpression:
		val, err := strconv.Atoi(e.Value)
		if err != nil {
			return 0, false
		}
		return val, true
	case *UnaryExpression:
		val, ok := IntegerLiteralValue(e.Right)
		if !ok {
			return 0, false
		}
		switch e.Operator {
		case TokenPlus:
			return val, true
		case TokenMinus:
			return -val, true
		default:
			return 0, false
		}
	default:
		return 0, false
	}
}

func IsTrivial(expr Expression) bool {
	switch expr.(type) {
	case *VariableExpression, *NumberExpression:
		return true
	default:
		return false
	}
}

func ExpressionKey(expr Expression) string {
	switch e := expr.(type) {
	case *NumberExpression:
		return "num(" + e.Value + ")"

	case *VariableExpression:
		return "var(" + e.Name + ")"

	case *UnaryExpression:
		return "unary(" + tokenTypeKey(e.Operator) + "," + ExpressionKey(e.Right) + ")"

	case *BinaryExpression:
		return "binary(" + tokenTypeKey(e.Operator) + "," +
			ExpressionKey(e.Left) + "," +
			ExpressionKey(e.Right) + ")"

	case *FunctionCallExpression:
		key := "call(" + e.Name
		for _, arg := range e.Arguments {
			key += "," + ExpressionKey(arg)
		}
		key += ")"
		return key

	default:
		panic("unsupported expression type in ExpressionKey")
	}
}
