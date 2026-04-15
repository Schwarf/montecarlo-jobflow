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

func TestLexer_SingleIdentifier(t *testing.T) {
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

func TestLexer_BasicOperators(t *testing.T) {
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

func TestLexer_ParenthesesAndUnaryMinus(t *testing.T) {
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

func TestLexer_FunctionCall(t *testing.T) {
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

func TestLexer_WhitespaceIsIgnored(t *testing.T) {
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
