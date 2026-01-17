package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tmc/nlm/internal/batchexecute"
)

func TestIsValidCommand(t *testing.T) {
	valid := []string{"list", "ls", "auth", cmdChatList, cmdRefresh}
	for _, cmd := range valid {
		if !isValidCommand(cmd) {
			t.Fatalf("expected command %q to be valid", cmd)
		}
	}

	if isValidCommand("not-a-command") {
		t.Fatalf("expected invalid command to be false")
	}
}

func TestIsAuthCommand(t *testing.T) {
	noAuth := []string{"help", "-h", "--help", "auth", cmdRefresh, cmdChatList}
	for _, cmd := range noAuth {
		if isAuthCommand(cmd) {
			t.Fatalf("expected command %q to not require auth", cmd)
		}
	}

	if !isAuthCommand("list") {
		t.Fatalf("expected list to require auth")
	}
}

func TestValidateArgs(t *testing.T) {
	if err := validateArgs("create", []string{"title"}); err != nil {
		t.Fatalf("expected create to be valid: %v", err)
	}
	if err := validateArgs("create", []string{}); err == nil {
		t.Fatalf("expected create without args to fail")
	}

	if err := validateArgs("rm", []string{"id"}); err != nil {
		t.Fatalf("expected rm to be valid: %v", err)
	}
	if err := validateArgs("rm", []string{}); err == nil {
		t.Fatalf("expected rm without args to fail")
	}

	if err := validateArgs("notes", []string{"notebook"}); err != nil {
		t.Fatalf("expected notes to be valid: %v", err)
	}
	if err := validateArgs("notes", []string{}); err == nil {
		t.Fatalf("expected notes without args to fail")
	}
}

func TestIsAuthenticationError(t *testing.T) {
	if !isAuthenticationError(batchexecute.ErrUnauthorized) {
		t.Fatalf("expected unauthorized to be auth error")
	}

	errs := []string{
		"unauthenticated",
		"authentication failed",
		"unauthorized",
		"api error 16",
		"error 16",
		"status: 401",
		"status: 403",
		"session.*invalid",
		"session.*expired",
		"login.*required",
		"auth.*required",
		"invalid.*credentials",
		"token.*expired",
		"cookie.*invalid",
	}
	for _, msg := range errs {
		if !isAuthenticationError(fmt.Errorf("%s", msg)) {
			t.Fatalf("expected %q to be auth error", msg)
		}
	}
}

func TestSaveCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	chromeProfile = "Profile"
	if err := saveCredentials("token", "cookies"); err != nil {
		t.Fatalf("saveCredentials error: %v", err)
	}

	path := filepath.Join(home, ".nlm", "env")
	//nolint:gosec // test file path is controlled
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read env file: %v", err)
	}
	if !strings.Contains(string(data), "NLM_AUTH_TOKEN") {
		t.Fatalf("expected token in env file")
	}
}

func TestGetChatSessionPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := getChatSessionPath("nb-xyz")
	if !strings.Contains(path, filepath.Join(home, ".nlm")) {
		t.Fatalf("expected path in HOME/.nlm, got %s", path)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Fatalf("expected .nlm directory to exist: %v", err)
	}
}
