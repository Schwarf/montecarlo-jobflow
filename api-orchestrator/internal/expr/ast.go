package expr

type Expression interface {
	expressionNode()
}

type NumberExpression struct {
	Value string
}

func (n *NumberExpression) expressionNode() {}

type VariableExpression struct {
	Name string
}

func (v *VariableExpression) expressionNode() {}
