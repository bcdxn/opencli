package gen

import (
	"bytes"
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcdxn/opencli/codec"
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
