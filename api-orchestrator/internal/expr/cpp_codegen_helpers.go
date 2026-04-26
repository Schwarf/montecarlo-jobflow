package expr

import "fmt"

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

func cppIdentifier(name string) string {
	switch name {
	case "Pi":
		return "3.141592653589793238462643383279502884"
	case "E":
		return "2.718281828459045235360287471352662498"
	default:
		return name
	}
}
