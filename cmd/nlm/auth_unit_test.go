package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMaskProfileName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "short", input: "ab", expected: "****"},
		{name: "medium", input: "abcd", expected: "ab****"},
		{name: "long", input: "longprofile", expected: "long****file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maskProfileName(tt.input); got != tt.expected {
				t.Fatalf("maskProfileName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseAuthFlagsDefaults(t *testing.T) {
	t.Setenv("NLM_BROWSER_PROFILE", "EnvProfile")
	chromeProfile = ""

	opts, remaining, err := parseAuthFlags([]string{})
	if err != nil {
		t.Fatalf("parseAuthFlags returned error: %v", err)
	}
	if len(remaining) != 0 {
		t.Fatalf("expected no remaining args, got %v", remaining)
	}
	if opts.ProfileName != "EnvProfile" {
		t.Fatalf("expected profile EnvProfile, got %q", opts.ProfileName)
	}
}

func TestParseAuthFlagsHelp(t *testing.T) {
	_, _, err := parseAuthFlags([]string{"--help"})
	if !errors.Is(err, errHelpShown) {
		t.Fatalf("expected errHelpShown, got %v", err)
	}
}

func TestDetectAuthInfoPersistsEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	input := "curl -H 'cookie: a=1; b=2' https://example.com?at=testtoken"
	authToken, cookies, err := detectAuthInfo(input)
	if err != nil {
		t.Fatalf("detectAuthInfo returned error: %v", err)
	}
	if authToken != "testtoken" {
		t.Fatalf("expected auth token, got %q", authToken)
	}
	if cookies != "a=1; b=2" {
		t.Fatalf("expected cookies, got %q", cookies)
	}

	envPath := filepath.Join(home, ".nlm", "env")
	//nolint:gosec // test file path is controlled
	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("expected env file written: %v", err)
	}
	if !strings.Contains(string(data), "NLM_AUTH_TOKEN=\"testtoken\"") {
		t.Fatalf("env file missing auth token")
	}
}

func TestPersistAuthToDiskLoadStoredEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if _, _, err := persistAuthToDisk("cookie=1", "token", "Profile"); err != nil {
		t.Fatalf("persistAuthToDisk returned error: %v", err)
	}

	if err := os.Unsetenv("NLM_COOKIES"); err != nil {
		t.Fatalf("unset NLM_COOKIES: %v", err)
	}
	if err := os.Unsetenv("NLM_AUTH_TOKEN"); err != nil {
		t.Fatalf("unset NLM_AUTH_TOKEN: %v", err)
	}
	if err := os.Unsetenv("NLM_BROWSER_PROFILE"); err != nil {
		t.Fatalf("unset NLM_BROWSER_PROFILE: %v", err)
	}
	loadStoredEnv()

	if got := os.Getenv("NLM_COOKIES"); got != "cookie=1" {
		t.Fatalf("expected cookies loaded, got %q", got)
	}
	if got := os.Getenv("NLM_AUTH_TOKEN"); got != "token" {
		t.Fatalf("expected token loaded, got %q", got)
	}
	if got := os.Getenv("NLM_BROWSER_PROFILE"); got != "Profile" {
		t.Fatalf("expected profile loaded, got %q", got)
	}
}

func TestHandleAuthHelp(t *testing.T) {
	_, _, err := handleAuth([]string{"--help"}, false)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
