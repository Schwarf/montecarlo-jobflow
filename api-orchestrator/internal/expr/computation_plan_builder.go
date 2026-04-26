package expr

import (
	"fmt"
)

type assignment struct {
	Name string
	Expr Expression
}

type ComputationPlanBuilder struct {
	tempCounter int
	Assignments []assignment
}

func (b *ComputationPlanBuilder) newTempVariable() string {
	b.tempCounter++
	return fmt.Sprintf("h%d", b.tempCounter)
}

func (b *ComputationPlanBuilder) assignToTempVariable(expr Expression) *VariableExpression {
	name := b.newTempVariable()
	b.Assignments = append(b.Assignments, assignment{
		Name: name,
		Expr: expr,
	})
	return &VariableExpression{Name: name}
}

func (b *ComputationPlanBuilder) assignNonTrivialToTempVariable(expr Expression) Expression {
	if IsTrivial(expr) {
		return expr
	}
	return b.assignToTempVariable(expr)
}

func (b *ComputationPlanBuilder) buildPositiveIntegerPower(base Expression, n int) Expression {
	if n <= 0 {
		panic("buildPositiveIntegerPower requires n >= 1")
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
		return b.assignToTempVariable(mul)
	}

	if (n & 1) == 0 {
		half := b.buildPositiveIntegerPower(base, n/2)
		mul := &BinaryExpression{
			Left:     half,
			Operator: TokenMultiply,
			Right:    half,
		}
		return b.assignToTempVariable(mul)
	}

	prev := b.buildPositiveIntegerPower(base, n-1)
	mul := &BinaryExpression{
		Left:     prev,
		Operator: TokenMultiply,
		Right:    base,
	}
	return b.assignToTempVariable(mul)
}

func (b *ComputationPlanBuilder) buildIntegerPower(base Expression, n int) Expression {
	if n == 0 {
		panic("power 0 not implemented")
	}

	if n > 0 {
		return b.buildPositiveIntegerPower(base, n)
	}

	positive := b.buildPositiveIntegerPower(base, -n)

	div := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    positive,
	}
	return b.assignToTempVariable(div)
}

func (b *ComputationPlanBuilder) Build(expr Expression) Expression {
	switch e := expr.(type) {
	case *BinaryExpression:
		if e.Operator == TokenPower {
			n, ok := IntegerLiteralValue(e.Right)
			if ok && n != 0 {
				base := b.Build(e.Left)
				base = b.assignNonTrivialToTempVariable(base)
				return b.buildIntegerPower(base, n)
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
