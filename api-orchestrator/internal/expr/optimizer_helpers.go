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
