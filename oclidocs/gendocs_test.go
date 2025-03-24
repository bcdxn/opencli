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

	actual := Generate(ocs)
	if string(actual) != expected {
		t.Errorf("generated documentation did not match expectation, FOUND:\n[[%s]]\n==========\nEXPECTED:\n[[%s]]\n==========\n", string(actual), expected)
	}
}
