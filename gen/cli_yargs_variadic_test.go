package gen

import (
	"strings"
	"testing"

	"github.com/bcdxn/opencli/spec"
)

// TestCLI_YargsVariadicNonString ensures that variadic flags with non-string types
// generate correct TypeScript: the local argv interface uses the proper array type
// (number[], boolean[]) instead of hardcoded string[], and the .option() builder
// emits a coerce function rather than string:true so runtime values are correctly
// typed. String variadics must remain byte-identical to prior output.
func TestCLI_YargsVariadicNonString(t *testing.T) {
	doc := &spec.Document{
		OpenCLIVersion: spec.SchemaVersion,
		Info:           spec.Info{Title: "VarTest CLI", Binary: "vartest"},
		Commands: &spec.CommandItem{
			Segment: "greet",
			Flags: []spec.FlagItem{
				{Name: "ids", Type: "integer", Variadic: true},
				{Name: "names", Type: "string", Variadic: true},
				{Name: "verbose", Type: "boolean"}, // non-variadic control
			},
		},
	}

	files, err := CLI(doc, GenCLIWithFramework(YargsFramework))
	if err != nil {
		t.Fatalf("unexpected error generating yargs output: %v", err)
	}

	var all strings.Builder
	for _, content := range files {
		all.Write(content)
		all.WriteByte('\n')
	}
	got := all.String()

	// Integer variadic: interface must say number[] (not string[]) and builder must coerce.
	if !strings.Contains(got, `number[] | undefined`) {
		t.Error("expected 'number[] | undefined' in local argv interface for integer variadic flag")
	}
	if strings.Contains(got, `"ids"?: string[]`) {
		t.Error("integer variadic flag must NOT be typed as string[] in the local argv interface")
	}
	expectedCoerce := `coerce: (v) => Array.isArray(v) ? v.map(Number) : v`
	if !strings.Contains(got, expectedCoerce) {
		t.Errorf("expected coerce function %q for integer variadic flag", expectedCoerce)
	}

	// String variadic: must remain string[] | undefined AND keep string:true (golden parity).
	if !strings.Contains(got, `string[] | undefined`) {
		t.Error("expected 'string[] | undefined' in local argv interface for string variadic flag")
	}
	if !strings.Contains(got, "string: true,") {
		t.Error("string variadic flag must still emit 'string: true,' in .option() builder (golden parity)")
	}

	// The integer variadic must NOT have string:true on it. Verify by checking that
	// the coerce line is present and there's no "string: true" immediately after a
	// type:"array" for ids. Simplest check: count occurrences of 'string: true,' —
	// should be exactly 1 (only from the names flag).
	if n := strings.Count(got, "string: true,"); n != 1 {
		t.Errorf("expected exactly 1 occurrence of 'string: true,' (for string variadic only), got %d", n)
	}

	// Non-variadic boolean control: should render as plain boolean (no array suffix).
	if !strings.Contains(got, `"verbose"?: boolean;`) {
		t.Error("expected non-variadic boolean flag to render as 'boolean' in interface")
	}
}

// TestCLI_YargsGlobalVariadicNonString ensures that global variadic flags with
// non-string types are registered in run.ts with the same coerce functions used
// for command-level flags, so their runtime values match the typed arrays declared
// on GlobalFlags. String variadics must keep string:true (golden parity).
func TestCLI_YargsGlobalVariadicNonString(t *testing.T) {
	doc := &spec.Document{
		OpenCLIVersion: spec.SchemaVersion,
		Info:           spec.Info{Title: "GVarTest CLI", Binary: "gvar"},
		Global: &spec.Global{
			Config: spec.Configuration{},
			Flags: []spec.FlagItem{
				{Name: "ids", Type: "integer", Variadic: true},
				{Name: "ratios", Type: "number", Variadic: true},
				{Name: "names", Type: "string", Variadic: true},
			},
		},
		Commands: &spec.CommandItem{
			Segment: "ping",
		},
	}

	files, err := CLI(doc, GenCLIWithFramework(YargsFramework))
	if err != nil {
		t.Fatalf("unexpected error generating yargs output: %v", err)
	}

	runTS, ok := files["gencli/run.ts"]
	if !ok {
		t.Fatal("expected gencli/run.ts in generated files")
	}
	got := string(runTS)

	intCoerce := `coerce: (v) => Array.isArray(v) ? v.map(Number) : v`
	boolCoerce := `(v) => Array.isArray(v) ? v.map((x) => x === true || x === "true") : v`

	// Integer and number variadics must coerce to numbers.
	if n := strings.Count(got, intCoerce); n != 2 {
		t.Errorf("expected %d occurrences of the numeric coerce function (ids + ratios), got %d", 2, n)
	}
	// String variadic keeps string:true for golden parity; no boolean global here.
	if !strings.Contains(got, "string: true,") {
		t.Error("global string variadic flag must still emit 'string: true,' in run.ts (golden parity)")
	}
	if strings.Count(got, boolCoerce) != 0 {
		t.Errorf("unexpected boolean coerce function %q for non-boolean global flags", boolCoerce)
	}

	paramsTS, ok := files["gencli/params.ts"]
	if !ok {
		t.Fatal("expected gencli/params.ts in generated files")
	}
	gotParams := string(paramsTS)
	for _, want := range []string{`ids?: number[] | undefined;`, `ratios?: number[] | undefined;`, `names?: string[] | undefined;`} {
		if !strings.Contains(gotParams, want) {
			t.Errorf("expected %q in GlobalFlags interface", want)
		}
	}
}

// TestCLI_YargsGlobalVariadicBoolean ensures a global boolean variadic flag gets
// the boolean coerce function rather than string:true.
func TestCLI_YargsGlobalVariadicBoolean(t *testing.T) {
	doc := &spec.Document{
		OpenCLIVersion: spec.SchemaVersion,
		Info:           spec.Info{Title: "GBoolTest CLI", Binary: "gbool"},
		Global: &spec.Global{
			Config: spec.Configuration{},
			Flags: []spec.FlagItem{
				{Name: "verbose", Type: "boolean", Variadic: true},
			},
		},
		Commands: &spec.CommandItem{
			Segment: "ping",
		},
	}

	files, err := CLI(doc, GenCLIWithFramework(YargsFramework))
	if err != nil {
		t.Fatalf("unexpected error generating yargs output: %v", err)
	}

	runTS, ok := files["gencli/run.ts"]
	if !ok {
		t.Fatal("expected gencli/run.ts in generated files")
	}
	got := string(runTS)

	boolCoerce := `coerce: (v) => Array.isArray(v) ? v.map((x) => x === true || x === "true") : v`
	if !strings.Contains(got, boolCoerce) {
		t.Errorf("expected boolean coerce function %q for global boolean variadic flag", boolCoerce)
	}
	if strings.Contains(got, "string: true,") {
		t.Error("global boolean variadic flag must NOT emit 'string: true,' in run.ts")
	}

	paramsTS := string(files["gencli/params.ts"])
	if !strings.Contains(paramsTS, `verbose?: boolean[] | undefined;`) {
		t.Error(`expected "verbose?: boolean[] | undefined;" in GlobalFlags interface`)
	}
}

// TestCLI_YargsRequiredVariadicFlag ensures a flag that is both required and
// variadic is registered as a demanded array option, so yargs rejects zero
// occurrences while still accepting many values.
func TestCLI_YargsRequiredVariadicFlag(t *testing.T) {
	doc := &spec.Document{
		OpenCLIVersion: "1.0.0-alpha.14",
		Info:           spec.Info{Title: "VarTest CLI", Binary: "vartest"},
		Commands: &spec.CommandItem{
			Segment: "greet",
			Flags: []spec.FlagItem{
				{Name: "items", Type: "string", Variadic: true, Required: true},
			},
		},
	}

	files, err := CLI(doc, GenCLIWithFramework(YargsFramework))
	if err != nil {
		t.Fatalf("unexpected error generating yargs output: %v", err)
	}

	var all strings.Builder
	for _, content := range files {
		all.Write(content)
		all.WriteByte('\n')
	}
	got := all.String()

	if !strings.Contains(got, `type: "array",`) {
		t.Error("expected required variadic flag to be registered as an array option")
	}
	if !strings.Contains(got, "demandOption: true,") {
		t.Error("expected required variadic flag to emit demandOption: true")
	}
	if !strings.Contains(got, `"items": string[]`) {
		t.Error("expected required variadic flag to be a non-optional string[] field in the argv interface")
	}
}
