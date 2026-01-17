package batchexecute

import (
	"strings"
	"testing"
)

func TestSanitizeJSONInvalidEscapes(t *testing.T) {
	input := `["bad\xescape","good\nline","unicode\u123G"]`
	got := sanitizeJSON(input)

	// Invalid \x and invalid unicode should be escaped to preserve JSON.
	if !strings.Contains(got, `\\xescape`) {
		t.Fatalf("expected invalid escape to be sanitized: %s", got)
	}
	if !strings.Contains(got, `\\u123G`) {
		t.Fatalf("expected invalid unicode escape to be sanitized: %s", got)
	}
	// Valid escapes should remain.
	if !strings.Contains(got, `\nline`) {
		t.Fatalf("expected valid escape to remain: %s", got)
	}
}

func TestNormalizeResponseArray(t *testing.T) {
	if _, err := normalizeResponseArray("not an array"); err == nil {
		t.Fatalf("expected error for non-array response")
	}

	decoded := []interface{}{
		[]interface{}{"wrb.fr", "id", "[1,2,3]"},
	}
	out, err := normalizeResponseArray(decoded)
	if err != nil {
		t.Fatalf("normalizeResponseArray error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 response, got %d", len(out))
	}
}

func TestDecodeResponseNumeric(t *testing.T) {
	if _, err := decodeResponseWithOptions("277567", false); err == nil {
		t.Fatalf("expected error for numeric-only response without chunked parsing")
	}
}

func TestDecodeResponseEmpty(t *testing.T) {
	if _, err := decodeResponseWithOptions(")]}'", false); err == nil {
		t.Fatalf("expected error for empty response")
	}
}
