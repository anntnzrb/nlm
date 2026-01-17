//go:build darwin

package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetProfilePathPrefersChrome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	chromePath := filepath.Join(home, "Library", "Application Support", "Google", "Chrome")
	if err := os.MkdirAll(chromePath, 0o700); err != nil {
		t.Fatalf("mkdir chrome: %v", err)
	}

	got := getProfilePath()
	if got != chromePath {
		t.Fatalf("expected chrome path %q, got %q", chromePath, got)
	}
}

func TestGetProfilePathFallsBackToBrave(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	bravePath := filepath.Join(home, "Library", "Application Support", "BraveSoftware", "Brave-Browser")
	if err := os.MkdirAll(bravePath, 0o700); err != nil {
		t.Fatalf("mkdir brave: %v", err)
	}

	got := getProfilePath()
	if got != bravePath {
		t.Fatalf("expected brave path %q, got %q", bravePath, got)
	}
}

func TestProfilePathHelpersUseHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	expectedCanary := filepath.Join(home, "Library", "Application Support", "Google", "Chrome Canary")
	if got := getCanaryProfilePath(); got != expectedCanary {
		t.Fatalf("expected canary path %q, got %q", expectedCanary, got)
	}

	expectedBrave := filepath.Join(home, "Library", "Application Support", "BraveSoftware", "Brave-Browser")
	if got := getBraveProfilePath(); got != expectedBrave {
		t.Fatalf("expected brave path %q, got %q", expectedBrave, got)
	}
}

func TestGetMostRecentPath(t *testing.T) {
	root := t.TempDir()

	older := filepath.Join(root, "old")
	if err := os.MkdirAll(older, 0o700); err != nil {
		t.Fatalf("mkdir old: %v", err)
	}
	if err := os.Chtimes(older, time.Now().Add(-2*time.Hour), time.Now().Add(-2*time.Hour)); err != nil {
		t.Fatalf("chtimes old: %v", err)
	}

	newer := filepath.Join(root, "new")
	if err := os.MkdirAll(newer, 0o700); err != nil {
		t.Fatalf("mkdir new: %v", err)
	}
	if err := os.Chtimes(newer, time.Now().Add(-1*time.Hour), time.Now().Add(-1*time.Hour)); err != nil {
		t.Fatalf("chtimes new: %v", err)
	}

	got := getMostRecentPath([]string{older, newer})
	if got != newer {
		t.Fatalf("expected newest path %q, got %q", newer, got)
	}
}
