package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestAudioOverviewResultGetAudioBytes(t *testing.T) {
	result := &AudioOverviewResult{AudioData: ""}
	if _, err := result.GetAudioBytes(); err == nil {
		t.Fatalf("expected error for empty audio data")
	}

	encoded := base64.StdEncoding.EncodeToString([]byte("audio"))
	result.AudioData = encoded
	decoded, err := result.GetAudioBytes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(decoded) != "audio" {
		t.Fatalf("unexpected decoded data: %s", string(decoded))
	}
}

func TestAudioOverviewResultSaveAudioToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audio.wav")
	encoded := base64.StdEncoding.EncodeToString([]byte("audio-bytes"))
	result := &AudioOverviewResult{AudioData: encoded}

	if err := result.SaveAudioToFile(path); err != nil {
		t.Fatalf("save audio file error: %v", err)
	}
	//nolint:gosec // test file path is controlled
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	if string(data) != "audio-bytes" {
		t.Fatalf("unexpected file content: %s", string(data))
	}
}

func TestVideoOverviewResultSaveBase64(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "video.bin")
	encoded := base64.StdEncoding.EncodeToString([]byte("video-bytes"))
	result := &VideoOverviewResult{VideoData: encoded}

	if err := result.SaveVideoToFile(path); err != nil {
		t.Fatalf("save video file error: %v", err)
	}
	//nolint:gosec // test file path is controlled
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	if string(data) != "video-bytes" {
		t.Fatalf("unexpected file content: %s", string(data))
	}
}

func TestVideoOverviewResultSaveFromURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("video-data"))
	}))
	defer server.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "video.bin")
	result := &VideoOverviewResult{VideoData: server.URL}

	if err := result.SaveVideoToFile(path); err != nil {
		t.Fatalf("save video file error: %v", err)
	}
	//nolint:gosec // test file path is controlled
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file error: %v", err)
	}
	if string(data) != "video-data" {
		t.Fatalf("unexpected file content: %s", string(data))
	}
}

func TestFindVideoURL(t *testing.T) {
	client := &Client{}
	url := "https://lh3.googleusercontent.com/rd-notebooklm/test"
	data := []interface{}{
		[]interface{}{
			"irrelevant",
			[]interface{}{url},
		},
	}

	if got := client.extractVideoURLFromResponse(data); got != url {
		t.Fatalf("expected %q, got %q", url, got)
	}
}
