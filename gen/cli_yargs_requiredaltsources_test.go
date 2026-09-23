package gen

import (
	"strings"
	"testing"

	"github.com/bcdxn/opencli/spec"
)

// TestCLI_YargsRequiredAltSourcesValidation ensures that a required flag with alternative
// sources gets post-resolution validation in the handler: demandOption is skipped for such
// flags (a CLI-absent value may still resolve from env/config), so without an explicit check
// the action would receive undefined when every source is absent. Optional alt-sourced flags
// must not be validated, and required flags WITHOUT alternative sources keep plain demandOption.
func TestCLI_YargsRequiredAltSourcesValidation(t *testing.T) {
	doc := &spec.Document{
		OpenCLIVersion: spec.SchemaVersion,
		Info:           spec.Info{Title: "ReqAlt CLI", Binary: "reqalt"},
		Commands: &spec.CommandItem{
			Segment: "login",
			Flags: []spec.FlagItem{
				{Name: "username", Type: "string", Required: true, AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "REQALT_USER"}}},
				{Name: "password", Type: "string", AltSources: []spec.AlternativeSource{{Type: "$FILE", Property: "$.auth.pass"}}},
				{Name: "region", Type: "string", Required: true, Default: "us-east-1"},
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

	// Required + alt-sourced flag: validated after resolution, no demandOption.
	if !strings.Contains(got, `if (cmdFlags.username === undefined) _missing.push("--username");`) {
		t.Error("expected post-resolution validation for required alt-sourced flag 'username'")
	}

	// Required non-alt-sourced flag: keeps plain demandOption and is NOT validated.
	if n := strings.Count(got, "demandOption: true,"); n != 1 {
		t.Errorf("expected exactly one demandOption (for 'region'), got %d", n)
	}
	if strings.Contains(got, `_missing.push("--region")`) {
		t.Error("required non-alt-sourced flag 'region' must not be validated post-resolution")
	}

	// Optional alt-sourced flag: no validation, and its option block has no demandOption.
	if strings.Contains(got, `_missing.push("--password")`) {
		t.Error("optional alt-sourced flag 'password' must not be validated as required")
	}

	// The rejection path reports the missing flags with yargs-style wording so the root
	// .fail() handler maps it to a bad-user-input exit code.
	if !strings.Contains(got, `Missing required arguments: ${_missing.join(", ")}`) {
		t.Error("expected 'Missing required arguments' error message in validation block")
	}

	// Typing must reflect that alt-sourced flags may be absent from the parsed argv even
	// when required (demandOption is skipped), so both the local argv interface and the
	// action params type them as T | undefined. Required non-alt-sourced flags keep a
	// plain non-optional type, and optional alt-sourced flags stay optional in params.
	if !strings.Contains(got, `"username": string | undefined;`) {
		t.Error("local argv interface must type required alt-sourced flag as 'string | undefined'")
	}
	if !strings.Contains(got, "  username: string | undefined;\n") {
		t.Error("params flags interface must type required alt-sourced flag as 'string | undefined'")
	}
	if !strings.Contains(got, `"password": string | undefined;`) || strings.Contains(got, `"password"?:`) {
		t.Error("local argv interface must type optional alt-sourced flag non-optionally with an undefined union")
	}
	if !strings.Contains(got, "  password?: string | undefined;\n") {
		t.Error("params flags interface must keep optional alt-sourced flag optional with an undefined union")
	}
	if strings.Contains(got, "  region: string | undefined;") || !strings.Contains(got, "  region: string;\n") {
		t.Error("required non-alt-sourced flag must stay a plain non-optional type in params")
	}

	// The whole block must be absent when no flag is both required and alt-sourced.
	doc2 := *doc
	doc2.Commands = &spec.CommandItem{
		Segment: "login",
		Flags: []spec.FlagItem{
			{Name: "username", Type: "string", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "REQALT_USER"}}},
		},
	}
	files2, err := CLI(&doc2, GenCLIWithFramework(YargsFramework))
	if err != nil {
		t.Fatalf("unexpected error generating yargs output (no required alt-sources): %v", err)
	}
	var all2 strings.Builder
	for _, content := range files2 {
		all2.Write(content)
		all2.WriteByte('\n')
	}
	if strings.Contains(all2.String(), "_missing") {
		t.Error("validation block must be omitted when no flag is both required and alt-sourced")
	}
}

// TestCLI_YargsRequiredAltSourcesGlobalValidation covers the symmetric case for global flags:
// a global flag that is required + alt-sourced must be post-resolution validated, and its
// check must share the same _missing array as any required+alt command flags so the user
// receives one error listing everything absent.
func TestCLI_YargsRequiredAltSourcesGlobalValidation(t *testing.T) {
	makeDoc := func(globalFlags []spec.FlagItem, cmdFlags []spec.FlagItem) *spec.Document {
		return &spec.Document{
			OpenCLIVersion: spec.SchemaVersion,
			Info:           spec.Info{Title: "ReqAlt Global CLI", Binary: "reqaltg"},
			Global:         &spec.Global{Flags: globalFlags},
			Commands: &spec.CommandItem{
				Segment: "deploy",
				Flags:   cmdFlags,
			},
		}
	}

	t.Run("required global alt-sourced flag is validated post-resolution", func(t *testing.T) {
		files, err := CLI(makeDoc(
			[]spec.FlagItem{
				{Name: "org", Type: "string", Required: true, AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "CLI_ORG"}}},
				{Name: "region", Type: "string"}, // optional, no alt-sources
			},
			[]spec.FlagItem{{Name: "env", Type: "string"}},
		), GenCLIWithFramework(YargsFramework))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var buf strings.Builder
		for _, c := range files {
			buf.Write(c)
			buf.WriteByte('\n')
		}
		got := buf.String()

		if !strings.Contains(got, "const _gFlags: GlobalFlags = {") {
			t.Error("expected global flags to be extracted into _gFlags for post-resolution check")
		}
		if !strings.Contains(got, `if (_gFlags.org === undefined) _missing.push("--org");`) {
			t.Error("expected post-resolution validation for required alt-sourced global flag 'org'")
		}
		if strings.Contains(got, `_missing.push("--region")`) {
			t.Error("optional global flag 'region' must not be validated as required")
		}
		if strings.Contains(got, `_missing.push("--env")`) {
			t.Error("optional command flag 'env' must not be validated as required")
		}
		if !strings.Contains(got, "setGlobalFlags(_gFlags)") {
			t.Error("expected setGlobalFlags to use pre-resolved _gFlags variable")
		}
		if strings.Contains(got, "setGlobalFlags({") {
			t.Error("inline setGlobalFlags must not appear when required alt-sourced global flags are present")
		}
	})

	t.Run("optional-only alt-sourced global flag does not produce _gFlags", func(t *testing.T) {
		files, err := CLI(makeDoc(
			[]spec.FlagItem{
				{Name: "org", Type: "string", AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "CLI_ORG"}}},
			},
			nil,
		), GenCLIWithFramework(YargsFramework))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var buf strings.Builder
		for _, c := range files {
			buf.Write(c)
		}
		got := buf.String()

		if strings.Contains(got, "_gFlags") {
			t.Error("_gFlags must not appear when no global flag is both required and alt-sourced")
		}
		if !strings.Contains(got, "setGlobalFlags({") {
			t.Error("expected inline setGlobalFlags when no global required alt-sources")
		}
	})

	t.Run("required+alt command flag and required+alt global flag share one _missing array", func(t *testing.T) {
		files, err := CLI(makeDoc(
			[]spec.FlagItem{
				{Name: "org", Type: "string", Required: true, AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "CLI_ORG"}}},
			},
			[]spec.FlagItem{
				{Name: "username", Type: "string", Required: true, AltSources: []spec.AlternativeSource{{Type: "$ENV", Property: "CLI_USER"}}},
			},
		), GenCLIWithFramework(YargsFramework))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		var buf strings.Builder
		for _, c := range files {
			buf.Write(c)
		}
		got := buf.String()

		if n := strings.Count(got, "const _missing: string[] = [];"); n != 1 {
			t.Errorf("expected exactly one _missing declaration, got %d", n)
		}
		if !strings.Contains(got, `if (cmdFlags.username === undefined) _missing.push("--username");`) {
			t.Error("expected command flag check in shared _missing block")
		}
		if !strings.Contains(got, `if (_gFlags.org === undefined) _missing.push("--org");`) {
			t.Error("expected global flag check in shared _missing block")
		}
		if n := strings.Count(got, "if (_missing.length > 0) {"); n != 1 {
			t.Errorf("expected exactly one _missing check, got %d", n)
		}
	})
}
