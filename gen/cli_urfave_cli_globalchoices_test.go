package gen

import (
	"strings"
	"testing"

	"github.com/bcdxn/opencli/spec"
)

// TestCLI_UrfaveCliGlobalChoicesValidation ensures that choice-constrained global flags get the
// same enum metadata as command-level flags: params.gen.go emits a <BinaryPascal><FlagName> type
// with constants and IsValid(), the GlobalFlags struct field is typed with it, and every leaf
// handler hoists the GlobalFlags literal into a local so both CLI values and alternative-sourced
// resolved values are validated before being injected into context. Non-choice globals must keep
// their plain Go types (no cast, no validation) to avoid golden churn for specs without choices.
func TestCLI_UrfaveCliGlobalChoicesValidation(t *testing.T) {
	doc := &spec.Document{
		OpenCLIVersion: "1.0.0-alpha.14",
		Info:           spec.Info{Title: "GChoice CLI", Binary: "gchoice"},
		Global: &spec.Global{
			Flags: []spec.FlagItem{
				{Name: "debug", Type: "boolean", Summary: "enable verbose debug output"},
				{Name: "output-format", Type: "string", Default: "text", Choices: []spec.Choice{{Value: "text"}, {Value: "json"}}, AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "GCHOICE_OUTPUT_FORMAT"}}},
			},
		},
		Commands: &spec.CommandItem{
			Segment: "ping",
			Summary: "Report that the CLI is alive",
		},
	}

	files, err := CLI(doc, GenCLIWithFramework(UrfaveCliFramework))
	if err != nil {
		t.Fatalf("unexpected error generating urfave/cli output: %v", err)
	}

	params, ok := files["gencli/params.gen.go"]
	if !ok {
		t.Fatal("expected gencli/params.gen.go in generated output")
	}
	cmdPing, ok := files["gencli/cmd_ping.gen.go"]
	if !ok {
		t.Fatalf("unexpected set of generated files: %v", fileNames(files))
	}

	paramsStr := string(params)

	// (a) params emits the enum type + constants for the choice-constrained global.
	for _, want := range []string{
		"type GchoiceOutputFormat string",
		`GchoiceOutputFormatText GchoiceOutputFormat = "text"`,
		`GchoiceOutputFormatJson GchoiceOutputFormat = "json"`,
	} {
		if !strings.Contains(paramsStr, want) {
			t.Errorf("params.gen.go missing %q\n%s", want, paramsStr)
		}
	}

	// (b) the GlobalFlags struct field is typed with the enum; non-choice globals keep their base type.
	structBody := ""
	if start := strings.Index(paramsStr, "type GlobalFlags struct {"); start >= 0 {
		end := strings.Index(paramsStr[start:], "\n}")
		if end > 0 {
			structBody = paramsStr[start : start+end]
		}
	} else {
		t.Fatal("params.gen.go missing GlobalFlags struct")
	}
	for _, want := range []string{"OutputFormat GchoiceOutputFormat", "Debug bool"} {
		found := false
		for _, line := range strings.Split(structBody, "\n") {
			if strings.Join(strings.Fields(line), " ") == want { // Fields normalizes gofmt alignment padding
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GlobalFlags struct missing field %q\n%s", want, structBody)
		}
	}

	cmdPingStr := string(cmdPing)

	// (c) the handler hoists the literal into a local and validates before WithGlobalFlags.
	for _, want := range []string{
		"globalFlags := GlobalFlags{",
		"GchoiceOutputFormat(resolveStringFlag(c.IsSet(\"output-format\"), c.String(\"output-format\"), ",
		`if globalFlags.OutputFormat != "" && !globalFlags.OutputFormat.IsValid() {`,
		`return BadUserInput("invalid value for --output-format flag: "+string(globalFlags.OutputFormat)`,
	} {
		if !strings.Contains(cmdPingStr, want) {
			t.Errorf("cmd_ping.gen.go missing %q\n%s", want, cmdPingStr)
		}
	}

	// The resolver call must be cast to the enum so alt-source-resolved values are validated too.
	if idx := strings.Index(cmdPingStr, "globalFlags.OutputFormat.IsValid()"); idx < 0 || !strings.Contains(cmdPingStr[:idx], "GchoiceOutputFormat(resolveStringFlag") {
		t.Error("resolver call for choice global must be cast to the enum type before validation")
	}

	// Non-choice globals keep their plain reads: no cast, no IsValid check.
	if strings.Contains(cmdPingStr, "Debug: Gchoice") || strings.Contains(cmdPingStr, "globalFlags.Debug.IsValid()") {
		t.Error("non-choice global 'debug' must not be cast or validated")
	}

	// The hoisted local is what gets injected into context.
	if !strings.Contains(cmdPingStr, "ctx = WithGlobalFlags(ctx, globalFlags)") {
		t.Errorf("handler must inject the hoisted local via WithGlobalFlags\n%s", cmdPingStr)
	}
}

func fileNames(files map[string][]byte) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	return names
}
