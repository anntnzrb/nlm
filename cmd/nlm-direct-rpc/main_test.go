package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainMissingEnv(t *testing.T) {
	//nolint:gosec // test helper uses the current test binary
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcessMissingEnv")
	cmd.Env = append(filteredEnv(os.Environ(), "NLM_AUTH_TOKEN", "NLM_COOKIES"), "GO_WANT_HELPER_PROCESS=1")

	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit, got none")
	}
	if !strings.Contains(string(output), "Please source ~/.nlm/env first") {
		t.Fatalf("unexpected output: %s", string(output))
	}
}

func TestHelperProcessMissingEnv(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	main()
}

func filteredEnv(env []string, keys ...string) []string {
	deny := map[string]struct{}{}
	for _, key := range keys {
		deny[key+"="] = struct{}{}
	}

	filtered := make([]string, 0, len(env))
	for _, entry := range env {
		skip := false
		for prefix := range deny {
			if strings.HasPrefix(entry, prefix) {
				skip = true
				break
			}
		}
		if !skip {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}
