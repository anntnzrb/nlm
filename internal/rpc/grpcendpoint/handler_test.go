package grpcendpoint

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/tmc/nlm/internal/rpc"
)

type rewriteTransport struct {
	base *url.URL
}

func (t rewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = t.base.Scheme
	req.URL.Host = t.base.Host
	return http.DefaultTransport.RoundTrip(req)
}

func TestBuildChatRequest(t *testing.T) {
	req := BuildChatRequest([]string{"src-a", "src-b"}, "hello")

	parts, ok := req.([]interface{})
	if !ok || len(parts) != 2 {
		t.Fatalf("expected two-part request, got %#v", req)
	}
	if parts[0] != nil {
		t.Fatalf("expected nil leading element, got %#v", parts[0])
	}

	innerJSON, ok := parts[1].(string)
	if !ok {
		t.Fatalf("expected inner JSON string, got %#v", parts[1])
	}

	var inner []interface{}
	if err := json.Unmarshal([]byte(innerJSON), &inner); err != nil {
		t.Fatalf("unmarshal inner JSON: %v", err)
	}
	if len(inner) != 4 {
		t.Fatalf("expected 4 inner elements, got %#v", inner)
	}

	expectedSources := []interface{}{
		[]interface{}{
			[]interface{}{"src-a", "src-b"},
		},
	}
	if !reflect.DeepEqual(inner[0], expectedSources) {
		t.Fatalf("unexpected sources: %#v", inner[0])
	}
	if inner[1] != "hello" {
		t.Fatalf("unexpected prompt: %#v", inner[1])
	}
	if inner[2] != nil {
		t.Fatalf("expected nil slot, got %#v", inner[2])
	}

	expectedMeta := []interface{}{float64(2), nil, []interface{}{float64(1)}}
	if !reflect.DeepEqual(inner[3], expectedMeta) {
		t.Fatalf("unexpected metadata: %#v", inner[3])
	}
}

func TestGenerateRequestID(t *testing.T) {
	original := requestCounter
	t.Cleanup(func() {
		requestCounter = original
	})

	requestCounter = 0
	if got := generateRequestID(); got != 1000001 {
		t.Fatalf("expected first ID 1000001, got %d", got)
	}
	if got := generateRequestID(); got != 1000002 {
		t.Fatalf("expected second ID 1000002, got %d", got)
	}
}

func TestClientExecuteSuccess(t *testing.T) {
	rpc.ClearAPIParamsCache()
	t.Setenv("NLM_BUILD_VERSION", "bl-test")
	t.Setenv("NLM_SESSION_ID", "sid-test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		if !strings.HasPrefix(r.URL.Path, "/_/LabsTailwindUi/data") {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("bl") != "bl-test" || query.Get("f.sid") != "sid-test" {
			t.Fatalf("missing API params: %s", r.URL.RawQuery)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		bodyStr := string(body)
		if !strings.Contains(bodyStr, "f.req=") || !strings.Contains(bodyStr, "at=token") {
			t.Fatalf("unexpected request body: %s", bodyStr)
		}

		response := ")]}'\n1\n[[\"wrb.fr\",null,\"{\\\"ok\\\":true}\"]]\n"
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}

	client := NewClient("token", "cookie=1")
	client.httpClient = &http.Client{
		Transport: rewriteTransport{base: base},
	}

	resp, err := client.Execute(Request{
		Endpoint: "/rpc",
		Body:     map[string]string{"hello": "world"},
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if string(resp) != `{"ok":true}` {
		t.Fatalf("unexpected response: %s", string(resp))
	}
}

func TestClientExecuteStatusError(t *testing.T) {
	rpc.ClearAPIParamsCache()
	t.Setenv("NLM_BUILD_VERSION", "bl-test")
	t.Setenv("NLM_SESSION_ID", "sid-test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("nope"))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}

	client := NewClient("token", "cookie=1")
	client.httpClient = &http.Client{
		Transport: rewriteTransport{base: base},
	}

	_, err = client.Execute(Request{
		Endpoint: "/rpc",
		Body:     map[string]string{"hello": "world"},
	})
	if err == nil || !strings.Contains(err.Error(), "status 403") {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestClientStream(t *testing.T) {
	rpc.ClearAPIParamsCache()
	t.Setenv("NLM_BUILD_VERSION", "bl-test")
	t.Setenv("NLM_SESSION_ID", "sid-test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}

	client := NewClient("token", "cookie=1")
	client.httpClient = &http.Client{
		Transport: rewriteTransport{base: base},
	}

	var collected []byte
	err = client.Stream(Request{
		Endpoint: "/stream",
		Body:     map[string]string{"hello": "world"},
	}, func(chunk []byte) error {
		collected = append(collected, chunk...)
		return nil
	})
	if err != nil {
		t.Fatalf("stream: %v", err)
	}
	if string(collected) != "hello" {
		t.Fatalf("unexpected streamed data: %s", string(collected))
	}
}

func TestClientStreamStatusError(t *testing.T) {
	rpc.ClearAPIParamsCache()
	t.Setenv("NLM_BUILD_VERSION", "bl-test")
	t.Setenv("NLM_SESSION_ID", "sid-test")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("bad"))
	}))
	defer server.Close()

	base, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("parse server URL: %v", err)
	}

	client := NewClient("token", "cookie=1")
	client.httpClient = &http.Client{
		Transport: rewriteTransport{base: base},
	}

	err = client.Stream(Request{
		Endpoint: "/stream",
		Body:     map[string]string{"hello": "world"},
	}, func(chunk []byte) error {
		return nil
	})
	if err == nil || !strings.Contains(err.Error(), "status 400") {
		t.Fatalf("expected status error, got %v", err)
	}
}
