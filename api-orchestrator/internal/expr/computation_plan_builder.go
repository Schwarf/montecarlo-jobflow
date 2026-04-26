package expr

import (
	"fmt"
)

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

func (b *ComputationPlanBuilder) AssignNonTrivialToTempVariable(expr Expression) Expression {
	if IsTrivial(expr) {
		return expr
	}
	return b.AssignToTempVariable(expr)
}

func (b *ComputationPlanBuilder) BuildPositiveIntegerPower(base Expression, n int) Expression {
	if n <= 0 {
		panic("BuildPositiveIntegerPower requires n >= 1")
	}

	if n == 1 {
		return base
	}

	if n == 2 {
		mul := &BinaryExpression{
			Left:     base,
			Operator: TokenMultiply,
			Right:    base,
		}
		return b.AssignToTempVariable(mul)
	}

	if (n & 1) == 0 {
		half := b.BuildPositiveIntegerPower(base, n/2)
		mul := &BinaryExpression{
			Left:     half,
			Operator: TokenMultiply,
			Right:    half,
		}
		return b.AssignToTempVariable(mul)
	}

	prev := b.BuildPositiveIntegerPower(base, n-1)
	mul := &BinaryExpression{
		Left:     prev,
		Operator: TokenMultiply,
		Right:    base,
	}
	return b.AssignToTempVariable(mul)
}

func (b *ComputationPlanBuilder) BuildIntegerPower(base Expression, n int) Expression {
	if n == 0 {
		panic("power 0 not implemented")
	}

	if n > 0 {
		return b.BuildPositiveIntegerPower(base, n)
	}

	positive := b.BuildPositiveIntegerPower(base, -n)

	div := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    positive,
	}
	return b.AssignToTempVariable(div)
}

func (b *ComputationPlanBuilder) Build(expr Expression) Expression {
	switch e := expr.(type) {
	case *BinaryExpression:
		if e.Operator == TokenPower {
			n, ok := IntegerLiteralValue(e.Right)
			if ok && n != 0 {
				base := b.Build(e.Left)
				base = b.AssignNonTrivialToTempVariable(base)
				return b.BuildIntegerPower(base, n)
			}
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
