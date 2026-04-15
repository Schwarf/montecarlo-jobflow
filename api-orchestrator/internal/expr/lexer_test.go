package expr

import "testing"

func collectTokenTypes(tokens []Token) []TokenType {
	result := make([]TokenType, 0, len(tokens))
	for _, tok := range tokens {
		result = append(result, tok.Type)
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
	tokens, err := LexAll("-z+x+y*z/2^a")
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
}
