package api

import (
	"strings"
	"testing"
)

func TestDetectMIMETypeProvided(t *testing.T) {
	got := detectMIMEType([]byte("hello"), "file.bin", "text/markdown")
	if got != "text/markdown" {
		t.Fatalf("expected provided type, got %q", got)
	}
}

func TestDetectMIMETypeJSON(t *testing.T) {
	got := detectMIMEType([]byte(`{"k":"v"}`), "file.txt", "")
	if got != contentTypeJSON {
		t.Fatalf("expected json type, got %q", got)
	}
}

func TestDetectMIMETypeByContent(t *testing.T) {
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0}
	got := detectMIMEType(jpegHeader, "file.bin", "")
	if got != "image/jpeg" {
		t.Fatalf("expected image/jpeg, got %q", got)
	}
}

func TestDetectMIMETypeByExtension(t *testing.T) {
	got := detectMIMEType([]byte("plain text"), "file.csv", "")
	if !strings.HasPrefix(got, "text/csv") {
		t.Fatalf("expected text/csv prefix, got %q", got)
	}
}

func TestYouTubeHelpers(t *testing.T) {
	if !isYouTubeURL("https://youtube.com/watch?v=abc") {
		t.Fatalf("expected youtube.com to be detected")
	}
	if !isYouTubeURL("https://youtu.be/abc") {
		t.Fatalf("expected youtu.be to be detected")
	}
	if isYouTubeURL("https://example.com") {
		t.Fatalf("expected example.com to be non-youtube")
	}
}

func TestExtractYouTubeVideoID(t *testing.T) {
	id, err := extractYouTubeVideoID("https://youtu.be/xyz123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "xyz123" {
		t.Fatalf("expected xyz123, got %q", id)
	}

	id, err = extractYouTubeVideoID("https://www.youtube.com/watch?v=abc456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "abc456" {
		t.Fatalf("expected abc456, got %q", id)
	}

	if _, err := extractYouTubeVideoID("https://example.com/watch?v=1"); err == nil {
		t.Fatalf("expected error for unsupported url")
	}
}
