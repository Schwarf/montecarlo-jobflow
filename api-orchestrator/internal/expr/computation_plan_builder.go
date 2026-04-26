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

func (b *ComputationPlanBuilder) SimplifyPowerOfOne(expr *BinaryExpression) (Expression, bool) {
	if expr.Operator != TokenPower {
		return nil, false
	}

	val, ok := IntegerLiteralValue(expr.Right)
	if !ok || val != 1 {
		return nil, false
	}

	return b.Build(expr.Left), true
}

func (b *ComputationPlanBuilder) SimplifyPowerOfMinusOne(expr *BinaryExpression) (Expression, bool) {
	if expr.Operator != TokenPower {
		return nil, false
	}

	val, ok := IntegerLiteralValue(expr.Right)
	if !ok || val != -1 {
		return nil, false
	}

	base := b.Build(expr.Left)
	base = b.AssignNonTrivialToTempVariable(base)

	div := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    base,
	}

	return b.AssignToTempVariable(div), true
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

func (b *ComputationPlanBuilder) BuildSquare(expr *BinaryExpression) (Expression, bool) {
	if expr.Operator != TokenPower {
		return nil, false
	}

	val, ok := IntegerLiteralValue(expr.Right)
	if !ok || val != 2 {
		return nil, false
	}

	base := b.Build(expr.Left)
	base = b.AssignNonTrivialToTempVariable(base)

	return b.BuildPositiveIntegerPower(base, val), true
}

func (b *ComputationPlanBuilder) BuildInverseSquare(expr *BinaryExpression) (Expression, bool) {
	if expr.Operator != TokenPower {
		return nil, false
	}

	val, ok := IntegerLiteralValue(expr.Right)
	if !ok || val != -2 {
		return nil, false
	}

	squareExpr := &BinaryExpression{
		Left:     expr.Left,
		Operator: TokenPower,
		Right:    &NumberExpression{Value: "2"},
	}

	squared, ok := b.BuildSquare(squareExpr)
	if !ok {
		return nil, false
	}

	div := &BinaryExpression{
		Left:     &NumberExpression{Value: "1"},
		Operator: TokenDivide,
		Right:    squared,
	}

	return b.AssignToTempVariable(div), true
}

func (b *ComputationPlanBuilder) Build(expr Expression) Expression {
	switch e := expr.(type) {
	case *BinaryExpression:
		one, ok := b.SimplifyPowerOfOne(e)
		if ok {
			return one
		}
		inverse, ok := b.SimplifyPowerOfMinusOne(e)
		if ok {
			return inverse
		}
		if e.Operator == TokenPower {
			n, ok := IntegerLiteralValue(e.Right)
			if ok && n > 0 {
				base := b.Build(e.Left)
				base = b.AssignNonTrivialToTempVariable(base)
				return b.BuildPositiveIntegerPower(base, n)
			}
		}

		inverseSquare, ok := b.BuildInverseSquare(e)
		if ok {
			return inverseSquare
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
