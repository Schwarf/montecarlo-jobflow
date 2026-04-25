package expr

import "testing"

func TestIntegerLiteralValueNumber(t *testing.T) {
	expr := &NumberExpression{Value: "2"}

	val, ok := IntegerLiteralValue(expr)
	if !ok {
		t.Fatal("expected integer literal, got ok=false")
	}
	if val != 2 {
		t.Fatalf("expected 2, got %d", val)
	}
}

func TestIntegerLiteralValueUnaryMinus(t *testing.T) {
	expr := &UnaryExpression{
		Operator: TokenMinus,
		Right:    &NumberExpression{Value: "16"},
	}

	val, ok := IntegerLiteralValue(expr)
	if !ok {
		t.Fatal("expected integer literal, got ok=false")
	}
	if val != -16 {
		t.Fatalf("expected -16, got %d", val)
	}
}

func TestIntegerLiteralValueUnaryPlus(t *testing.T) {
	expr := &UnaryExpression{
		Operator: TokenPlus,
		Right:    &NumberExpression{Value: "3"},
	}

	value, ok := IntegerLiteralValue(expr)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if value != 3 {
		t.Fatalf("expected 3, got %d", value)
	}
}

func TestIntegerLiteralValueFloatFails(t *testing.T) {
	expr := &NumberExpression{Value: "2.5"}

	_, ok := IntegerLiteralValue(expr)
	if ok {
		t.Fatal("expected ok=false")
	}
}

func TestIntegerLiteralValueVariableFails(t *testing.T) {
	expr := &VariableExpression{Name: "x"}

	_, ok := IntegerLiteralValue(expr)
	if ok {
		t.Fatal("expected ok=false")
	}
}

func TestIntegerLiteralValueFractionFails(t *testing.T) {
	expr := &BinaryExpression{
		Left:     &NumberExpression{Value: "4"},
		Operator: TokenDivide,
		Right:    &NumberExpression{Value: "3"},
	}

	_, ok := IntegerLiteralValue(expr)
	if ok {
		t.Fatal("expected ok=false for fraction expression 4/3")
	}
}

func TestTrivialValidCases(t *testing.T) {
	number := &NumberExpression{Value: "4"}
	variable := &VariableExpression{Name: "x"}

	if !IsTrivial(number) {
		t.Fatal("expected trivial expression: number")
	}
	if !IsTrivial(variable) {
		t.Fatal("expected trivial expression: variable")
	}

}

func TestTrivialInvalidCases(t *testing.T) {
	unary := &UnaryExpression{Right: &NumberExpression{Value: "4"}, Operator: TokenMinus}
	binary := &BinaryExpression{Right: &NumberExpression{Value: "4"}, Operator: TokenMinus, Left: unary}
	function := &FunctionCallExpression{Name: "sin", Arguments: []Expression{&VariableExpression{Name: "x"}}}

	if IsTrivial(unary) {
		t.Fatal("expected nontrivial expression: unary")
	}

	if IsTrivial(binary) {
		t.Fatal("expected nontrivial expression: binary")
	}
	if IsTrivial(function) {
		t.Fatal("expected nontrivial expression: function")
	}
}
