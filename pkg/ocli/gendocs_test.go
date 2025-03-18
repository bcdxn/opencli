package ocli

import (
	_ "embed"
	"testing"
)

//go:embed testdata/ocli-docs.md
var expected string

func TestGenDocsMarkdown(t *testing.T) {
	ocs, err := UnmarshalYAML("testdata/cli.ocs.yaml")
	if err != nil {
		t.Fatal(err)
	}

	actual := GenDocs(ocs, "markdown")
	if string(actual) != expected {
		t.Errorf("generated documentation did not match expectation, FOUND:\n%s\n==========\nEXPECTED:\n%s\n==========\n", string(actual), expected)
	}
}
