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
	memo        map[string]*VariableExpression
}

func (b *ComputationPlanBuilder) newTempVariable() string {
	b.tempCounter++
	return fmt.Sprintf("h%d", b.tempCounter)
}

func (b *ComputationPlanBuilder) assignOrReuseTempVariable(expr Expression) *VariableExpression {
	if b.memo == nil {
		b.memo = make(map[string]*VariableExpression)
	}
	key := ExpressionKey(expr)
	if val, ok := b.memo[key]; ok {
		return val
	}

	name := b.newTempVariable()
	b.Assignments = append(b.Assignments, assignment{
		Name: name,
		Expr: expr,
	})
	val := &VariableExpression{Name: name}
	b.memo[key] = val
	return val
}

func (b *ComputationPlanBuilder) assignNonTrivialToTempVariable(expr Expression) Expression {
	if IsTrivial(expr) {
		return expr
	}
	return b.assignOrReuseTempVariable(expr)
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
		return b.assignOrReuseTempVariable(mul)
	}

	if (n & 1) == 0 {
		half := b.buildPositiveIntegerPower(base, n/2)
		mul := &BinaryExpression{
			Left:     half,
			Operator: TokenMultiply,
			Right:    half,
		}
		return b.assignOrReuseTempVariable(mul)
	}

	prev := b.buildPositiveIntegerPower(base, n-1)
	mul := &BinaryExpression{
		Left:     prev,
		Operator: TokenMultiply,
		Right:    base,
	}
	return b.assignOrReuseTempVariable(mul)
}

func (b *ComputationPlanBuilder) buildIntegerPower(base Expression, n int) Expression {
	if n == 0 {
		return &NumberExpression{Value: "1"}
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
	return b.assignOrReuseTempVariable(div)
}

func (b *ComputationPlanBuilder) Build(expr Expression) Expression {
	switch e := expr.(type) {
	case *BinaryExpression:
		if e.Operator == TokenPower {
			n, ok := IntegerLiteralValue(e.Right)
			if ok {
				if n == 0 {
					return b.buildIntegerPower(nil, 0)
				}
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
