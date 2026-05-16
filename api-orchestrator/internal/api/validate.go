package api

import (
	"fmt"
	"unicode"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/expr"
	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

func (r *CreateJobRequest) ValidateBasic() error {
	if r.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	j := job.Job{
		Integrand:            r.Integrand,
		IntegrationVariables: r.IntegrationVariables,
		Evaluations:          r.Evaluations,
	}

	return job.ValidateForComputation(j)
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
