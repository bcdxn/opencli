package oclidocs

import (
	_ "embed"
	"os"
	"testing"

	"github.com/bcdxn/opencli/oclispec"
)

func TestGenDocsMarkdown(t *testing.T) {
	ocs, err := oclispec.UnmarshalYAML("../examples/cli.ocs.yaml")
	if err != nil {
		t.Fatal(err)
	}

	docs, err := Generate(ocs)
	if err != nil {
		t.Errorf("error generating docs - %v", err)
	}

	expectedDocs, err := os.ReadFile("../examples/markdown-docs/docs.gen.md")
	if err != nil {
		t.Fatalf("error loading expected docs file - %v", err)
	}

	actual := string(docs[0].Contents)
	expected := string(expectedDocs)
	if actual != expected {
		t.Errorf("generated documentation did not match expectation, FOUND:\n[[%s]]\n==========\nEXPECTED:\n[[%s]]\n==========\n", actual, expected)
	}
}
