package expr

import "testing"

func parseForTest(t *testing.T, input string) (Expression, error) {
	t.Helper()

	tokens, err := LexAll(input)
	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)
	return parser.Parse()
}

func TestParseIdentifier(t *testing.T) {
	expr, err := parseForTest(t, "x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variable, ok := expr.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", expr)
	}

	if variable.Name != "x" {
		t.Fatalf("expected name %q, got %q", "x", variable.Name)
	}
}

func TestParseIdentifierWithDigit(t *testing.T) {
	expr, err := parseForTest(t, "x2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variable, ok := expr.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", expr)
	}

	if variable.Name != "x2" {
		t.Fatalf("expected name %q, got %q", "x2", variable.Name)
	}
}

func TestParseIdentifierWithUnderscoreAndDigit(t *testing.T) {
	expr, err := parseForTest(t, "x_11")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variable, ok := expr.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", expr)
	}

	if variable.Name != "x_11" {
		t.Fatalf("expected name %q, got %q", "x_11", variable.Name)
	}
}

func TestParseInteger(t *testing.T) {
	expr, err := parseForTest(t, "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	num, ok := expr.(*NumberExpression)
	if !ok {
		t.Fatalf("expected *NumberExpression, got %T", expr)
	}

	if num.Value != "42" {
		t.Fatalf("expected value %q, got %q", "42", num.Value)
	}
}

func TestParseScientificNumber(t *testing.T) {
	expr, err := parseForTest(t, "1.23E-3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := expr.(*NumberExpression); !ok {
		t.Fatalf("expected *NumberExpression, got %T", expr)
	}
}

func TestParseParenthesizedNumber(t *testing.T) {
	expr, err := parseForTest(t, "(7)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	num, ok := expr.(*NumberExpression)
	if !ok {
		t.Fatalf("expected *NumberExpression, got %T", expr)
	}

	if num.Value != "7" {
		t.Fatalf("expected value %q, got %q", "7", num.Value)
	}
}

func TestParseDoubleParenthesizedNumber(t *testing.T) {
	expr, err := parseForTest(t, "((7))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	num, ok := expr.(*NumberExpression)
	if !ok {
		t.Fatalf("expected *NumberExpression, got %T", expr)
	}

	if num.Value != "7" {
		t.Fatalf("expected value %q, got %q", "7", num.Value)
	}
}

func TestParseParenthesizedIdentifier(t *testing.T) {
	expr, err := parseForTest(t, "(x)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variable, ok := expr.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", expr)
	}

	if variable.Name != "x" {
		t.Fatalf("expected name %q, got %q", "x", variable.Name)
	}
}

func TestParseDoubleParenthesizedIdentifier(t *testing.T) {
	expr, err := parseForTest(t, "((x_0))")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	variable, ok := expr.(*VariableExpression)
	if !ok {
		t.Fatalf("expected *VariableExpression, got %T", expr)
	}

	if variable.Name != "x_0" {
		t.Fatalf("expected name %q, got %q", "x_0", variable.Name)
	}
}

func TestParseNumberWithTrailingTokenFails(t *testing.T) {
	_, err := parseForTest(t, "42 x")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseIdentifierWithTrailingTokenFails(t *testing.T) {
	_, err := parseForTest(t, "x 7")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMissingRightParenFails(t *testing.T) {
	_, err := parseForTest(t, "(1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMissingLeftParenFails(t *testing.T) {
	_, err := parseForTest(t, "x)")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseEmptyInputFails(t *testing.T) {
	_, err := parseForTest(t, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseUnexpectedTokenFails(t *testing.T) {
	_, err := parseForTest(t, "!")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseAdditionNumbers(t *testing.T) {
	expr, err := parseForTest(t, "(1+2)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := binaryExpr.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected Left to be *NumberExpression, got %T", binaryExpr.Left)
	}

	right, ok := binaryExpr.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected Right to be *NumberExpression, got %T", binaryExpr.Left)
	}

	if left.Value != "1" {
		t.Fatalf("expected left %q, got %q", "1", left.Value)
	}

	if right.Value != "2" {
		t.Fatalf("expected right %q, got %q", "1", right.Value)
	}

	if binaryExpr.Operator.String() != "+" {
		t.Fatalf("expected operator %q, got %q", "+", binaryExpr.Operator.String())
	}

}
