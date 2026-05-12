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
