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
		OpenCLIVersion: "1.0.0-alpha.14",
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
