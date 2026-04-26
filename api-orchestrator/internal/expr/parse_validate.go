package expr

import (
	"fmt"
	"strings"
)

func ParseAndValidate(input string, context ValidationContext) (Expression, error) {
	tokens, err := LexAll(input)
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	validationErrors := Validate(expression, context)
	if len(validationErrors) > 0 {
		var messages []string
		for _, ve := range validationErrors {
			messages = append(messages, ve.Message)
		}
		return nil, fmt.Errorf("semantic validation failed: %s", strings.Join(messages, "; "))
	}

	return expression, nil
}
