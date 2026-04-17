package api

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/expr"
)

func (r *CreateJobRequest) ValidateBasic() error {
	if r.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	if r.Integrand == "" {
		return fmt.Errorf("integrand must not be empty")
	}
	if len(r.IntegrationVariables) == 0 {
		return fmt.Errorf("at least one variable is required")
	}
	if r.Evaluations <= 0 {
		return fmt.Errorf("evaluations must be > 0")
	}

	seen := make(map[string]struct{})
	for _, v := range r.IntegrationVariables {
		if v.Name == "" {
			return fmt.Errorf("variable name must not be empty")
		}
		if !isValidIdentifier(v.Name) {
			return fmt.Errorf("invalid variable name: %q", v.Name)
		}
		if _, ok := seen[v.Name]; ok {
			return fmt.Errorf("duplicate variable name: %q", v.Name)
		}
		seen[v.Name] = struct{}{}

		if v.Lower == "" {
			return fmt.Errorf("variable %q has empty lower bound", v.Name)
		}
		if v.Upper == "" {
			return fmt.Errorf("variable %q has empty upper bound", v.Name)
		}
	}
	return nil
}

func (r *CreateJobRequest) ValidateSemantics() error {
	tokens, err := expr.LexAll(r.Integrand)
	if err != nil {
		return err
	}

	parser := expr.NewParser(tokens)
	expression, err := parser.Parse()
	if err != nil {
		return err
	}

	context := expr.DefaultValidationContext()
	for _, variable := range r.IntegrationVariables {
		context.UserVariables[variable.Name] = struct{}{}
	}

	validationErrors := expr.Validate(expression, context)
	if len(validationErrors) > 0 {
		var messages []string
		for _, ve := range validationErrors {
			messages = append(messages, ve.Message)
		}
		return fmt.Errorf("semantic validation failed: %s", strings.Join(messages, "; "))
	}

	return nil
}

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) {
				return false
			}
			continue
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
