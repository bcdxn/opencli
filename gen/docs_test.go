package gen_test

import (
	_ "embed"
	"os"
	"testing"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/gen"
)

//go:generate mkdir -p out
//go:generate cp -r ../examples/petstore-cli.ocs.yaml ./out/petstore-cli.ocs.yaml

//go:embed out/petstore-cli.ocs.yaml
var exampleYAML []byte

func TestDocs(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	contents, err := gen.Docs(doc, gen.DocsWithFormat(gen.Markdown))
	if err != nil {
		t.Fatalf("unexpected error generated documentation: %v", err)
	}

	os.WriteFile("out/petstore-cli.md", contents, 0644)
}
