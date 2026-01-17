package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestChatSessionPathsAndPersistence(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	session := &ChatSession{
		NotebookID: "nb-123",
		Messages: []ChatMessage{
			{Role: "user", Content: "hello", Timestamp: time.Now().UTC()},
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := saveChatSession(session); err != nil {
		t.Fatalf("saveChatSession error: %v", err)
	}

	path := getChatSessionPath(session.NotebookID)
	if !strings.HasPrefix(path, filepath.Join(home, ".nlm")) {
		t.Fatalf("expected chat session under HOME/.nlm, got %s", path)
	}

	loaded, err := loadChatSession(session.NotebookID)
	if err != nil {
		t.Fatalf("loadChatSession error: %v", err)
	}
	if loaded.NotebookID != session.NotebookID {
		t.Fatalf("notebook id mismatch: %q", loaded.NotebookID)
	}
	if len(loaded.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(loaded.Messages))
	}
}

func TestListChatSessionsEmpty(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := listChatSessions(); err != nil {
		t.Fatalf("listChatSessions returned error: %v", err)
	}
}

func TestBuildContextualPrompt(t *testing.T) {
	session := &ChatSession{
		NotebookID: "nb",
		Messages: []ChatMessage{
			{Role: "user", Content: "One"},
			{Role: "assistant", Content: "Two"},
			{Role: "user", Content: "Three"},
			{Role: "assistant", Content: "Four"},
			{Role: "user", Content: "Five"},
		},
	}

	prompt := buildContextualPrompt(session, "Six")
	if !strings.Contains(prompt, "Previous conversation") {
		t.Fatalf("expected contextual prompt")
	}
	if !strings.Contains(prompt, "User: Six") {
		t.Fatalf("expected current input in prompt")
	}
}

func TestLoadChatSessionInvalidJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	path := filepath.Join(home, ".nlm")
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatalf("mkdir error: %v", err)
	}
	filePath := filepath.Join(path, "chat-bad.json")
	if err := os.WriteFile(filePath, []byte("{invalid json"), 0o600); err != nil {
		t.Fatalf("write error: %v", err)
	}

	if _, err := loadChatSession("bad"); err == nil {
		t.Fatalf("expected error for invalid json")
	}
}

func TestChatSessionJSONRoundTrip(t *testing.T) {
	msgTime := time.Now().UTC().Truncate(time.Second)
	session := ChatSession{
		NotebookID: "nb",
		Messages:   []ChatMessage{{Role: "user", Content: "hello", Timestamp: msgTime}},
		CreatedAt:  msgTime,
		UpdatedAt:  msgTime,
	}

	data, err := json.Marshal(session)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var roundTrip ChatSession
	if err := json.Unmarshal(data, &roundTrip); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if roundTrip.NotebookID != session.NotebookID {
		t.Fatalf("expected %q, got %q", session.NotebookID, roundTrip.NotebookID)
	}
}
