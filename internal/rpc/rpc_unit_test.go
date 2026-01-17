package rpc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAPIParamsEnvOverride(t *testing.T) {
	paramsMutex.Lock()
	cachedParams = nil
	paramsMutex.Unlock()

	t.Setenv("NLM_BUILD_VERSION", "boq_labs-tailwind-frontend_test")
	t.Setenv("NLM_SESSION_ID", "12345")

	params := GetAPIParams("")
	if params.BuildVersion != "boq_labs-tailwind-frontend_test" {
		t.Fatalf("expected env build version, got %q", params.BuildVersion)
	}
	if params.SessionID != "12345" {
		t.Fatalf("expected env session id, got %q", params.SessionID)
	}
}

func TestFetchAPIParamsFromPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"cfb2h":"boq_labs-tailwind-frontend_x","FdrFJe":"-999"}`))
	}))
	defer server.Close()

	orig := notebookLMURL
	defer func() { notebookLMURL = orig }()
	notebookLMURL = server.URL

	params := fetchAPIParamsFromPage("cookie=1")
	if params == nil {
		t.Fatalf("expected params")
	}
	if params.BuildVersion != "boq_labs-tailwind-frontend_x" {
		t.Fatalf("build version mismatch: %q", params.BuildVersion)
	}
	if params.SessionID != "-999" {
		t.Fatalf("session id mismatch: %q", params.SessionID)
	}
}

func TestFetchAPIParamsFromPageAltPatterns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`bl='boq_labs-tailwind-frontend_alt' f.sid='42'`))
	}))
	defer server.Close()

	orig := notebookLMURL
	defer func() { notebookLMURL = orig }()
	notebookLMURL = server.URL

	params := fetchAPIParamsFromPage("cookie=1")
	if params == nil {
		t.Fatalf("expected params")
	}
	if params.BuildVersion != "boq_labs-tailwind-frontend_alt" {
		t.Fatalf("build version mismatch: %q", params.BuildVersion)
	}
	if params.SessionID != "42" {
		t.Fatalf("session id mismatch: %q", params.SessionID)
	}
}
