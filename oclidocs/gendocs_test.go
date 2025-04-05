package oclidocs

import (
	_ "embed"
	"testing"

	"github.com/bcdxn/opencli/oclispec"
)

//go:embed testdata/ocli-docs.md
var expected string

func TestGenDocsMarkdown(t *testing.T) {
	ocs, err := oclispec.UnmarshalYAML("testdata/cli.ocs.yaml")
	if err != nil {
		t.Fatal(err)
	}

	docs, err := Generate(ocs)
	if err != nil {
		t.Errorf("error generating docs - %v", err)
	}
	actual := string(docs[0].Contents)
	if actual != expected {
		t.Errorf("generated documentation did not match expectation, FOUND:\n[[%s]]\n==========\nEXPECTED:\n[[%s]]\n==========\n", actual, expected)
	}
}
