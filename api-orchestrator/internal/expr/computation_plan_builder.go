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

func (b *ComputationPlanBuilder) AssignToTempVariable(expr Expression) *VariableExpression {
	name := b.NewTempVariable()
	b.Assignments = append(b.Assignments, Assignment{
		Name: name,
		Expr: expr,
	})
	return &VariableExpression{Name: name}
}

func (b *ComputationPlanBuilder) BuildSquare(expr *BinaryExpression) (Expression, bool) {
	if expr.Operator != TokenPower {
		return nil, false
	}

	val, ok := IntegerLiteralValue(expr.Right)
	if !ok || val != 2 {
		return nil, false
	}

	mul := &BinaryExpression{
		Left:     expr.Left,
		Operator: TokenMultiply,
		Right:    expr.Left,
	}
	result := b.AssignToTempVariable(mul)
	return result, true
}
