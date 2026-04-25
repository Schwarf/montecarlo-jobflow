package expr

import "fmt"

type Assignment struct {
	Name string
	Expr Expression
}

type ComputationPlanBuilder struct {
	tempCounter int
	Assignments []Assignment
}

func (b *ComputationPlanBuilder) NewTempVariable() string {
	b.tempCounter++
	return fmt.Sprintf("h%d", b.tempCounter)
}

func (b *ComputationPlanBuilder) Emit(expr Expression) *VariableExpression {
	name := b.NewTempVariable()
	b.Assignments = append(b.Assignments, Assignment{
		Name: name,
		Expr: expr,
	})
	return &VariableExpression{Name: name}
}
