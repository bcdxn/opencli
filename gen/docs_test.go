package gen_test

import (
	"bytes"
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

func TestDocs_Markdown(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	actual, err := gen.Docs(doc, gen.DocsWithFormat(gen.Markdown))
	if err != nil {
		t.Fatalf("unexpected error generated documentation: %v", err)
	}

	expected, err := os.ReadFile("testdata/petstore-cli.md")

	if !bytes.Equal(actual, expected) {
		t.Fatalf("output docs does not match expected docs")
	}
}

func TestDocs_HTMLComponentBundle(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	actual, err := gen.Docs(
		doc,
		gen.DocsWithFormat(gen.HTML),
		gen.DocsWithHTMLFlavor(gen.EmbeddableComponent),
	)
	if err != nil {
		t.Fatalf("unexpected error generated html embed docs: %v", err)
	}

	if !bytes.Contains(actual, []byte("global.OcliDocs = OcliDocs")) {
		t.Fatalf("expected embed HTML flavor output to include OcliDocs initializer")
	}

	if !bytes.Contains(actual, []byte("container.innerHTML = EMBED_MARKUP")) {
		t.Fatalf("expected embed HTML flavor output to include embeddable component markup")
	}
}
