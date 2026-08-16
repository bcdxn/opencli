package gen

import (
	"bytes"
	_ "embed"
	"flag"
	"os"
	"path/filepath"
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
