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
