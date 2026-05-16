package jobdispatch

import (
	"fmt"

	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/expr"
	"github.com/Schwarf/montecarlo-jobflow/api-orchestrator/internal/job"
)

const IntegrandFunctionName = "integrand"

type HeaderCompiler interface {
	CompileToHeader(
		functionName string,
		integrand string,
		variableNames []string,
		context expr.ValidationContext) (string, error)
}

type Preparer struct {
	Compiler HeaderCompiler
}

func NewPreparer() Preparer {
	return Preparer{
		Compiler: expr.NewIntegrandCompiler(),
	}
}

func (p Preparer) Prepare(j job.Job) (WorkerSubmission, error) {
	if p.Compiler == nil {
		return WorkerSubmission{}, fmt.Errorf("compiler must not be nil")
	}

	if err := validateJobForDispatch(j); err != nil {
		return WorkerSubmission{}, err
	}

	variableNames := make([]string, 0, len(j.IntegrationVariables))
	context := expr.DefaultValidationContext()

	for _, v := range j.IntegrationVariables {
		variableNames = append(variableNames, v.Name)
		context.UserVariables[v.Name] = struct{}{}
	}

	header, err := p.Compiler.CompileToHeader(
		IntegrandFunctionName,
		j.Integrand,
		variableNames,
		context,
	)
	if err != nil {
		return WorkerSubmission{}, fmt.Errorf("compile integrand header: %w", err)
	}

	return WorkerSubmission{
		JobID:        j.ID,
		FunctionName: IntegrandFunctionName,
		Integrand:    j.Integrand,
		Variables:    append([]job.VariableSpec(nil), j.IntegrationVariables...),
		Evaluations:  j.Evaluations,
		Header:       header,
	}, nil
}

func validateJobForDispatch(j job.Job) error {
	if j.ID == "" {
		return fmt.Errorf("job id must not be empty")
	}

	return job.ValidateForComputation(j)
}
