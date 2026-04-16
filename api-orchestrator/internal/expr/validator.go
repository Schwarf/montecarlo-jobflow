package expr

import "fmt"

type ValidationError struct {
	Message string
}

type ValidationContext struct {
	AllowedFunctions map[string]struct{}
	AllowedConstants map[string]struct{}
}

func DefaultValidationContext() ValidationContext {
	return ValidationContext{
		AllowedFunctions: map[string]struct{}{
			"sin":   {},
			"cos":   {},
			"tan":   {},
			"asin":  {},
			"acos":  {},
			"atan":  {},
			"sinh":  {},
			"cosh":  {},
			"tanh":  {},
			"asinh": {},
			"acosh": {},
			"atanh": {},
			"log10": {},
			"ln":    {},
			"log2":  {},
			"exp":   {},
		},
		AllowedConstants: map[string]struct{}{
			"Pi": {},
			"E":  {},
		},
	}
}

func Validate(expr Expression) []ValidationError {
	context := DefaultValidationContext()
	var errors []ValidationError
	validateExpr(expr, context, &errors)
	return errors
}

func validateExpr(expr Expression, context ValidationContext, errors *[]ValidationError) {
	switch e := expr.(type) {
	case *NumberExpression:
		return
	case *VariableExpression:
		return
	case *UnaryExpression:
		validateExpr(e.Right, context, errors)
	case *BinaryExpression:
		validateExpr(e.Left, context, errors)
		validateExpr(e.Right, context, errors)
	case *FunctionCallExpression:
		if _, ok := context.AllowedFunctions[e.Name]; !ok {
			*errors = append(*errors, ValidationError{
				Message: fmt.Sprintf("unknown function %q", e.Name),
			})
		}
		for _, arg := range e.Arguments {
			validateExpr(arg, context, errors)
		}
	default:
		*errors = append(*errors, ValidationError{
			Message: fmt.Sprintf("unknown expression type %T", expr),
		})
	}
}
