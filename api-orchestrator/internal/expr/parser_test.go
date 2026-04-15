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
	expr, err := parseForTest(t, "1+2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := binaryExpr.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left to be *NumberExpression, got %T", binaryExpr.Left)
	}

	right, ok := binaryExpr.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right to be *NumberExpression, got %T", binaryExpr.Right)
	}

	if left.Value != "1" {
		t.Fatalf("expected left.Value %q, got %q", "1", left.Value)
	}

	if right.Value != "2" {
		t.Fatalf("expected right.Value %q, got %q", "2", right.Value)
	}

	if binaryExpr.Operator != TokenPlus {
		t.Fatalf("expected operator %v, got %v", TokenPlus, binaryExpr.Operator)
	}
}

func TestParseSubtractIdentifierFromNumber(t *testing.T) {
	expr, err := parseForTest(t, "(2-x)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := binaryExpr.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left to be *NumberExpression, got %T", binaryExpr.Left)
	}

	right, ok := binaryExpr.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right to be *VariableExpression, got %T", binaryExpr.Right)
	}

	if left.Value != "2" {
		t.Fatalf("expected left.Value %q, got %q", "2", left.Value)
	}

	if right.Name != "x" {
		t.Fatalf("expected right.Name %q, got %q", "x", right.Name)
	}

	if binaryExpr.Operator != TokenMinus {
		t.Fatalf("expected operator %q, got %q", TokenMinus, binaryExpr.Operator)
	}
}

func TestParseAssociativity(t *testing.T) {
	expr, err := parseForTest(t, "a+b-c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := binaryExpr.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected left to be *BinaryExpression, got %T", binaryExpr.Left)
	}

	right, ok := binaryExpr.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right to be *VariableExpression, got %T", binaryExpr.Right)
	}

	leftLeft, ok := left.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left.Left to be *VariableExpression, got %T", left.Left)
	}

	leftRight, ok := left.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected leftRight to be *VariableExpression, got %T", leftRight)
	}

	if right.Name != "c" {
		t.Fatalf("expected left %q, got %q", "c", right.Name)
	}

	if binaryExpr.Operator != TokenMinus {
		t.Fatalf("expected operator %v, got %v", TokenMinus, binaryExpr.Operator)
	}

	if leftLeft.Name != "a" {
		t.Fatalf("expected leftLeft.Name %q, got %q", "a", leftLeft.Name)
	}

	if leftRight.Name != "b" {
		t.Fatalf("expected leftRight.Name %q, got %q", "b", leftRight.Name)
	}

	if left.Operator != TokenPlus {
		t.Fatalf("expected left.Operator %v, got %v", TokenPlus, left.Operator)
	}

}

func TestParseMissingRightOperandAfterPlusFails(t *testing.T) {
	_, err := parseForTest(t, "x+a+1+")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMissingRightOperandAfterMinusFails(t *testing.T) {
	_, err := parseForTest(t, "2-")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMultiplicationNumbers(t *testing.T) {
	expr, err := parseForTest(t, "2*3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := binaryExpr.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left to be *NumberExpression, got %T", binaryExpr.Left)
	}

	right, ok := binaryExpr.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right to be *NumberExpression, got %T", binaryExpr.Right)
	}

	if left.Value != "2" {
		t.Fatalf("expected left.Value %q, got %q", "2", left.Value)
	}

	if right.Value != "3" {
		t.Fatalf("expected right.Value %q, got %q", "3", right.Value)
	}

	if binaryExpr.Operator != TokenMultiply {
		t.Fatalf("expected operator %v, got %v", TokenMultiply, binaryExpr.Operator)
	}
}

func TestParseDivisionIdentifierByNumber(t *testing.T) {
	expr, err := parseForTest(t, "x/2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := binaryExpr.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left to be *VariableExpression, got %T", binaryExpr.Left)
	}

	right, ok := binaryExpr.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right to be *NumberExpression, got %T", binaryExpr.Right)
	}

	if left.Name != "x" {
		t.Fatalf("expected left.Name %q, got %q", "x", left.Name)
	}

	if right.Value != "2" {
		t.Fatalf("expected right.Value %q, got %q", "2", right.Value)
	}

	if binaryExpr.Operator != TokenDivide {
		t.Fatalf("expected operator %v, got %v", TokenDivide, binaryExpr.Operator)
	}
}

func TestParseMultiplicationLeftAssociative(t *testing.T) {
	expr, err := parseForTest(t, "a*b*c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := root.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected root.Left to be *BinaryExpression, got %T", root.Left)
	}

	right, ok := root.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected root.Right to be *VariableExpression, got %T", root.Right)
	}

	leftLeft, ok := left.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left.Left to be *VariableExpression, got %T", left.Left)
	}

	leftRight, ok := left.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left.Right to be *VariableExpression, got %T", left.Right)
	}

	if root.Operator != TokenMultiply {
		t.Fatalf("expected root.Operator %v, got %v", TokenMultiply, root.Operator)
	}

	if left.Operator != TokenMultiply {
		t.Fatalf("expected left.Operator %v, got %v", TokenMultiply, left.Operator)
	}

	if leftLeft.Name != "a" {
		t.Fatalf("expected leftLeft.Name %q, got %q", "a", leftLeft.Name)
	}

	if leftRight.Name != "b" {
		t.Fatalf("expected leftRight.Name %q, got %q", "b", leftRight.Name)
	}

	if right.Name != "c" {
		t.Fatalf("expected right.Name %q, got %q", "c", right.Name)
	}
}

func TestParseDivisionLeftAssociative(t *testing.T) {
	expr, err := parseForTest(t, "a/b/c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := root.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected root.Left to be *BinaryExpression, got %T", root.Left)
	}

	right, ok := root.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected root.Right to be *VariableExpression, got %T", root.Right)
	}

	leftLeft, ok := left.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left.Left to be *VariableExpression, got %T", left.Left)
	}

	leftRight, ok := left.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left.Right to be *VariableExpression, got %T", left.Right)
	}

	if root.Operator != TokenDivide {
		t.Fatalf("expected root.Operator %v, got %v", TokenDivide, root.Operator)
	}

	if left.Operator != TokenDivide {
		t.Fatalf("expected left.Operator %v, got %v", TokenDivide, left.Operator)
	}

	if leftLeft.Name != "a" {
		t.Fatalf("expected leftLeft.Name %q, got %q", "a", leftLeft.Name)
	}

	if leftRight.Name != "b" {
		t.Fatalf("expected leftRight.Name %q, got %q", "b", leftRight.Name)
	}

	if right.Name != "c" {
		t.Fatalf("expected right.Name %q, got %q", "c", right.Name)
	}
}

func TestParsePrecedenceMultiplyBeforePlus(t *testing.T) {
	expr, err := parseForTest(t, "2+3*4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := root.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected root.Left to be *NumberExpression, got %T", root.Left)
	}

	right, ok := root.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected root.Right to be *BinaryExpression, got %T", root.Right)
	}

	rightLeft, ok := right.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right.Left to be *NumberExpression, got %T", right.Left)
	}

	rightRight, ok := right.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right.Right to be *NumberExpression, got %T", right.Right)
	}

	if root.Operator != TokenPlus {
		t.Fatalf("expected root.Operator %v, got %v", TokenPlus, root.Operator)
	}

	if right.Operator != TokenMultiply {
		t.Fatalf("expected right.Operator %v, got %v", TokenMultiply, right.Operator)
	}

	if left.Value != "2" {
		t.Fatalf("expected left.Value %q, got %q", "2", left.Value)
	}

	if rightLeft.Value != "3" {
		t.Fatalf("expected rightLeft.Value %q, got %q", "3", rightLeft.Value)
	}

	if rightRight.Value != "4" {
		t.Fatalf("expected rightRight.Value %q, got %q", "4", rightRight.Value)
	}
}

func TestParsePrecedenceMultiplyBeforeMinus(t *testing.T) {
	expr, err := parseForTest(t, "a-b*c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := root.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected root.Left to be *VariableExpression, got %T", root.Left)
	}

	right, ok := root.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected root.Right to be *BinaryExpression, got %T", root.Right)
	}

	rightLeft, ok := right.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right.Left to be *VariableExpression, got %T", right.Left)
	}

	rightRight, ok := right.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right.Right to be *VariableExpression, got %T", right.Right)
	}

	if root.Operator != TokenMinus {
		t.Fatalf("expected root.Operator %v, got %v", TokenMinus, root.Operator)
	}

	if right.Operator != TokenMultiply {
		t.Fatalf("expected right.Operator %v, got %v", TokenMultiply, right.Operator)
	}

	if left.Name != "a" {
		t.Fatalf("expected left.Name %q, got %q", "a", left.Name)
	}

	if rightLeft.Name != "b" {
		t.Fatalf("expected rightLeft.Name %q, got %q", "b", rightLeft.Name)
	}

	if rightRight.Name != "c" {
		t.Fatalf("expected rightRight.Name %q, got %q", "c", rightRight.Name)
	}
}

func TestParseParenthesizedAdditionTimesNumber(t *testing.T) {
	expr, err := parseForTest(t, "(1+2)*3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	root, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	left, ok := root.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected root.Left to be *BinaryExpression, got %T", root.Left)
	}

	right, ok := root.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected root.Right to be *NumberExpression, got %T", root.Right)
	}

	leftLeft, ok := left.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left.Left to be *NumberExpression, got %T", left.Left)
	}

	leftRight, ok := left.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left.Right to be *NumberExpression, got %T", left.Right)
	}

	if root.Operator != TokenMultiply {
		t.Fatalf("expected root.Operator %v, got %v", TokenMultiply, root.Operator)
	}

	if left.Operator != TokenPlus {
		t.Fatalf("expected left.Operator %v, got %v", TokenPlus, left.Operator)
	}

	if leftLeft.Value != "1" {
		t.Fatalf("expected leftLeft.Value %q, got %q", "1", leftLeft.Value)
	}

	if leftRight.Value != "2" {
		t.Fatalf("expected leftRight.Value %q, got %q", "2", leftRight.Value)
	}

	if right.Value != "3" {
		t.Fatalf("expected right.Value %q, got %q", "3", right.Value)
	}
}

func TestParseMultiplyMissingRightOperandFails(t *testing.T) {
	_, err := parseForTest(t, "2*")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseDivideMissingRightOperandFails(t *testing.T) {
	_, err := parseForTest(t, "x/")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMultiplyMissingLeftOperandFails(t *testing.T) {
	_, err := parseForTest(t, "*2")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseDivideMissingLeftOperandFails(t *testing.T) {
	_, err := parseForTest(t, "/2")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMalformedMixedOperatorsFails(t *testing.T) {
	_, err := parseForTest(t, "2+*3")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseEmptyParenthesizedTermFails(t *testing.T) {
	_, err := parseForTest(t, "()")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseMissingRightParenAfterTermFails(t *testing.T) {
	_, err := parseForTest(t, "(2*3")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseUnaryMinusVariable(t *testing.T) {
	expr, err := parseForTest(t, "-x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unaryExpr, ok := expr.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected *UnaryExpression, got %T", expr)
	}

	if unaryExpr.Operator != TokenMinus {
		t.Fatalf("expected operator TokenMinus, got %v", unaryExpr.Operator)
	}

	right, ok := unaryExpr.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right to be *VariableExpression, got %T", unaryExpr.Right)
	}

	if right.Name != "x" {
		t.Fatalf("expected variable name x, got %q", right.Name)
	}
}

func TestParseUnaryMinusNumber(t *testing.T) {
	expr, err := parseForTest(t, "-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unaryExpr, ok := expr.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected *UnaryExpression, got %T", expr)
	}

	if unaryExpr.Operator != TokenMinus {
		t.Fatalf("expected operator TokenMinus, got %v", unaryExpr.Operator)
	}

	right, ok := unaryExpr.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right to be *NumberExpression, got %T", unaryExpr.Right)
	}

	if right.Value != "2" {
		t.Fatalf("expected number value 2, got %q", right.Value)
	}
}

func TestParseUnaryMinusParenthesizedExpression(t *testing.T) {
	expr, err := parseForTest(t, "-(a+b)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unaryExpr, ok := expr.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected *UnaryExpression, got %T", expr)
	}

	if unaryExpr.Operator != TokenMinus {
		t.Fatalf("expected operator TokenMinus, got %v", unaryExpr.Operator)
	}

	right, ok := unaryExpr.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected right to be *BinaryExpression, got %T", unaryExpr.Right)
	}

	if right.Operator != TokenPlus {
		t.Fatalf("expected operator TokenPlus, got %v", right.Operator)
	}
}

func TestParseMultiplyByUnaryMinus(t *testing.T) {
	expr, err := parseForTest(t, "3*-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	if binaryExpr.Operator != TokenMultiply {
		t.Fatalf("expected operator TokenMultiply, got %v", binaryExpr.Operator)
	}

	left, ok := binaryExpr.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left to be *NumberExpression, got %T", binaryExpr.Left)
	}

	if left.Value != "3" {
		t.Fatalf("expected left value 3, got %q", left.Value)
	}

	right, ok := binaryExpr.Right.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected right to be *UnaryExpression, got %T", binaryExpr.Right)
	}

	if right.Operator != TokenMinus {
		t.Fatalf("expected unary operator TokenMinus, got %v", right.Operator)
	}

	rightNumber, ok := right.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected unary right to be *NumberExpression, got %T", right.Right)
	}

	if rightNumber.Value != "2" {
		t.Fatalf("expected unary right value 2, got %q", rightNumber.Value)
	}
}

func TestParseDivideByUnaryMinus(t *testing.T) {
	expr, err := parseForTest(t, "x/-y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	if binaryExpr.Operator != TokenDivide {
		t.Fatalf("expected operator TokenDivide, got %v", binaryExpr.Operator)
	}

	left, ok := binaryExpr.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left to be *VariableExpression, got %T", binaryExpr.Left)
	}

	if left.Name != "x" {
		t.Fatalf("expected left name x, got %q", left.Name)
	}

	right, ok := binaryExpr.Right.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected right to be *UnaryExpression, got %T", binaryExpr.Right)
	}

	if right.Operator != TokenMinus {
		t.Fatalf("expected unary operator TokenMinus, got %v", right.Operator)
	}

	rightVar, ok := right.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected unary right to be *VariableExpression, got %T", right.Right)
	}

	if rightVar.Name != "y" {
		t.Fatalf("expected unary right name y, got %q", rightVar.Name)
	}
}

func TestParseDoubleUnaryMinus(t *testing.T) {
	expr, err := parseForTest(t, "--x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outer, ok := expr.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected outer to be *UnaryExpression, got %T", expr)
	}

	if outer.Operator != TokenMinus {
		t.Fatalf("expected outer operator TokenMinus, got %v", outer.Operator)
	}

	inner, ok := outer.Right.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected inner to be *UnaryExpression, got %T", outer.Right)
	}

	if inner.Operator != TokenMinus {
		t.Fatalf("expected inner operator TokenMinus, got %v", inner.Operator)
	}

	rightVar, ok := inner.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected inner right to be *VariableExpression, got %T", inner.Right)
	}

	if rightVar.Name != "x" {
		t.Fatalf("expected variable name x, got %q", rightVar.Name)
	}
}

func TestParseUnaryPlusVariable(t *testing.T) {
	expr, err := parseForTest(t, "+x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unaryExpr, ok := expr.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected *UnaryExpression, got %T", expr)
	}

	if unaryExpr.Operator != TokenPlus {
		t.Fatalf("expected operator TokenPlus, got %v", unaryExpr.Operator)
	}

	right, ok := unaryExpr.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right to be *VariableExpression, got %T", unaryExpr.Right)
	}

	if right.Name != "x" {
		t.Fatalf("expected variable name x, got %q", right.Name)
	}
}

func TestParseMultiplyByUnaryPlus(t *testing.T) {
	expr, err := parseForTest(t, "x*+y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	if binaryExpr.Operator != TokenMultiply {
		t.Fatalf("expected operator TokenMultiply, got %v", binaryExpr.Operator)
	}

	right, ok := binaryExpr.Right.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected right to be *UnaryExpression, got %T", binaryExpr.Right)
	}

	if right.Operator != TokenPlus {
		t.Fatalf("expected unary operator TokenPlus, got %v", right.Operator)
	}
}

func TestParsePowerSimple(t *testing.T) {
	expr, err := parseForTest(t, "2^3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	binaryExpr, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected *BinaryExpression, got %T", expr)
	}

	if binaryExpr.Operator != TokenPower {
		t.Fatalf("expected operator TokenPower, got %v", binaryExpr.Operator)
	}

	left, ok := binaryExpr.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected left to be *NumberExpression, got %T", binaryExpr.Left)
	}
	if left.Value != "2" {
		t.Fatalf("expected left value 2, got %q", left.Value)
	}

	right, ok := binaryExpr.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected right to be *NumberExpression, got %T", binaryExpr.Right)
	}
	if right.Value != "3" {
		t.Fatalf("expected right value 3, got %q", right.Value)
	}
}

func TestParsePowerRightAssociative(t *testing.T) {
	expr, err := parseForTest(t, "2^3^4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outer, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer to be *BinaryExpression, got %T", expr)
	}

	if outer.Operator != TokenPower {
		t.Fatalf("expected outer operator TokenPower, got %v", outer.Operator)
	}

	left, ok := outer.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected outer left to be *NumberExpression, got %T", outer.Left)
	}
	if left.Value != "2" {
		t.Fatalf("expected outer left value 2, got %q", left.Value)
	}

	right, ok := outer.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer right to be *BinaryExpression, got %T", outer.Right)
	}

	if right.Operator != TokenPower {
		t.Fatalf("expected inner operator TokenPower, got %v", right.Operator)
	}

	innerLeft, ok := right.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected inner left to be *NumberExpression, got %T", right.Left)
	}
	if innerLeft.Value != "3" {
		t.Fatalf("expected inner left value 3, got %q", innerLeft.Value)
	}

	innerRight, ok := right.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected inner right to be *NumberExpression, got %T", right.Right)
	}
	if innerRight.Value != "4" {
		t.Fatalf("expected inner right value 4, got %q", innerRight.Value)
	}
}

func TestParseMultiplyWithPowerPrecedence(t *testing.T) {
	expr, err := parseForTest(t, "2*3^4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outer, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer to be *BinaryExpression, got %T", expr)
	}

	if outer.Operator != TokenMultiply {
		t.Fatalf("expected outer operator TokenMultiply, got %v", outer.Operator)
	}

	left, ok := outer.Left.(*NumberExpression)
	if !ok {
		t.Fatalf("expected outer left to be *NumberExpression, got %T", outer.Left)
	}
	if left.Value != "2" {
		t.Fatalf("expected outer left value 2, got %q", left.Value)
	}

	right, ok := outer.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer right to be *BinaryExpression, got %T", outer.Right)
	}

	if right.Operator != TokenPower {
		t.Fatalf("expected right operator TokenPower, got %v", right.Operator)
	}
}

func TestParsePowerBeforeMultiply(t *testing.T) {
	expr, err := parseForTest(t, "2^3*4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outer, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer to be *BinaryExpression, got %T", expr)
	}

	if outer.Operator != TokenMultiply {
		t.Fatalf("expected outer operator TokenMultiply, got %v", outer.Operator)
	}

	left, ok := outer.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer left to be *BinaryExpression, got %T", outer.Left)
	}

	if left.Operator != TokenPower {
		t.Fatalf("expected left operator TokenPower, got %v", left.Operator)
	}

	right, ok := outer.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected outer right to be *NumberExpression, got %T", outer.Right)
	}
	if right.Value != "4" {
		t.Fatalf("expected outer right value 4, got %q", right.Value)
	}
}

func TestParsePowerWithParenthesizedExponent(t *testing.T) {
	expr, err := parseForTest(t, "x^(-4/3)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outer, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer to be *BinaryExpression, got %T", expr)
	}

	if outer.Operator != TokenPower {
		t.Fatalf("expected operator TokenPower, got %v", outer.Operator)
	}

	left, ok := outer.Left.(*VariableExpression)
	if !ok {
		t.Fatalf("expected left to be *VariableExpression, got %T", outer.Left)
	}
	if left.Name != "x" {
		t.Fatalf("expected variable name x, got %q", left.Name)
	}

	right, ok := outer.Right.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected right to be *BinaryExpression, got %T", outer.Right)
	}

	if right.Operator != TokenDivide {
		t.Fatalf("expected exponent operator TokenDivide, got %v", right.Operator)
	}

	rightLeft, ok := right.Left.(*UnaryExpression)
	if !ok {
		t.Fatalf("expected exponent left to be *UnaryExpression, got %T", right.Left)
	}
	if rightLeft.Operator != TokenMinus {
		t.Fatalf("expected unary minus in exponent left, got %v", rightLeft.Operator)
	}

	rightLeftNumber, ok := rightLeft.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected unary right to be *NumberExpression, got %T", rightLeft.Right)
	}
	if rightLeftNumber.Value != "4" {
		t.Fatalf("expected unary right value 4, got %q", rightLeftNumber.Value)
	}

	rightRight, ok := right.Right.(*NumberExpression)
	if !ok {
		t.Fatalf("expected exponent right to be *NumberExpression, got %T", right.Right)
	}
	if rightRight.Value != "3" {
		t.Fatalf("expected exponent right value 3, got %q", rightRight.Value)
	}
}

func TestParsePowerWithParenthesizedBase(t *testing.T) {
	expr, err := parseForTest(t, "(a+b)^c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outer, ok := expr.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected outer to be *BinaryExpression, got %T", expr)
	}

	if outer.Operator != TokenPower {
		t.Fatalf("expected operator TokenPower, got %v", outer.Operator)
	}

	left, ok := outer.Left.(*BinaryExpression)
	if !ok {
		t.Fatalf("expected left to be *BinaryExpression, got %T", outer.Left)
	}
	if left.Operator != TokenPlus {
		t.Fatalf("expected left operator TokenPlus, got %v", left.Operator)
	}

	right, ok := outer.Right.(*VariableExpression)
	if !ok {
		t.Fatalf("expected right to be *VariableExpression, got %T", outer.Right)
	}
	if right.Name != "c" {
		t.Fatalf("expected variable name c, got %q", right.Name)
	}
}

func TestParsePowerMissingRightOperand(t *testing.T) {
	_, err := parseForTest(t, "2^")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerMissingLeftOperand(t *testing.T) {
	_, err := parseForTest(t, "^2")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerMissingExponentAfterUnaryMinus(t *testing.T) {
	_, err := parseForTest(t, "2^-")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerEmptyParenthesizedExponent(t *testing.T) {
	_, err := parseForTest(t, "x^()")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerMissingClosingParenInExponent(t *testing.T) {
	_, err := parseForTest(t, "x^(-4/3")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerUnexpectedTokenAfterExponent(t *testing.T) {
	_, err := parseForTest(t, "2^3)")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerDoubleOperator(t *testing.T) {
	_, err := parseForTest(t, "2^^3")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParsePowerOperatorWithoutValidRightSide(t *testing.T) {
	_, err := parseForTest(t, "2^*3")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
