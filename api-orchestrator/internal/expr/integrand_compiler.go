package expr

type IntegrandCompiler struct {
	CodeGenerator *CppCodeGenerator
}

func NewIntegrandCompiler() *IntegrandCompiler {
	return &IntegrandCompiler{
		CodeGenerator: &CppCodeGenerator{},
	}
}

func (c *IntegrandCompiler) CompileToHeader(
	functionName string,
	integrand string,
	variableNames []string,
	context ValidationContext,
) (string, error) {
	expression, err := ParseAndValidate(integrand, context)
	if err != nil {
		return "", err
	}

	planBuilder := &ComputationPlanBuilder{}
	result := planBuilder.Build(expression)

	return c.CodeGenerator.GenerateIntegrandHeader(
		functionName,
		variableNames,
		planBuilder.Assignments,
		result,
	)
}

// TODO no TESTs
