package gen

import (
	"bytes"
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
)

//go:generate mkdir -p out
//go:generate cp -r ../examples/petstore-cli.ocs.yaml ./out/petstore-cli.ocs.yaml

//go:embed out/petstore-cli.ocs.yaml
var exampleYAML []byte // used by docs_test.go as well

var update = flag.Bool("update", false, "update golden files for CLI generation tests") // used by docs_test.go as well

func TestCLI_Cobra(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	files, err := CLI(
		doc,
		GenCLIWithFramework(CobraFramework),
	)
	if err != nil {
		t.Fatalf("unexpected error generating CLI: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("expected generated files but got none")
	}

	for relPath, content := range files {
		goldenPath := filepath.Join("testdata/cobra", relPath)

		if *update {
			if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
				t.Fatalf("failed to create golden dir for %s: %v", goldenPath, err)
			}
			if err := os.WriteFile(goldenPath, content, 0644); err != nil {
				t.Fatalf("failed to write golden file %s: %v", goldenPath, err)
			}
			continue
		}

		expected, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("failed to read golden file %s (run with -update to generate): %v", goldenPath, err)
		}

		if !bytes.Equal(content, expected) {
			t.Errorf("generated file %s does not match golden file", relPath)
		}
	}
}

func TestCLI_Yargs(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	files, err := CLI(
		doc,
		GenCLIWithFramework(YargsFramework),
	)
	if err != nil {
		t.Fatalf("unexpected error generating Yargs CLI: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("expected generated files but got none")
	}

	for relPath, content := range files {
		goldenPath := filepath.Join("testdata/yargs", relPath)

		if *update {
			if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
				t.Fatalf("failed to create golden dir for %s: %v", goldenPath, err)
			}
			if err := os.WriteFile(goldenPath, content, 0644); err != nil {
				t.Fatalf("failed to write golden file %s: %v", goldenPath, err)
			}
			continue
		}

		expected, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("failed to read golden file %s (run with -update to generate): %v", goldenPath, err)
		}

		if !bytes.Equal(content, expected) {
			t.Errorf("generated file %s does not match golden file", relPath)
		}
	}
}

func TestCLI_UrfaveCli(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	files, err := CLI(
		doc,
		GenCLIWithFramework(UrfaveCliFramework),
	)
	if err != nil {
		t.Fatalf("unexpected error generating UrfaveCli CLI: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("expected generated files but got none")
	}

	for relPath, content := range files {
		goldenPath := filepath.Join("testdata/urfavecli", relPath)

		if *update {
			if err := os.MkdirAll(filepath.Dir(goldenPath), 0755); err != nil {
				t.Fatalf("failed to create golden dir for %s: %v", goldenPath, err)
			}
			if err := os.WriteFile(goldenPath, content, 0644); err != nil {
				t.Fatalf("failed to write golden file %s: %v", goldenPath, err)
			}
			continue
		}

		expected, err := os.ReadFile(goldenPath)
		if err != nil {
			t.Fatalf("failed to read golden file %s (run with -update to generate): %v", goldenPath, err)
		}

		if !bytes.Equal(content, expected) {
			t.Errorf("generated file %s does not match golden file", relPath)
		}
	}
}

func TestCobraDefaultVal(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		t        string
		variadic bool
		want     string
	}{
		// Provided values
		{"string_val", "hello", "", false, `"hello"`},
		{"int_val", int(42), "", false, `42`},
		{"int32_val", int32(42), "", false, `42`},
		{"int64_val", int64(42), "", false, `42`},
		{"float_val", float64(3.14), "", false, `3.140000`},
		{"bool_true", true, "", false, `true`},
		{"bool_false", false, "", false, `false`},

		// Variadic defaults (canonical shapes produced by the codec)
		{"variadic_string_list", []string{"a", "b"}, "string", true, `[]string{"a", "b"}`},
		{"variadic_integer_list", []int64{1, 2, 3}, "integer", true, `[]int64{1, 2, 3}`},
		{"variadic_number_list", []float64{1.5, 2.5}, "number", true, `[]float64{1.5, 2.5}`},
		{"variadic_boolean_list", []bool{true, false}, "boolean", true, `[]bool{true, false}`},

		// Variadic zero values
		{"variadic_string", nil, "string", true, `[]string{}`},
		{"variadic_integer", nil, "integer", true, `[]int64{}`},
		{"variadic_boolean", nil, "boolean", true, `[]bool{}`},
		{"variadic_number", nil, "number", true, `[]float64{}`},

		// Non-variadic zero values
		{"zero_string", nil, "string", false, `""`},
		{"zero_integer", nil, "integer", false, `0`},
		{"zero_boolean", nil, "boolean", false, `false`},
		{"zero_number", nil, "number", false, `0`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cobraDefaultVal(tt.val, tt.t, tt.variadic)
			if got != tt.want {
				t.Errorf("cobraDefaultVal(%v, %q, %v) = %q, want %q", tt.val, tt.t, tt.variadic, got, tt.want)
			}
		})
	}
}

func TestUrfaveCliFlagStruct(t *testing.T) {
	tests := []struct {
		name     string
		t        string
		variadic bool
		want     string
	}{
		{"string", "string", false, "cli.StringFlag"},
		{"integer", "integer", false, "cli.Int64Flag"},
		{"boolean", "boolean", false, "cli.BoolFlag"},
		{"number", "number", false, "cli.Float64Flag"},
		{"variadic_string", "string", true, "cli.StringSliceFlag"},
		{"variadic_integer", "integer", true, "cli.Int64SliceFlag"},
		{"variadic_boolean", "boolean", true, "cli.BoolSliceFlag"},
		{"variadic_number", "number", true, "cli.Float64SliceFlag"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urfaveCliFlagStruct(tt.t, tt.variadic)
			if got != tt.want {
				t.Errorf("urfaveCliFlagStruct(%q, %v) = %q, want %q", tt.t, tt.variadic, got, tt.want)
			}
		})
	}
}

func TestUrfaveCliAccessor(t *testing.T) {
	tests := []struct {
		name     string
		t        string
		variadic bool
		want     string
	}{
		{"string", "string", false, "String"},
		{"integer", "integer", false, "Int64"},
		{"boolean", "boolean", false, "Bool"},
		{"number", "number", false, "Float64"},
		{"variadic_string", "string", true, "StringSlice"},
		{"variadic_integer", "integer", true, "Int64Slice"},
		{"variadic_boolean", "boolean", true, "BoolSlice"},
		{"variadic_number", "number", true, "Float64Slice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urfaveCliAccessor(tt.t, tt.variadic)
			if got != tt.want {
				t.Errorf("urfaveCliAccessor(%q, %v) = %q, want %q", tt.t, tt.variadic, got, tt.want)
			}
		})
	}
}

func TestUrfaveCliZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		t        string
		variadic bool
		want     string
	}{
		{"string", "string", false, `""`},
		{"integer", "integer", false, "0"},
		{"boolean", "boolean", false, "false"},
		{"number", "number", false, "0.0"},
		{"variadic_string", "string", true, "[]string{}"},
		{"variadic_integer", "integer", true, "[]int64{}"},
		{"variadic_boolean", "boolean", true, "[]bool{}"},
		{"variadic_number", "number", true, "[]float64{}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urfaveCliZeroValue(tt.t, tt.variadic)
			if got != tt.want {
				t.Errorf("urfaveCliZeroValue(%q, %v) = %q, want %q", tt.t, tt.variadic, got, tt.want)
			}
		})
	}
}

func TestUrfaveCliDefaultVal(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		t        string
		variadic bool
		want     string
	}{
		// Provided scalars (canonical shapes produced by the codec)
		{"string_val", "hello", "", false, `"hello"`},
		{"int64_val", int64(42), "", false, `42`},
		{"uint64_val", uint64(42), "", false, `42`},
		{"float_val", float64(3.5), "", false, `3.500000`},
		{"bool_true", true, "", false, `true`},

		// Variadic defaults (canonical shapes produced by the codec)
		{"variadic_string_list", []string{"a", "b"}, "string", true, `[]string{"a", "b"}`},
		{"variadic_integer_list", []int64{1, 2, 3}, "integer", true, `[]int64{1, 2, 3}`},
		{"variadic_number_list", []float64{1.5, 2.5}, "number", true, `[]float64{1.5, 2.5}`},
		{"variadic_boolean_list", []bool{true, false}, "boolean", true, `[]bool{true, false}`},

		// No default -> zero value
		{"zero_string", nil, "string", false, `""`},
		{"zero_integer", nil, "integer", false, `0`},
		{"zero_boolean", nil, "boolean", false, `false`},
		{"zero_number", nil, "number", false, `0.0`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urfaveCliDefaultVal(tt.val, tt.t, tt.variadic)
			if got != tt.want {
				t.Errorf("urfaveCliDefaultVal(%v, %q, %v) = %q, want %q", tt.val, tt.t, tt.variadic, got, tt.want)
			}
		})
	}
}

func TestYargsDefaultVal(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
	}{
		{"string", "hello", `"hello"`},
		{"int", int(42), `42`},
		{"int32", int32(42), `42`},
		{"int64", int64(42), `42`},
		{"float", float64(3.14), `3.140000`},
		{"bool_true", true, `true`},
		{"bool_false", false, `false`},

		// Variadic defaults (canonical shapes produced by the codec)
		{"string_list", []string{"a", "b"}, `["a", "b"]`},
		{"int64_list", []int64{1, 2, 3}, `[1, 2, 3]`},
		{"float_list", []float64{1.5, 2.5}, `[1.5, 2.5]`},
		{"bool_list", []bool{true, false}, `[true, false]`},

		{"nil", nil, ``},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := yargsDefaultVal(tt.val)
			if got != tt.want {
				t.Errorf("yargsDefaultVal(%v) = %q, want %q", tt.val, got, tt.want)
			}
		})
	}
}

func TestCLIFrameworkIsValid(t *testing.T) {
	tests := []struct {
		f    CLIFramework
		want bool
	}{
		{CobraFramework, true},
		{YargsFramework, true},
		{UrfaveCliFramework, true},
		{"INVALID", false},
		{"", false},
	}
	for _, tt := range tests {
		got := tt.f.IsValid()
		if got != tt.want {
			t.Errorf("CLIFramework(%q).IsValid() = %v, want %v", tt.f, got, tt.want)
		}
	}
}

func TestCLI_NilDoc(t *testing.T) {
	_, err := CLI(nil)
	if err == nil || err.Error() != "provided specification document was nil" {
		t.Fatalf("expected nil document error but got %s", err)
	}
}

func TestCLI_InvalidFramework(t *testing.T) {
	doc := &spec.Document{}
	_, err := CLI(doc, GenCLIWithFramework("INVALID"))
	if err == nil {
		t.Fatal("expected error for invalid framework")
	}
}

func TestDocFormatIsValid(t *testing.T) {
	tests := []struct {
		f    DocFormat
		want bool
	}{
		{Markdown, true},
		{HTML_PAGE, true},
		{HTML_EMBED, true},
		{MAN, true},
		{"INVALID", false},
	}
	for _, tt := range tests {
		got := tt.f.IsValid()
		if got != tt.want {
			t.Errorf("DocFormat(%q).IsValid() = %v, want %v", tt.f, got, tt.want)
		}
	}
}

func TestDocs_InvalidFormat(t *testing.T) {
	doc := &spec.Document{}
	_, err := Docs(doc, DocsWithFormat("INVALID"))
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestDocs_ManPageNotImplemented(t *testing.T) {
	doc := &spec.Document{}
	_, err := Docs(doc, DocsWithFormat(MAN))
	if err == nil {
		t.Fatal("expected error for ManPage format")
	}
}

func TestDocsOptions(t *testing.T) {
	// Verify functional options work without panic
	doc := &spec.Document{}
	_, _ = Docs(doc, DocsWithoutFooter())
	_, _ = Docs(doc, DocsWithoutBadge())
	_, _ = Docs(doc, DocsWithoutFooter(), DocsWithoutBadge())
}

func TestDocs_HTMLPage(t *testing.T) {
	doc := &spec.Document{}
	out, err := Docs(doc, DocsWithFormat(HTML_PAGE))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("expected HTML output")
	}
}

func TestCobraPlainFlagExpr(t *testing.T) {
	tests := []struct {
		name string
		f    cobraFlagEntry
		want string
	}{
		{"bare_var", cobraFlagEntry{VarName: "flagUsername"}, "flagUsername"},
		{"cast_to_type", cobraFlagEntry{VarName: "flagStatus", TypeName: "PetstoreStatus"}, "PetstoreStatus(flagStatus)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := plainCobraFlagExpr(tt.f); got != tt.want {
				t.Errorf("plainCobraFlagExpr(%+v) = %q, want %q", tt.f, got, tt.want)
			}
		})
	}
}

// TestCobraResolveFlagValue exercises the resolveFlagValue template function exactly as the
// command.tmpl uses it (pulled from cobraTemplateFuncMap so there is a single source of truth).
func TestCobraResolveFlagValue(t *testing.T) {
	resolve := func(f cobraFlagEntry) string {
		fn, ok := cobraTemplateFuncMap()["resolveFlagValue"].(func(cobraFlagEntry) string)
		if !ok {
			t.Fatal("resolveFlagValue not found in cobra template func map")
		}
		return fn(f)
	}

	envUser := spec.AlternativeSource{Type: "$ENV", Property: "PETSTORE_USER"}
	fileAuth := spec.AlternativeSource{Type: "$FILE", Property: "$.auth.user"}

	tests := []struct {
		name string
		f    cobraFlagEntry
		want string
	}{
		// No alternative sources -> plain bound variable (optionally cast).
		{"no_alt_bare", cobraFlagEntry{VarName: "flagUsername"}, "flagUsername"},
		{"no_alt_cast", cobraFlagEntry{VarName: "flagStatus", TypeName: "PetstoreStatus"}, "PetstoreStatus(flagStatus)"},

		// Single $ENV source.
		{"env_only",
			cobraFlagEntry{VarName: "flagUsername", FlagName: "username", GoType: "string", AltSources: []spec.AlternativeSource{envUser}},
			`resolveStringFlag(c.Flags(), "username", []AltSource{{Type: "$ENV", Property: "PETSTORE_USER"}})`},

		// Single $FILE source.
		{"file_only",
			cobraFlagEntry{VarName: "flagUsername", FlagName: "username", GoType: "string", AltSources: []spec.AlternativeSource{fileAuth}},
			`resolveStringFlag(c.Flags(), "username", []AltSource{{Type: "$FILE", Property: "$.auth.user"}})`},

		// Mixed sources preserve declared order ($ENV before $FILE).
		{"env_then_file_order",
			cobraFlagEntry{VarName: "flagUsername", FlagName: "username", GoType: "string", AltSources: []spec.AlternativeSource{envUser, fileAuth}},
			`resolveStringFlag(c.Flags(), "username", []AltSource{{Type: "$ENV", Property: "PETSTORE_USER"}, {Type: "$FILE", Property: "$.auth.user"}})`},

		// Resolver is selected by Go type; a generated choices type wraps the call.
		{"int64_env",
			cobraFlagEntry{VarName: "flagLimit", FlagName: "limit", GoType: "int64", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "LIMIT"}}},
			`resolveInt64Flag(c.Flags(), "limit", []AltSource{{Type: "$ENV", Property: "LIMIT"}})`},

		{"bool_env_cast",
			cobraFlagEntry{VarName: "flagVerbose", FlagName: "verbose", GoType: "bool", TypeName: "Verbosity", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "VERBOSE"}}},
			`Verbosity(resolveBoolFlag(c.Flags(), "verbose", []AltSource{{Type: "$ENV", Property: "VERBOSE"}}))`},

		{"float64_env",
			cobraFlagEntry{VarName: "flagRate", FlagName: "rate", GoType: "float64", AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.rate"}}},
			`resolveFloat64Flag(c.Flags(), "rate", []AltSource{{Type: "$FILE", Property: "$.rate"}})`},

		// Variadic types map to their slice resolvers (GetStringArray-backed for strings).
		{"string_slice_env",
			cobraFlagEntry{VarName: "flagTags", FlagName: "tags", GoType: "[]string", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "TAGS"}}},
			`resolveStringSliceFlag(c.Flags(), "tags", []AltSource{{Type: "$ENV", Property: "TAGS"}})`},

		{"int64_slice_env",
			cobraFlagEntry{VarName: "flagIds", FlagName: "ids", GoType: "[]int64", AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.ids"}}},
			`resolveInt64SliceFlag(c.Flags(), "ids", []AltSource{{Type: "$FILE", Property: "$.ids"}})`},

		{"bool_slice_env",
			cobraFlagEntry{VarName: "flagFlags", FlagName: "flags", GoType: "[]bool", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "FLAGS"}}},
			`resolveBoolSliceFlag(c.Flags(), "flags", []AltSource{{Type: "$ENV", Property: "FLAGS"}})`},

		{"float64_slice_env",
			cobraFlagEntry{VarName: "flagRates", FlagName: "rates", GoType: "[]float64", AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.rates"}}},
			`resolveFloat64SliceFlag(c.Flags(), "rates", []AltSource{{Type: "$FILE", Property: "$.rates"}})`},

		// Unknown Go type with alt sources falls back to the plain expression (defensive).
		{"unknown_type_falls_back",
			cobraFlagEntry{VarName: "flagWeird", FlagName: "weird", GoType: "int32", AltSources: []spec.AlternativeSource{envUser}},
			"flagWeird"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.f); got != tt.want {
				t.Errorf("resolveFlagValue(%+v)\n  = %q\nwant %q", tt.f, got, tt.want)
			}
		})
	}
}

// TestCobraScanAltSources verifies the has-alt / has-file scan that gates emission of
// gencli/config.gen.go and its JSONPath import.
func TestCobraScanAltSources(t *testing.T) {
	tests := []struct {
		name     string
		flags    []cobraFlagEntry
		wantHas  bool
		wantFile bool
	}{
		{"none", nil, false, false},
		{
			"env_only",
			[]cobraFlagEntry{{AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "X"}}}},
			true, false,
		},
		{
			"file_only",
			[]cobraFlagEntry{{AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.x"}}}},
			true, true,
		},
		{
			"mixed_across_flags",
			[]cobraFlagEntry{
				{}, // no alt sources
				{AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "X"}, {Type: "$FILE", Property: "$.y"}}},
			},
			true, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAlt, hasFile := scanCobraAltSources(tt.flags)
			if hasAlt != tt.wantHas || hasFile != tt.wantFile {
				t.Errorf("scanCobraAltSources(%+v) = (%v, %v), want (%v, %v)", tt.flags, hasAlt, hasFile, tt.wantHas, tt.wantFile)
			}
		})
	}
}

// TestUrfaveCliScanAltSources verifies the has-alt / has-file scan that gates emission of
// gencli/config.gen.go and its JSONPath import for the urfave/cli framework.
func TestUrfaveCliScanAltSources(t *testing.T) {
	tests := []struct {
		name     string
		flags    []urfaveCliFlagEntry
		wantHas  bool
		wantFile bool
	}{
		{"none", nil, false, false},
		{
			"env_only",
			[]urfaveCliFlagEntry{{AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "X"}}}},
			true, false,
		},
		{
			"file_only",
			[]urfaveCliFlagEntry{{AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.x"}}}},
			true, true,
		},
		{
			"mixed_across_flags",
			[]urfaveCliFlagEntry{
				{}, // no alt sources
				{AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "X"}, {Type: "$FILE", Property: "$.y"}}},
			},
			true, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAlt, hasFile := scanUrfaveCliAltSources(tt.flags)
			if hasAlt != tt.wantHas || hasFile != tt.wantFile {
				t.Errorf("scanUrfaveCliAltSources(%+v) = (%v, %v), want (%v, %v)", tt.flags, hasAlt, hasFile, tt.wantHas, tt.wantFile)
			}
		})
	}
}

// TestCLI_UrfaveCli_ConfigEmissionConditional verifies that gencli/config.gen.go is only
// emitted when at least one flag declares an alternative source. Specs without this feature
// must produce no extra config file (and therefore no unused imports), matching the cobra
// framework's behavior.
func TestCLI_UrfaveCli_ConfigEmissionConditional(t *testing.T) {
	generate := func(t *testing.T, docYAML string) map[string][]byte {
		t.Helper()
		doc, err := codec.UnmarshalYAML([]byte(docYAML))
		if err != nil {
			t.Fatalf("unexpected error unmarshaling doc: %v", err)
		}
		files, err := CLI(
			doc,
			GenCLIWithFramework(UrfaveCliFramework),
		)
		if err != nil {
			t.Fatalf("unexpected error generating UrfaveCli CLI: %v", err)
		}
		return files
	}

	fileNames := func(files map[string][]byte) []string {
		names := make([]string, 0, len(files))
		for n := range files {
			names = append(names, n)
		}
		return names
	}

	t.Run("no_alt_sources", func(t *testing.T) {
		files := generate(t, `opencliVersion: `+spec.SchemaVersion+`
info:
  title: minimal cli for alt-source emission test
  binary: minicli
commands:
  minicli greet [flags]:
    summary: say hello
`)

		if _, hasConfig := files["gencli/config.gen.go"]; hasConfig {
			t.Errorf("config.gen.go should not be emitted without alternative sources; generated files: %v", fileNames(files))
		}

		runContent, ok := files["gencli/run.gen.go"]
		if !ok {
			t.Fatal("expected gencli/run.gen.go in output")
		}
		if bytes.Contains(runContent, []byte("loadConfig()")) {
			t.Error("run.gen.go should not call loadConfig without alternative sources")
		}
	})

	t.Run("with_alt_sources", func(t *testing.T) {
		files := generate(t, `opencliVersion: `+spec.SchemaVersion+`
info:
  title: minimal cli for alt-source emission test
  binary: minicli
commands:
  minicli greet [flags]:
    summary: say hello
    flags:
      - name: username
        type: string
        alternativeSources:
          - type: $ENV
            property: MINI_USER
`)

		if _, hasConfig := files["gencli/config.gen.go"]; !hasConfig {
			t.Errorf("expected config.gen.go when alt sources present; generated files: %v", fileNames(files))
		}

		runContent, ok := files["gencli/run.gen.go"]
		if !ok {
			t.Fatal("expected gencli/run.gen.go in output")
		}
		if !bytes.Contains(runContent, []byte("loadConfig()")) {
			t.Error("expected run.gen.go to call loadConfig when alt sources present")
		}
	})
}

// TestYargsScanAltSources verifies the has-alt / has-file scan that gates emission of
// gencli/config.ts and its JSONPath import for the yargs framework.
func TestYargsScanAltSources(t *testing.T) {
	tests := []struct {
		name     string
		flags    []yargsFlagEntry
		wantHas  bool
		wantFile bool
	}{
		{"none", nil, false, false},
		{
			"env_only",
			[]yargsFlagEntry{{AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "X"}}}},
			true, false,
		},
		{
			"file_only",
			[]yargsFlagEntry{{AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.x"}}}},
			true, true,
		},
		{
			"mixed_across_flags",
			[]yargsFlagEntry{
				{}, // no alt sources
				{AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "X"}, {Type: "$FILE", Property: "$.y"}}},
			},
			true, true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasAlt, hasFile := scanYargsAltSources(tt.flags)
			if hasAlt != tt.wantHas || hasFile != tt.wantFile {
				t.Errorf("scanYargsAltSources(%+v) = (%v, %v), want (%v, %v)", tt.flags, hasAlt, hasFile, tt.wantHas, tt.wantFile)
			}
		})
	}
}

// TestYargsAltSourceNames verifies the accepted-option-name list passed to the generated
// wasSetOnCli scanner: field name first (used to read argv), then raw name and shorthand
// when present, then extra aliases — with duplicates collapsed. The shorthand must be
// included so that a flag set via its single-character form is still detected as CLI-set.
func TestYargsAltSourceNames(t *testing.T) {
	tests := []struct {
		name string
		flag yargsFlagEntry
		want []string
	}{
		{"field_only", yargsFlagEntry{FieldName: "verbose"}, []string{"verbose"}},
		{"kebab_raw_name", yargsFlagEntry{FieldName: "dryRun", RawName: "dry-run"}, []string{"dryRun", "dry-run"}},
		{"shorthand_included", yargsFlagEntry{FieldName: "verbose", RawName: "verbose", Shorthand: "v"}, []string{"verbose", "v"}},
		{"all_names_deduped", yargsFlagEntry{FieldName: "output", RawName: "out-file", Shorthand: "o", ExtraAliases: []string{"dest"}}, []string{"output", "out-file", "o", "dest"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := yargsAltSourceNames(tt.flag)
			if !slices.Equal(got, tt.want) {
				t.Errorf("yargsAltSourceNames(%+v) = %v, want %v", tt.flag, got, tt.want)
			}
		})
	}
}

// TestYargsResolveFlagValue exercises the resolveFlagValue template function exactly as the
// command.tmpl uses it (pulled from yargsTemplateFuncMap so there is a single source of
// truth). Choice-typed flags with alternative sources must be wrapped in assertChoice:
// yargs only validates values parsed from the command line against .choices(), so a value
// resolved from $ENV or $FILE has to be checked before it reaches the action.
func TestYargsResolveFlagValue(t *testing.T) {
	resolve := func(f yargsFlagEntry) string {
		fn, ok := yargsTemplateFuncMap()["resolveFlagValue"].(func(yargsFlagEntry) string)
		if !ok {
			t.Fatal("resolveFlagValue not found in yargs template func map")
		}
		return fn(f)
	}

	envUser := spec.AlternativeSource{Type: "$ENV", Property: "PETSTORE_USER"}
	fileAuth := spec.AlternativeSource{Type: "$FILE", Property: "$.auth.user"}

	tests := []struct {
		name string
		f    yargsFlagEntry
		want string
	}{
		// No alternative sources -> plain argv field (optionally cast).
		{"no_alt_bare", yargsFlagEntry{FieldName: "username", TSType: "string"}, `argv.username as string`},
		{"no_alt_cast", yargsFlagEntry{FieldName: "status", TypeName: "PetstoreStatus"}, `argv.status as PetstoreStatus`},

		// Single $ENV source.
		{"env_only",
			yargsFlagEntry{FieldName: "username", RawName: "username", TSType: "string", AltSources: []spec.AlternativeSource{envUser}},
			`resolveStringFlag(argv, ["username"], [{ type: "$ENV", property: "PETSTORE_USER" }])`},

		// Single $FILE source.
		{"file_only",
			yargsFlagEntry{FieldName: "username", RawName: "username", TSType: "string", AltSources: []spec.AlternativeSource{fileAuth}},
			`resolveStringFlag(argv, ["username"], [{ type: "$FILE", property: "$.auth.user" }])`},

		// Mixed sources preserve declared order ($ENV before $FILE).
		{"env_then_file_order",
			yargsFlagEntry{FieldName: "username", RawName: "username", TSType: "string", AltSources: []spec.AlternativeSource{envUser, fileAuth}},
			`resolveStringFlag(argv, ["username"], [{ type: "$ENV", property: "PETSTORE_USER" }, { type: "$FILE", property: "$.auth.user" }])`},

		// Resolver is selected by TS type; non-choice flags are not wrapped in assertChoice.
		{"number_env",
			yargsFlagEntry{FieldName: "limit", RawName: "limit", TSType: "number", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "LIMIT"}}},
			`resolveNumberFlag(argv, ["limit"], [{ type: "$ENV", property: "LIMIT" }])`},

		{"bool_env",
			yargsFlagEntry{FieldName: "verbose", RawName: "verbose", TSType: "boolean", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "VERBOSE"}}},
			`resolveBoolFlag(argv, ["verbose"], [{ type: "$ENV", property: "VERBOSE" }])`},

		// Choice-typed flag with alt sources is wrapped in assertChoice so $ENV/$FILE-resolved
		// values are validated against the declared choices before reaching the action.
		{"choices_env_wrapped",
			yargsFlagEntry{FieldName: "status", RawName: "status", TSType: "string", TypeName: "PetstoreStatus", Choices: []yargsChoiceEntry{{EnumKey: "AVAILABLE", Value: "available"}, {EnumKey: "SOLD_OUT", Value: "sold-out"}}, AltSources: []spec.AlternativeSource{envUser}},
			`assertChoice("status", (resolveStringFlag(argv, ["status"], [{ type: "$ENV", property: "PETSTORE_USER" }])) as PetstoreStatus | undefined, ["available", "sold-out"])`},

		// Choice-typed flag with both $ENV and $FILE sources.
		{"choices_env_file_wrapped",
			yargsFlagEntry{FieldName: "status", RawName: "status", TSType: "string", TypeName: "PetstoreStatus", Choices: []yargsChoiceEntry{{EnumKey: "AVAILABLE", Value: "available"}}, AltSources: []spec.AlternativeSource{envUser, fileAuth}},
			`assertChoice("status", (resolveStringFlag(argv, ["status"], [{ type: "$ENV", property: "PETSTORE_USER" }, { type: "$FILE", property: "$.auth.user" }])) as PetstoreStatus | undefined, ["available"])`},

		// Unknown TS type with alt sources falls back to the plain expression (defensive).
		{"unknown_type_falls_back",
			yargsFlagEntry{FieldName: "weird", RawName: "weird", TSType: "int32", AltSources: []spec.AlternativeSource{envUser}},
			`argv.weird as int32`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolve(tt.f); got != tt.want {
				t.Errorf("resolveFlagValue(%+v)\n  = %q\nwant %q", tt.f, got, tt.want)
			}
		})
	}
}

// TestCLI_Yargs_ConfigEmissionConditional verifies that gencli/config.ts is only emitted when
// at least one flag declares an alternative source. Specs without this feature must produce no
// extra config file (and therefore no unused imports), matching the cobra/urfave behavior.
func TestCLI_Yargs_ConfigEmissionConditional(t *testing.T) {
	generate := func(t *testing.T, docYAML string) map[string][]byte {
		t.Helper()
		doc, err := codec.UnmarshalYAML([]byte(docYAML))
		if err != nil {
			t.Fatalf("unexpected error unmarshaling doc: %v", err)
		}
		files, err := CLI(
			doc,
			GenCLIWithFramework(YargsFramework),
		)
		if err != nil {
			t.Fatalf("unexpected error generating Yargs CLI: %v", err)
		}
		return files
	}

	fileNames := func(files map[string][]byte) []string {
		names := make([]string, 0, len(files))
		for n := range files {
			names = append(names, n)
		}
		return names
	}

	t.Run("no_alt_sources", func(t *testing.T) {
		files := generate(t, `opencliVersion: 1.0.0-alpha.14
info:
  title: minimal cli for alt-source emission test
  binary: minicli
commands:
  minicli greet [flags]:
    summary: say hello
`)

		if _, hasConfig := files["gencli/config.ts"]; hasConfig {
			t.Errorf("config.ts should not be emitted without alternative sources; generated files: %v", fileNames(files))
		}

		runContent, ok := files["gencli/run.ts"]
		if !ok {
			t.Fatal("expected gencli/run.ts in output")
		}
		if bytes.Contains(runContent, []byte("loadConfig()")) {
			t.Error("run.ts should not call loadConfig without alternative sources")
		}
	})

	t.Run("with_alt_sources", func(t *testing.T) {
		files := generate(t, `opencliVersion: 1.0.0-alpha.14
info:
  title: minimal cli for alt-source emission test
  binary: minicli
commands:
  minicli greet [flags]:
    summary: say hello
    flags:
      - name: username
        type: string
        alternativeSources:
          - type: $ENV
            property: MINI_USER
`)

		if _, hasConfig := files["gencli/config.ts"]; !hasConfig {
			t.Errorf("expected config.ts when alt sources present; generated files: %v", fileNames(files))
		}

		runContent, ok := files["gencli/run.ts"]
		if !ok {
			t.Fatal("expected gencli/run.ts in output")
		}
		if !bytes.Contains(runContent, []byte("loadConfig()")) {
			t.Error("expected run.ts to call loadConfig when alt sources present")
		}
	})
}
