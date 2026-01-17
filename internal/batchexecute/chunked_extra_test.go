package batchexecute

import (
	"strconv"
	"strings"
	"testing"
)

func TestDecodeChunkedResponseNumeric(t *testing.T) {
	resp, err := decodeChunkedResponse(strings.NewReader("277567"))
	if err != nil {
		t.Fatalf("decodeChunkedResponse error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resp))
	}
	if resp[0].ID != "numeric" {
		t.Fatalf("expected numeric id, got %q", resp[0].ID)
	}
}

func TestDecodeChunkedResponseValid(t *testing.T) {
	chunk := `[["wrb.fr","id","[1]",null,null,null,"generic"]]`
	raw := strings.Join([]string{strconv.Itoa(len(chunk)), chunk}, "\n")
	resp, err := decodeChunkedResponse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("decodeChunkedResponse error: %v", err)
	}
	if len(resp) != 1 {
		t.Fatalf("expected 1 response, got %d", len(resp))
	}
	if resp[0].ID != "id" {
		t.Fatalf("expected id, got %q", resp[0].ID)
	}
}
