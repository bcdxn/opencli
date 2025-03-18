package ocli

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

// GenDocs generates documentation for the given OpenCLI document in the specified format.
// It accepts the OpenCLI document domain object and the desired output format for the documentation as parameters.
// It returns the generated documentation as a byte slice.
func GenDocs(doc OpenCliDocument, format string) []byte {
	tmpl := getTemplate("markdown")
	buf := bytes.NewBuffer([]byte{})

	// Note that various docs formats may have multiple template files depending on the complexity.
	// `docs.tmpl`, however, always serves as the entrypoint.
	err := tmpl.ExecuteTemplate(buf, "docs.tmpl", doc)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

/* Private functions
------------------------------------------------------------------------------------------------- */

//go:embed templates/docs/*
var docsTemplates embed.FS

func getTemplate(format string) *template.Template {
	// Templates for a specific format are stored in a subdirectory with the name of the format nested within `templates/docs/`.
	t, err := template.New("tmpl").ParseFS(
		docsTemplates,
		fmt.Sprintf("templates/docs/%s/*", format),
	)
	if err != nil {
		panic(err)
	}

	return t
}
