package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/md/markdown.tmpl
var mdTemplate []byte

// genDocsMarkdown executes the template in combination with the spec data + options.
// The result is a fully rendered markdown document.
func genDocsMarkdown(data docsTmplData) ([]byte, error) {
	t, err := template.New("markdown.tmpl").Parse(string(mdTemplate))
	if err != nil {
		return []byte{}, fmt.Errorf("unable to generate markdown docs template: %w", err)
	}

	buf := bytes.NewBuffer([]byte{})

	err = t.ExecuteTemplate(buf, "markdown.tmpl", data)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to render markdown docs: %w", err)
	}

	return buf.Bytes(), nil
}
