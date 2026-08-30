package codec_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
)

// defaultTestSpecYAML builds a minimal single-command spec with one flag. The
// fragment is the body of that flag entry (name, type, variadic, default...).
func defaultTestSpecYAML(t *testing.T, frag string) []byte {
	t.Helper()

	doc := strings.Join([]string{
		"opencliVersion: 1.0.0-alpha.14",
		"info:",
		"  title: default test",
		"  version: 0.1.0",
		"  binary: deftest",
		"commands:",
		"  root:",
		"    flags:",
		// the fragment's continuation lines are pre-indented to align under "name:"
		"      - " + frag,
	}, "\n")

	return []byte(doc)
}

// defaultTestSpecJSON builds the JSON equivalent of a single-flag spec.
func defaultTestSpecJSON(t *testing.T, flagObj string) []byte {
	t.Helper()

	doc := `{"opencliVersion":"1.0.0-alpha.14","info":{"title":"default test","version":"0.1.0","binary":"deftest"},"commands":{"root":{"flags":[` + flagObj + `]}}}`

	return []byte(doc)
}

// globalFlagSpecYAML builds a minimal spec with one flag in the global section.
func globalFlagSpecYAML(t *testing.T, frag string) []byte {
	t.Helper()

	doc := strings.Join([]string{
		"opencliVersion: 1.0.0-alpha.14",
		"info:",
		"  title: default test",
		"  version: 0.1.0",
		"  binary: deftest",
		"global:",
		"  flags:",
		// the fragment's continuation lines are pre-indented to align under "name:"
		"    - " + frag,
		"commands:",
		"  root: {}",
	}, "\n")

	return []byte(doc)
}

// globalFlagSpecJSON builds a minimal spec with one flag in the global section.
func globalFlagSpecJSON(t *testing.T, flagObj string) []byte {
	t.Helper()

	doc := `{"opencliVersion":"1.0.0-alpha.14","info":{"title":"default test","version":"0.1.0","binary":"deftest"},"global":{"flags":[` + flagObj + `]},"commands":{"root":{}}}`

	return []byte(doc)
}

// rootFlagDefault returns the Default value of the named flag on the root command.
func rootFlagDefault(t *testing.T, d *spec.Document, name string) any {
	t.Helper()

	if d.Commands == nil || len(d.Commands.Flags) == 0 {
		t.Fatal("expected a root command with flags")
	}

	for _, f := range d.Commands.Flags {
		if f.Name == name {
			return f.Default
		}
	}

	t.Fatalf("flag %q not found on root command", name)

	return nil
}

// globalFlagDefault returns the Default value of the named flag in the global section.
func globalFlagDefault(t *testing.T, d *spec.Document, name string) any {
	t.Helper()

	if d.Global == nil || len(d.Global.Flags) == 0 {
		t.Fatal("expected a global section with flags")
	}

	for _, f := range d.Global.Flags {
		if f.Name == name {
			return f.Default
		}
	}

	t.Fatalf("flag %q not found in the global section", name)

	return nil
}

// TestUnmarshalFlagDefaultNormalization verifies that flag default values are
// coerced to canonical Go types regardless of the source format. goccy/go-yaml
// decodes integers as uint64 while encoding/json uses float64; both must end up
// as int64 for integer flags so downstream emitters see one stable shape.
func TestUnmarshalFlagDefaultNormalization(t *testing.T) {
	cases := []struct {
		name     string
		yamlFrag string // flag fragment, e.g.: `name: n\n        type: number`
		jsonObj  string // JSON object for the same flag
		wantVal  any    // expected normalized value (nil means "no default")
	}{
		{
			name:     "integer from yaml is int64",
			yamlFrag: "name: n\n        type: integer\n        default: 42",
			jsonObj:  `{"name":"n","type":"integer","default":42}`,
			wantVal:  int64(42),
		},
		{
			name:     "number from yaml is float64",
			yamlFrag: "name: n\n        type: number\n        default: 3.5",
			jsonObj:  `{"name":"n","type":"number","default":3.5}`,
			wantVal:  3.5,
		},
		{
			name:     "string stays string",
			yamlFrag: "name: n\n        type: string\n        default: hello",
			jsonObj:  `{"name":"n","type":"string","default":"hello"}`,
			wantVal:  "hello",
		},
		{
			name:     "boolean stays bool",
			yamlFrag: "name: n\n        type: boolean\n        default: true",
			jsonObj:  `{"name":"n","type":"boolean","default":true}`,
			wantVal:  true,
		},
		{
			name:     "no default stays nil",
			yamlFrag: "name: n\n        type: integer",
			jsonObj:  `{"name":"n","type":"integer"}`,
			wantVal:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+" (yaml)", func(t *testing.T) {
			d, err := codec.UnmarshalYAML(defaultTestSpecYAML(t, tc.yamlFrag))
			if err != nil {
				t.Fatalf("unmarshal yaml failed: %v", err)
			}

			got := rootFlagDefault(t, d, "n")
			if !reflect.DeepEqual(got, tc.wantVal) {
				t.Errorf("yaml default = %#v (%T), want %#v (%T)", got, got, tc.wantVal, tc.wantVal)
			}
		})

		t.Run(tc.name+" (json)", func(t *testing.T) {
			d, err := codec.UnmarshalJSON(defaultTestSpecJSON(t, tc.jsonObj))
			if err != nil {
				t.Fatalf("unmarshal json failed: %v", err)
			}

			got := rootFlagDefault(t, d, "n")
			if !reflect.DeepEqual(got, tc.wantVal) {
				t.Errorf("json default = %#v (%T), want %#v (%T)", got, got, tc.wantVal, tc.wantVal)
			}
		})
	}
}

// TestUnmarshalVariadicFlagDefaultNormalization covers list defaults on
// variadic flags. The schema currently rejects these (default is oneOf scalars
// only), but the codec decodes them anyway and must produce canonical typed
// slices rather than decoder-specific []interface{} shapes that panic in the
// emitters.
func TestUnmarshalVariadicFlagDefaultNormalization(t *testing.T) {
	cases := []struct {
		name     string
		yamlFrag string
		jsonObj  string
		wantVal  any
	}{
		{
			name:     "integer list",
			yamlFrag: "name: n\n        type: integer\n        variadic: true\n        default:\n          - 1\n          - 2\n          - 3",
			jsonObj:  `{"name":"n","type":"integer","variadic":true,"default":[1,2,3]}`,
			wantVal:  []int64{1, 2, 3},
		},
		{
			name:     "number list",
			yamlFrag: "name: n\n        type: number\n        variadic: true\n        default:\n          - 1.5\n          - 2.5",
			jsonObj:  `{"name":"n","type":"number","variadic":true,"default":[1.5,2.5]}`,
			wantVal:  []float64{1.5, 2.5},
		},
		{
			name:     "string list",
			yamlFrag: "name: n\n        type: string\n        variadic: true\n        default:\n          - a\n          - b",
			jsonObj:  `{"name":"n","type":"string","variadic":true,"default":["a","b"]}`,
			wantVal:  []string{"a", "b"},
		},
		{
			name:     "boolean list",
			yamlFrag: "name: n\n        type: boolean\n        variadic: true\n        default:\n          - true\n          - false",
			jsonObj:  `{"name":"n","type":"boolean","variadic":true,"default":[true,false]}`,
			wantVal:  []bool{true, false},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+" (yaml)", func(t *testing.T) {
			d, err := codec.UnmarshalYAML(defaultTestSpecYAML(t, tc.yamlFrag))
			if err != nil {
				t.Fatalf("unmarshal yaml failed: %v", err)
			}

			got := rootFlagDefault(t, d, "n")
			if !reflect.DeepEqual(got, tc.wantVal) {
				t.Errorf("yaml default = %#v (%T), want %#v (%T)", got, got, tc.wantVal, tc.wantVal)
			}
		})

		t.Run(tc.name+" (json)", func(t *testing.T) {
			d, err := codec.UnmarshalJSON(defaultTestSpecJSON(t, tc.jsonObj))
			if err != nil {
				t.Fatalf("unmarshal json failed: %v", err)
			}

			got := rootFlagDefault(t, d, "n")
			if !reflect.DeepEqual(got, tc.wantVal) {
				t.Errorf("json default = %#v (%T), want %#v (%T)", got, got, tc.wantVal, tc.wantVal)
			}
		})
	}
}

// TestUnmarshalFlagDefaultTypeMismatch verifies that a default value which
// cannot be represented as the flag's declared type is reported at decode time
// instead of being silently dropped or mis-emitted downstream.
func TestUnmarshalFlagDefaultTypeMismatch(t *testing.T) {
	cases := []struct {
		name     string
		yamlFrag string
		jsonObj  string
	}{
		{
			name:     "string default on integer flag",
			yamlFrag: "name: n\n        type: integer\n        default: not-a-number",
			jsonObj:  `{"name":"n","type":"integer","default":"not-a-number"}`,
		},
		{
			name:     "string default on boolean flag",
			yamlFrag: "name: n\n        type: boolean\n        default: maybe",
			jsonObj:  `{"name":"n","type":"boolean","default":"maybe"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+" (yaml)", func(t *testing.T) {
			if _, err := codec.UnmarshalYAML(defaultTestSpecYAML(t, tc.yamlFrag)); err == nil {
				t.Error("expected an error for a type-mismatched flag default, got nil")
			} else if !strings.Contains(err.Error(), `flag "n"`) {
				t.Errorf("error should name the offending flag: %v", err)
			}
		})

		t.Run(tc.name+" (json)", func(t *testing.T) {
			if _, err := codec.UnmarshalJSON(defaultTestSpecJSON(t, tc.jsonObj)); err == nil {
				t.Error("expected an error for a type-mismatched flag default, got nil")
			} else if !strings.Contains(err.Error(), `flag "n"`) {
				t.Errorf("error should name the offending flag: %v", err)
			}
		})
	}
}

// TestUnmarshalGlobalFlagDefaultNormalization verifies that global flag default
// values are coerced to canonical Go types just like command-level flags. Global
// flags live outside the command Trie, so they must be normalized explicitly at
// decode time; otherwise decoder-native types (uint64 from YAML, float64 from
// JSON) leak into code generation and produce wrong or uncompilable defaults.
func TestUnmarshalGlobalFlagDefaultNormalization(t *testing.T) {
	cases := []struct {
		name     string
		yamlFrag string // flag fragment, e.g.: `name: n\n      type: number`
		jsonObj  string // JSON object for the same flag
		wantVal  any    // expected normalized value (nil means "no default")
	}{
		{
			name:     "integer from yaml is int64",
			yamlFrag: "name: n\n      type: integer\n      default: 30",
			jsonObj:  `{"name":"n","type":"integer","default":30}`,
			wantVal:  int64(30),
		},
		{
			name:     "number from yaml is float64",
			yamlFrag: "name: n\n      type: number\n      default: 2.5",
			jsonObj:  `{"name":"n","type":"number","default":2.5}`,
			wantVal:  2.5,
		},
		{
			name:     "string stays string",
			yamlFrag: "name: n\n      type: string\n      default: hello",
			jsonObj:  `{"name":"n","type":"string","default":"hello"}`,
			wantVal:  "hello",
		},
		{
			name:     "boolean stays bool",
			yamlFrag: "name: n\n      type: boolean\n      default: true",
			jsonObj:  `{"name":"n","type":"boolean","default":true}`,
			wantVal:  true,
		},
		{
			name:     "no default stays nil",
			yamlFrag: "name: n\n      type: integer",
			jsonObj:  `{"name":"n","type":"integer"}`,
			wantVal:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+" (yaml)", func(t *testing.T) {
			d, err := codec.UnmarshalYAML(globalFlagSpecYAML(t, tc.yamlFrag))
			if err != nil {
				t.Fatalf("unmarshal yaml failed: %v", err)
			}

			got := globalFlagDefault(t, d, "n")
			if !reflect.DeepEqual(got, tc.wantVal) {
				t.Errorf("yaml default = %#v (%T), want %#v (%T)", got, got, tc.wantVal, tc.wantVal)
			}
		})

		t.Run(tc.name+" (json)", func(t *testing.T) {
			d, err := codec.UnmarshalJSON(globalFlagSpecJSON(t, tc.jsonObj))
			if err != nil {
				t.Fatalf("unmarshal json failed: %v", err)
			}

			got := globalFlagDefault(t, d, "n")
			if !reflect.DeepEqual(got, tc.wantVal) {
				t.Errorf("json default = %#v (%T), want %#v (%T)", got, got, tc.wantVal, tc.wantVal)
			}
		})
	}
}

// TestUnmarshalGlobalFlagDefaultTypeMismatch verifies that a global flag whose
// default cannot be represented as its declared type is reported at decode time.
func TestUnmarshalGlobalFlagDefaultTypeMismatch(t *testing.T) {
	cases := []struct {
		name     string
		yamlFrag string
		jsonObj  string
	}{
		{
			name:     "string default on integer flag",
			yamlFrag: "name: n\n      type: integer\n      default: not-a-number",
			jsonObj:  `{"name":"n","type":"integer","default":"not-a-number"}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name+" (yaml)", func(t *testing.T) {
			if _, err := codec.UnmarshalYAML(globalFlagSpecYAML(t, tc.yamlFrag)); err == nil {
				t.Error("expected an error for a type-mismatched global flag default, got nil")
			} else if !strings.Contains(err.Error(), `flag "n"`) {
				t.Errorf("error should name the offending flag: %v", err)
			}
		})

		t.Run(tc.name+" (json)", func(t *testing.T) {
			if _, err := codec.UnmarshalJSON(globalFlagSpecJSON(t, tc.jsonObj)); err == nil {
				t.Error("expected an error for a type-mismatched global flag default, got nil")
			} else if !strings.Contains(err.Error(), `flag "n"`) {
				t.Errorf("error should name the offending flag: %v", err)
			}
		})
	}
}
