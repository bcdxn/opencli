package ds

import (
	_ "embed"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMemFSOpen(t *testing.T) {
	m := MemFS{
		"config/settings.json": []byte(`{"enabled":true}`),
	}

	f, err := m.Open("config/settings.json")
	if err != nil {
		t.Fatalf("open file: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })

	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(b) != `{"enabled":true}` {
		t.Fatalf("unexpected file content: %q", string(b))
	}

	info, err := f.Stat()
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if info.Name() != "config/settings.json" {
		t.Fatalf("unexpected file name: %q", info.Name())
	}
	if info.Size() != int64(len(b)) {
		t.Fatalf("unexpected file size: got %d want %d", info.Size(), len(b))
	}
	if info.Mode() != 0444 {
		t.Fatalf("unexpected file mode: %v", info.Mode())
	}
}

func TestMemFSOpenNotFound(t *testing.T) {
	m := MemFS{}

	_, err := m.Open(".")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected fs.ErrNotExist for root open, got: %v", err)
	}

	_, err = m.Open("missing.txt")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("expected fs.ErrNotExist for missing file, got: %v", err)
	}
}

func TestMemFSPersist(t *testing.T) {
	m := MemFS{
		"docs/readme.md": []byte("hello"),
		"/rooted.txt":    []byte("world"),
	}

	dir := t.TempDir()
	if err := m.Persist(dir); err != nil {
		t.Fatalf("persist failed: %v", err)
	}

	readme, err := os.ReadFile(filepath.Join(dir, "docs", "readme.md"))
	if err != nil {
		t.Fatalf("read persisted file: %v", err)
	}
	if string(readme) != "hello" {
		t.Fatalf("unexpected docs/readme.md contents: %q", string(readme))
	}

	rooted, err := os.ReadFile(filepath.Join(dir, "rooted.txt"))
	if err != nil {
		t.Fatalf("read rooted persisted file: %v", err)
	}
	if string(rooted) != "world" {
		t.Fatalf("unexpected rooted.txt contents: %q", string(rooted))
	}
}

func TestMemFSPersistBlocksPathTraversal(t *testing.T) {
	m := MemFS{
		"../escape.txt": []byte("blocked"),
	}

	err := m.Persist(t.TempDir())
	if err == nil {
		t.Fatal("expected path traversal error, got nil")
	}
	if !strings.Contains(err.Error(), "attempts to escape target directory") && !strings.Contains(err.Error(), "failed to localize") {
		t.Fatalf("unexpected error: %v", err)
	}
}
