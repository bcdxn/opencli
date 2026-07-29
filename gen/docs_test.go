package gen

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/bcdxn/opencli/codec"
)

var (
	goldenMarkdown = "testdata/petstore-cli.md"
	goldenManPage  = "testdata/petstore-cli.man"
)

func TestDocs_Markdown(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	actual, err := Docs(doc, DocsWithFormat(Markdown))
	if err != nil {
		t.Fatalf("unexpected error generated documentation: %v", err)
	}

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenMarkdown), 0755); err != nil {
			t.Fatalf("failed to create golden dir for %s: %v", goldenMarkdown, err)
		}
		if err := os.WriteFile(goldenMarkdown, actual, 0644); err != nil {
			t.Fatalf("failed to write golden file %s: %v", goldenMarkdown, err)
		}
	}

	expected, err := os.ReadFile(goldenMarkdown)

	if !bytes.Equal(actual, expected) {
		t.Fatalf("output docs does not match expected docs")
	}
}

func TestDocs_HTMLComponentBundle(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	actual, err := Docs(
		doc,
		DocsWithFormat(HTML_EMBED),
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

func TestDocs_ManPage(t *testing.T) {
	doc, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}

	actual, err := Docs(doc, DocsWithFormat(MAN))
	if err != nil {
		t.Fatalf("unexpected error generated documentation: %v", err)
	}

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenManPage), 0755); err != nil {
			t.Fatalf("failed to create golden dir for %s: %v", goldenManPage, err)
		}
		if err := os.WriteFile(goldenManPage, actual, 0644); err != nil {
			t.Fatalf("failed to write golden file %s: %v", goldenManPage, err)
		}
	}

	expected, err := os.ReadFile(goldenManPage)
	if err != nil {
		t.Fatalf("unexpected error reading golden file: %v", err)
	}

	if !bytes.Equal(actual, expected) {
		t.Fatalf("output docs does not match expected docs")
	}
}
