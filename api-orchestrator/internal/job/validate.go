package job

import (
	"fmt"
	"unicode"
)

func ValidateForComputation(j Job) error {
	if j.Integrand == "" {
		return fmt.Errorf("integrand must not be empty")
	}
	if len(j.IntegrationVariables) == 0 {
		return fmt.Errorf("at least one variable is required")
	}
	if j.Evaluations <= 0 {
		return fmt.Errorf("evaluations must be > 0")
	}

	seen := make(map[string]struct{}, len(j.IntegrationVariables))

	for _, v := range j.IntegrationVariables {
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

func isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}

	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) {
				return false
			}
		}
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}
