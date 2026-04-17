package expr

import "fmt"

type ValidationError struct {
	Message string
}

type ValidationContext struct {
	AllowedFunctions map[string]int
	BuiltInConstants map[string]struct{}
	UserVariables    map[string]struct{}
}

func DefaultValidationContext() ValidationContext {
	return ValidationContext{
		AllowedFunctions: map[string]int{
			"sin":   1,
			"cos":   1,
			"tan":   1,
			"asin":  1,
			"acos":  1,
			"atan":  1,
			"sinh":  1,
			"cosh":  1,
			"tanh":  1,
			"asinh": 1,
			"acosh": 1,
			"atanh": 1,
			"log10": 1,
			"ln":    1,
			"log2":  1,
			"exp":   1,
		},
		BuiltInConstants: map[string]struct{}{
			"Pi": {},
			"E":  {},
		},
		UserVariables: map[string]struct{}{},
	}
}

func Validate(expr Expression, context ValidationContext) []ValidationError {
	var errors []ValidationError
	validateExpr(expr, context, &errors)
	return errors
}

func validateExpr(expr Expression, context ValidationContext, errors *[]ValidationError) {
	switch e := expr.(type) {
	case *NumberExpression:
		return
	case *VariableExpression:
		if _, ok := context.BuiltInConstants[e.Name]; ok {
			return
		}
		if _, ok := context.UserVariables[e.Name]; ok {
			return
		}

		*errors = append(*errors, ValidationError{
			Message: fmt.Sprintf("unknown identifier %q", e.Name),
		})
	case *UnaryExpression:
		validateExpr(e.Right, context, errors)
	case *BinaryExpression:
		validateExpr(e.Left, context, errors)
		validateExpr(e.Right, context, errors)
	case *FunctionCallExpression:
		expectedArgCount, ok := context.AllowedFunctions[e.Name]
		if !ok {
			*errors = append(*errors, ValidationError{
				Message: fmt.Sprintf("unknown function %q", e.Name),
			})
		} else if len(e.Arguments) != expectedArgCount {
			*errors = append(*errors, ValidationError{
				Message: fmt.Sprintf("function %q expects %d argument(s), got %d", e.Name, expectedArgCount, len(e.Arguments)),
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
