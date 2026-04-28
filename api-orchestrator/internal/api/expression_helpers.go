package api

import "github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/expr"

func (r *CreateJobRequest) VariableNames() []string {
	names := make([]string, 0, len(r.IntegrationVariables))
	for _, variable := range r.IntegrationVariables {
		names = append(names, variable.Name)
	}
	return names
}

func (r *CreateJobRequest) ExpressionValidationContext() expr.ValidationContext {
	context := expr.DefaultValidationContext()
	for _, variable := range r.IntegrationVariables {
		context.UserVariables[variable.Name] = struct{}{}
	}
	return context
}
