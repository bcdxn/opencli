package gen

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcdxn/opencli/codec"
)

// globalFlagsYAML is a minimal spec that declares genuine (non-help/version)
// global flags plus one leaf per action-call shape, so the generated output for
// every framework exercises the context-injection / setGlobalFlags code paths.
//
//go:embed testdata/globalflags-cli.ocs.yaml
var globalFlagsYAML []byte

// TestCLI_GlobalFlags generates the global-flags fixture for each framework and
// compares against goldens under testdata/<framework>/globalflags/. This locks in
// the non-breaking "inject into context / setGlobalFlags" pattern: action method
// signatures carry no global param, while handlers inject via WithGlobalFlags (Go)
// or call setGlobalFlags before invoking the action (Yargs).
func TestCLI_GlobalFlags(t *testing.T) {
	doc, err := codec.UnmarshalYAML(globalFlagsYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling global-flags fixture: %v", err)
	}

	for _, framework := range []CLIFramework{CobraFramework, UrfaveCliFramework, YargsFramework} {
		fwDir := globalFlagsGoldenDir(framework) // cobra | urfavecli | yargs
		t.Run(string(framework), func(t *testing.T) {
			files, err := CLI(doc, GenCLIWithFramework(framework))
			if err != nil {
				t.Fatalf("unexpected error generating %s CLI: %v", framework, err)
			}
			if len(files) == 0 {
				t.Fatal("expected generated files but got none")
			}

			for relPath, content := range files {
				goldenPath := filepath.Join("testdata", fwDir, "globalflags", relPath)

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
		})
	}
}

// globalFlagsGoldenDir maps a framework to its golden-file directory name.
func globalFlagsGoldenDir(f CLIFramework) string {
	switch f {
	case CobraFramework:
		return "cobra"
	case UrfaveCliFramework:
		return "urfavecli"
	default: // YargsFramework
		return "yargs"
	}
}
