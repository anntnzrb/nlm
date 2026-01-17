package main

import "testing"

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
