package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCheckProfileForDomainCookies(t *testing.T) {
	dir := t.TempDir()
	cookiesPath := filepath.Join(dir, "Cookies")

	// Small file should fail.
	if err := os.WriteFile(cookiesPath, []byte("tiny"), 0o600); err != nil {
		t.Fatalf("write cookies: %v", err)
	}
	if checkProfileForDomainCookies(cookiesPath, "notebooklm.google.com") {
		t.Fatalf("expected small cookies file to be false")
	}

	// Large file but old mtime should fail.
	data := make([]byte, 2048)
	if err := os.WriteFile(cookiesPath, data, 0o600); err != nil {
		t.Fatalf("write cookies: %v", err)
	}
	old := time.Now().Add(-45 * 24 * time.Hour)
	if err := os.Chtimes(cookiesPath, old, old); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	if checkProfileForDomainCookies(cookiesPath, "notebooklm.google.com") {
		t.Fatalf("expected old cookies file to be false")
	}

	// Large file with recent mtime should pass.
	recent := time.Now().Add(-2 * 24 * time.Hour)
	if err := os.Chtimes(cookiesPath, recent, recent); err != nil {
		t.Fatalf("chtimes: %v", err)
	}
	if !checkProfileForDomainCookies(cookiesPath, "notebooklm.google.com") {
		t.Fatalf("expected recent cookies file to be true")
	}
}

func TestScanBrowserProfiles(t *testing.T) {
	root := t.TempDir()

	makeProfile := func(name string, withCookies bool) {
		profileDir := filepath.Join(root, name)
		if err := os.MkdirAll(profileDir, 0o700); err != nil {
			t.Fatalf("mkdir profile: %v", err)
		}
		// Required files for valid profile.
		if err := os.WriteFile(filepath.Join(profileDir, "Cookies"), make([]byte, 2048), 0o600); err != nil {
			t.Fatalf("write cookies: %v", err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "Login Data"), []byte("login"), 0o600); err != nil {
			t.Fatalf("write login: %v", err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "History"), []byte("history"), 0o600); err != nil {
			t.Fatalf("write history: %v", err)
		}

		if !withCookies {
			old := time.Now().Add(-60 * 24 * time.Hour)
			if err := os.Chtimes(filepath.Join(profileDir, "Cookies"), old, old); err != nil {
				t.Fatalf("chtimes: %v", err)
			}
		}
	}

	makeProfile("Default", true)
	makeProfile("Profile 1", false)
	if err := os.MkdirAll(filepath.Join(root, "System Profile"), 0o700); err != nil {
		t.Fatalf("mkdir system profile: %v", err)
	}

	profiles, err := scanBrowserProfiles(root, "Chrome", "notebooklm.google.com")
	if err != nil {
		t.Fatalf("scanBrowserProfiles error: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}

	for _, profile := range profiles {
		if profile.Name == "Default" && !profile.HasTargetCookies {
			t.Fatalf("expected Default to have cookies")
		}
		if profile.Name == "Profile 1" && profile.HasTargetCookies {
			t.Fatalf("expected Profile 1 to lack cookies")
		}
	}
}
