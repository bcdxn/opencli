package ds

import (
	_ "embed"
	"slices"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
)

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

	b, err := yaml.Marshal(&om)
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

func TestMemFileInfoMethods(t *testing.T) {
	m := MemFS{"test.txt": []byte("hello")}
	f, err := m.Open("test.txt")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = f.Close() })

	info, err := f.Stat()
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if time.Since(info.ModTime()) > time.Minute {
		t.Fatal("expected ModTime to be within last minute")
	}
	if info.IsDir() {
		t.Fatal("expected IsDir to be false for file")
	}
	if info.Sys() != nil {
		t.Fatal("expected Sys to return nil")
	}
}

func TestMapKeys(t *testing.T) {
	om := NewMap[string, int]()
	om.Set("a", 1)
	om.Set("b", 2)
	om.Set("c", 3)

	keys := om.Keys()
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	keySet := make(map[int]bool)
	for _, k := range keys {
		switch k {
		case "a", "b", "c":
			keySet[0] = true
		default:
			t.Fatalf("unexpected key: %q", k)
		}
	}
	if !keySet[0] {
		t.Fatal("expected all keys a, b, c to be present")
	}
}

//go:embed testdata/nested-indent-input.yaml
var indentTestingYaml []byte

//go:embed testdata/expected.yaml
var expectedIndentTestinYaml []byte

func TestYAMLMarshallingIndentLevel(t *testing.T) {
	var ts testStruct
	err := yaml.Unmarshal(indentTestingYaml, &ts)
	if err != nil {
		t.Fatal("error unmarshalling test struct", err)
	}

	// now marshall and ensure indent levels match expected levels
	actual, err := yaml.MarshalWithOptions(ts, yaml.UseLiteralStyleIfMultiline(true))
	if err != nil {
		t.Fatal("error marshalling test struct", actual)
	}

	if !slices.Equal(expectedIndentTestinYaml, actual) {
		t.Fatal("slices not equal!")
	}
}

type testStruct struct {
	OpencliVersion string                `yaml:"opencliVersion"`
	Commands       *Map[string, testCmd] `yaml:"commands"`
}

type testCmd struct {
	Description string `yaml:"description"`
}
