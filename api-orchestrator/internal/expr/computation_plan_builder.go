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

func (b *ComputationPlanBuilder) Build(expr Expression) Expression {
	switch e := expr.(type) {
	case *BinaryExpression:
		square, ok := b.BuildSquare(e)
		if ok {
			return square
		}
		left := b.Build(e.Left)
		right := b.Build(e.Right)
		return &BinaryExpression{
			Left:     left,
			Operator: e.Operator,
			Right:    right,
		}
	case *UnaryExpression:
		right := b.Build(e.Right)
		return &UnaryExpression{
			Operator: e.Operator,
			Right:    right,
		}
	case *FunctionCallExpression:
		result := FunctionCallExpression{Name: e.Name}
		for _, arg := range e.Arguments {
			res := b.Build(arg)
			result.Arguments = append(result.Arguments, res)
		}
		return &result
	default:
		return expr
	}
}
