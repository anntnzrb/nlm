package batchexecute

import "testing"

func TestReqIDGeneratorReset(t *testing.T) {
	gen := NewReqIDGenerator()
	first := gen.Next()
	second := gen.Next()
	if first == second {
		t.Fatalf("expected different ids")
	}

	gen.Reset()
	reset := gen.Next()
	if reset != first {
		t.Fatalf("expected reset id %q, got %q", first, reset)
	}
}
