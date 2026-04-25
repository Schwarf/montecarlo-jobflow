package expr

import "testing"

func TestComputationPlanBuilderNewTempVariable(t *testing.T) {
	var b ComputationPlanBuilder

	if got := b.NewTempVariable(); got != "h1" {
		t.Fatalf("expected h1, got %q", got)
	}
	if got := b.NewTempVariable(); got != "h2" {
		t.Fatalf("expected h2, got %q", got)
	}
}
