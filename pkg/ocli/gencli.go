package ocli

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

// GenCliOptions is a functional option to configure the GenCLI function
type GenCliOptions func(*genCliOptions)

// GenCLI generates CLI boilerplate for the given OpenCLI document using the specified framework.
func GenCLI(doc OpenCliDocument, options ...GenCliOptions) []byte {
	opts := &genCliOptions{
		Framework: "urfavecli",
	}

	for _, opt := range options {
		opt(opts)
	}

	tmpl := getCliTemplate(opts.Framework)
	buf := bytes.NewBuffer([]byte{})
	// Note that various cli frameworks may have multiple template files depending on the complexity.
	// `cli.gen.go.tmpl`, however, always serves as the entrypoint.
	err := tmpl.ExecuteTemplate(buf, "cli.gen.go.tmpl", cliTmplData{doc})
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func Framework(name string) GenCliOptions {
	return func(opts *genCliOptions) {
		opts.Framework = name
	}
}

/* Private functions and types
------------------------------------------------------------------------------------------------- */

// genDocsOptions represents the configurable options when generating documentation.
// The options are meant to be configured using the functional options pattern.
type genCliOptions struct {
	Framework string
}

// cliTmplData represents the data used to render the documentation template.
type cliTmplData struct {
	Doc OpenCliDocument
}

//go:embed templates/cli/*
var cliTemplates embed.FS

// getCliTemplate reads the framework-appropariate template file(s) into memory -- a prerequisite for generating the boilerplate code.
func getCliTemplate(framework string) *template.Template {
	// Templates for a specific framework are stored in a subdirectory with the name of the framework nested within `templates/cli/`.
	t, err := template.New("tmpl").ParseFS(
		cliTemplates,
		fmt.Sprintf("templates/cli/%s/*", framework),
	)
	if err != nil {
		panic(err)
	}

	return t
}
