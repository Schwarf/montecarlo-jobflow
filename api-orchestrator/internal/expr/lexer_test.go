package expr

import "testing"

func collectTokenTypes(tokens []Token) []TokenType {
	result := make([]TokenType, 0, len(tokens))
	for _, tok := range tokens {
		result = append(result, tok.Type)
	}
	return result
}
func collectLexemes(tokens []Token) []string {
	result := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		result = append(result, tok.Lexeme)
	}
	return result
}

func TestLexerSingleIdentifier(t *testing.T) {
	tokens, err := LexAll("x")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	got := collectTokenTypes(tokens)
	want := []TokenType{
		TokenIdentifier,
		TokenEOF,
	}

	if len(got) != len(want) {
		t.Fatalf("token count mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("token %d mismatch: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestLexerBasicOperators(t *testing.T) {
	tokens, err := LexAll("-z+x+y*z/2^a_1")
	if err != nil {
		t.Fatalf("LexAll returned error: %v", err)
	}
	got := collectTokenTypes(tokens)
	want := []TokenType{
		TokenMinus,      // -
		TokenIdentifier, // z
		TokenPlus,       // +
		TokenIdentifier, // x
		TokenPlus,       // +
		TokenIdentifier, // y
		TokenMultiply,   // *
		TokenIdentifier, // z
		TokenDivide,     // /
		TokenNumber,     // 2
		TokenExponent,   // ^
		TokenIdentifier, // a
		TokenEOF,
	}

	if len(got) != len(want) {
		t.Fatalf("token count mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("token %d mismatch: got %v, want %v", i, got[i], want[i])
		}
	}
	gotLexemes := collectLexemes(tokens)
	wantLexemes := []string{
		"-",
		"z",
		"+",
		"x",
		"+",
		"y",
		"*",
		"z",
		"/",
		"2",
		"^",
		"a_1",
		"",
	}

	for i := range wantLexemes {
		if gotLexemes[i] != wantLexemes[i] {
			t.Fatalf("lexeme %d mismatch: got %q, want %q", i, gotLexemes[i], wantLexemes[i])
		}
	}

}

func TestLexerParenthesesAndUnaryMinus(t *testing.T) {
	tokens, err := LexAll("x^(-4/3)")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	gotTypes := collectTokenTypes(tokens)
	wantTypes := []TokenType{
		TokenIdentifier,
		TokenExponent,
		TokenLeftParen,
		TokenMinus,
		TokenNumber,
		TokenDivide,
		TokenNumber,
		TokenRightParen,
		TokenEOF,
	}

	if len(gotTypes) != len(wantTypes) {
		t.Fatalf("token count mismatch: got %d, want %d", len(gotTypes), len(wantTypes))
	}

	for i := range wantTypes {
		if gotTypes[i] != wantTypes[i] {
			t.Fatalf("token %d mismatch: got %v, want %v", i, gotTypes[i], wantTypes[i])
		}
	}

	gotLexemes := collectLexemes(tokens)
	wantLexemes := []string{
		"x",
		"^",
		"(",
		"-",
		"4",
		"/",
		"3",
		")",
		"",
	}

	for i := range wantLexemes {
		if gotLexemes[i] != wantLexemes[i] {
			t.Fatalf("lexeme %d mismatch: got %q, want %q", i, gotLexemes[i], wantLexemes[i])
		}
	}
}

func TestLexerFunctionCall(t *testing.T) {
	tokens, err := LexAll("sin(x)")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	got := collectTokenTypes(tokens)
	want := []TokenType{
		TokenIdentifier,
		TokenLeftParen,
		TokenIdentifier,
		TokenRightParen,
		TokenEOF,
	}

	if len(got) != len(want) {
		t.Fatalf("token count mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("token %d mismatch: got %v, want %v", i, got[i], want[i])
		}
	}
	gotLexemes := collectLexemes(tokens)
	wantLexemes := []string{
		"sin",
		"(",
		"x",
		")",
		"",
	}

	for i := range wantLexemes {
		if gotLexemes[i] != wantLexemes[i] {
			t.Fatalf("lexeme %d mismatch: got %q, want %q", i, gotLexemes[i], wantLexemes[i])
		}
	}

}

func TestLexerWhitespaceIsIgnored(t *testing.T) {
	tokens, err := LexAll("  \t x  +   y \n ")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	got := collectTokenTypes(tokens)
	want := []TokenType{
		TokenIdentifier,
		TokenPlus,
		TokenIdentifier,
		TokenEOF,
	}

	if len(got) != len(want) {
		t.Fatalf("token count mismatch: got %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("token %d mismatch: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestLexer_Number_Integer(t *testing.T) {
	tokens, err := LexAll("42")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("token count mismatch: got %d, want %d", len(tokens), 2)
	}

	if tokens[0].Type != TokenNumber {
		t.Fatalf("got token type %v, want %v", tokens[0].Type, TokenNumber)
	}

	if tokens[0].Lexeme != "42" {
		t.Fatalf("got lexeme %q, want %q", tokens[0].Lexeme, "42")
	}
}

func TestLexer_Number_Decimal(t *testing.T) {
	tokens, err := LexAll("3.1415")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("token count mismatch: got %d, want %d", len(tokens), 2)
	}

	if tokens[0].Type != TokenNumber {
		t.Fatalf("got token type %v, want %v", tokens[0].Type, TokenNumber)
	}

	if tokens[0].Lexeme != "3.1415" {
		t.Fatalf("got lexeme %q, want %q", tokens[0].Lexeme, "3.1415")
	}
}

func TestLexer_Number_LeadingDot(t *testing.T) {
	tokens, err := LexAll(".5")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("token count mismatch: got %d, want %d", len(tokens), 2)
	}

	if tokens[0].Type != TokenNumber {
		t.Fatalf("got token type %v, want %v", tokens[0].Type, TokenNumber)
	}

	if tokens[0].Lexeme != ".5" {
		t.Fatalf("got lexeme %q, want %q", tokens[0].Lexeme, ".5")
	}
}

func TestLexer_Number_ScientificNotation(t *testing.T) {
	tokens, err := LexAll("1.2e-3")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("token count mismatch: got %d, want %d", len(tokens), 2)
	}

	if tokens[0].Type != TokenNumber {
		t.Fatalf("got token type %v, want %v", tokens[0].Type, TokenNumber)
	}

	if tokens[0].Lexeme != "1.2e-3" {
		t.Fatalf("got lexeme %q, want %q", tokens[0].Lexeme, "1.2e-3")
	}

	// Other version of exponents
	tokens, err = LexAll("4.567E+8")
	if err != nil {
		t.Fatalf("lexAll returned error: %v", err)
	}

	if len(tokens) != 2 {
		t.Fatalf("token count mismatch: got %d, want %d", len(tokens), 2)
	}

	if tokens[0].Type != TokenNumber {
		t.Fatalf("got token type %v, want %v", tokens[0].Type, TokenNumber)
	}

	if tokens[0].Lexeme != "4.567E+8" {
		t.Fatalf("got lexeme %q, want %q", tokens[0].Lexeme, "4.567E+8")
	}
}
