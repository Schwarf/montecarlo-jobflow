package expr

import "fmt"

type ValidationError struct {
	Message string
}

func Validate(expr Expression) []ValidationError {
	var errors []ValidationError
	validateExpr(expr, &errors)
	return errors
}

func validateExpr(expr Expression, errors *[]ValidationError) {
	switch e := expr.(type) {
	case *NumberExpression:
		return
	case *VariableExpression:
		return
	case *UnaryExpression:
		validateExpr(e.Right, errors)
	case *BinaryExpression:
		validateExpr(e.Left, errors)
		validateExpr(e.Right, errors)
	case *FunctionCallExpression:
		for _, arg := range e.Arguments {
			validateExpr(arg, errors)
		}
	default:
		*errors = append(*errors, ValidationError{
			Message: fmt.Sprintf("unknown expression type %T", expr),
		})
	}
}
