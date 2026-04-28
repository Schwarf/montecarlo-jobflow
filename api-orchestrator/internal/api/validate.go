package api

import (
	"fmt"
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
		// ToDo The lower bound shall be smaller than the upper bound. if not integrand gets a global minus
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
	_, err := expr.ParseAndValidate(r.Integrand, r.ExpressionValidationContext())
	return err
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
