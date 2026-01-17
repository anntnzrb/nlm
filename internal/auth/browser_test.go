package auth

import "testing"

func TestBrowserString(t *testing.T) {
	b := Browser{
		Name:    "TestBrowser",
		Version: "1.2.3",
	}
	if got := b.String(); got != "TestBrowser (1.2.3)" {
		t.Fatalf("unexpected string: %q", got)
	}
}
