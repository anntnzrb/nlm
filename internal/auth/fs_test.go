package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindMostRecentProfile(t *testing.T) {
	root := t.TempDir()

	makeProfile := func(name string, mod time.Time) {
		profileDir := filepath.Join(root, name)
		if err := os.MkdirAll(profileDir, 0o700); err != nil {
			t.Fatalf("mkdir profile: %v", err)
		}
		if err := os.WriteFile(filepath.Join(profileDir, "Cookies"), make([]byte, 10), 0o600); err != nil {
			t.Fatalf("write cookies: %v", err)
		}
		if err := os.Chtimes(profileDir, mod, mod); err != nil {
			t.Fatalf("chtimes: %v", err)
		}
	}

	makeProfile("Old", time.Now().Add(-48*time.Hour))
	makeProfile("New", time.Now().Add(-1*time.Hour))

	if err := os.MkdirAll(filepath.Join(root, "System Profile"), 0o700); err != nil {
		t.Fatalf("mkdir system: %v", err)
	}

	got := findMostRecentProfile(root)
	if got == "" || filepath.Base(got) != "New" {
		t.Fatalf("expected newest profile, got %q", got)
	}
}

func TestCopyDirectoryRecursiveWithCount(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	if err := os.MkdirAll(filepath.Join(src, "subdir"), 0o700); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "file1.txt"), []byte("one"), 0o600); err != nil {
		t.Fatalf("write file1: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "subdir", "file2.txt"), []byte("two"), 0o600); err != nil {
		t.Fatalf("write file2: %v", err)
	}

	var files, dirs int
	if err := copyDirectoryRecursiveWithCount(src, dst, false, &files, &dirs); err != nil {
		t.Fatalf("copyDirectoryRecursiveWithCount error: %v", err)
	}

	if files != 2 {
		t.Fatalf("expected 2 files copied, got %d", files)
	}
	if dirs != 1 {
		t.Fatalf("expected 1 dir copied, got %d", dirs)
	}
}
