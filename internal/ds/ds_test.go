package ds

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
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

func TestMapSetGetAndOrder(t *testing.T) {
	om := NewMap[string, int]()
	om.Set("first", 1)
	om.Set("second", 2)
	om.Set("first", 3) // update existing key should not change insertion order

	v, ok := om.Get("first")
	if !ok || v != 3 {
		t.Fatalf("unexpected get(first): v=%d ok=%v", v, ok)
	}

	if _, ok := om.Get("missing"); ok {
		t.Fatal("expected missing key lookup to return ok=false")
	}

	entries := om.Entries()
	if len(entries) != 2 {
		t.Fatalf("unexpected entry count: %d", len(entries))
	}
	if entries[0].Key != "first" || entries[0].Value != 3 {
		t.Fatalf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Key != "second" || entries[1].Value != 2 {
		t.Fatalf("unexpected second entry: %+v", entries[1])
	}

	builtin := om.ToBuiltInMap()
	if len(builtin) != 2 || builtin["first"] != 3 || builtin["second"] != 2 {
		t.Fatalf("unexpected built-in map conversion: %+v", builtin)
	}
}

func TestMapUnmarshalJSONPreservesOrder(t *testing.T) {
	var om Map[string, int]
	data := []byte(`{"a":1,"b":2,"a":3,"c":4}`)

	if err := om.UnmarshalJSON(data); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}

	entries := om.Entries()
	if len(entries) != 3 {
		t.Fatalf("unexpected entry count: %d", len(entries))
	}
	if entries[0].Key != "a" || entries[0].Value != 3 {
		t.Fatalf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Key != "b" || entries[1].Value != 2 {
		t.Fatalf("unexpected second entry: %+v", entries[1])
	}
	if entries[2].Key != "c" || entries[2].Value != 4 {
		t.Fatalf("unexpected third entry: %+v", entries[2])
	}
}

func TestMapUnmarshalJSONRejectsNonObject(t *testing.T) {
	var om Map[string, int]
	err := om.UnmarshalJSON([]byte(`[1,2,3]`))
	if err == nil {
		t.Fatal("expected non-object JSON error, got nil")
	}
}

func TestMapJSONRoundTrip(t *testing.T) {
	om := NewMap[string, int]()
	om.Set("one", 1)
	om.Set("two", 2)

	b, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal json: %v", err)
	}

	var decoded Map[string, int]
	if err := decoded.UnmarshalJSON(b); err != nil {
		t.Fatalf("unmarshal marshaled json: %v", err)
	}

	entries := decoded.Entries()
	if len(entries) != 2 || entries[0].Key != "one" || entries[1].Key != "two" {
		t.Fatalf("unexpected round-trip order: %+v", entries)
	}
}

func TestMapYAMLRoundTrip(t *testing.T) {
	input := "first: 1\nsecond: 2\nthird: 3\n"

	var om Map[string, int]
	if err := yaml.Unmarshal([]byte(input), &om); err != nil {
		t.Fatalf("unmarshal yaml: %v", err)
	}

	entries := om.Entries()
	if len(entries) != 3 || entries[0].Key != "first" || entries[1].Key != "second" || entries[2].Key != "third" {
		t.Fatalf("unexpected yaml unmarshal order: %+v", entries)
	}

	b, err := om.MarshalYAML()
	if err != nil {
		t.Fatalf("marshal yaml: %v", err)
	}

	var decoded Map[string, int]
	if err := yaml.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal marshaled yaml: %v", err)
	}

	roundTrip := decoded.Entries()
	if len(roundTrip) != 3 || roundTrip[0].Key != "first" || roundTrip[1].Key != "second" || roundTrip[2].Key != "third" {
		t.Fatalf("unexpected yaml round-trip order: %+v", roundTrip)
	}
}

func TestMapUnmarshalYAMLRejectsNonMapping(t *testing.T) {
	var om Map[string, int]
	err := yaml.Unmarshal([]byte("- 1\n- 2\n"), &om)
	if err == nil {
		t.Fatal("expected non-mapping YAML error, got nil")
	}
}
