package oclicode

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"text/template"

	"github.com/bcdxn/opencli/oclispec"
)

// GenCliOptions is a functional option to configure the GenCLI function
type GenCliOptions func(*genCliOptions)

// Generate generates CLI boilerplate for the given OpenCLI document using the specified framework.
func Generate(doc oclispec.Document, options ...GenCliOptions) ([]GenFile, error) {
	opts := &genCliOptions{}

	for _, opt := range options {
		opt(opts)
	}

	tmpl := getCliTemplate(opts.Framework)

	// Before generating the code, ensure the root-level command has a summary/description
	if doc.CommandTrie.Root.Command.Summary == "" {
		doc.CommandTrie.Root.Command.Summary = doc.Info.Summary
	}
	if doc.CommandTrie.Root.Command.Description == "" {
		doc.CommandTrie.Root.Command.Description = doc.Info.Description
	}

	switch opts.Framework {
	case "urfavecli":
		return genUrfaveCli(tmpl, cliTmplData{*opts, doc})
	case "yargs":
		return genYargsCli(tmpl, cliTmplData{*opts, doc})
	}

	return nil, errors.New("unsupported framework")
}

func GoPackage(pkg string) GenCliOptions {
	return func(opts *genCliOptions) {
		opts.GoPackage = pkg
	}
}

func Framework(name string) GenCliOptions {
	return func(opts *genCliOptions) {
		opts.Framework = name
	}
}

func ModuleType(moduleType string) GenCliOptions {
	return func(opts *genCliOptions) {
		opts.Framework = moduleType
	}
}

type GenFile struct {
	Name     string
	Contents []byte
}

/* Private functions and types
------------------------------------------------------------------------------------------------- */

// genDocsOptions represents the configurable options when generating documentation.
// The options are meant to be configured using the functional options pattern.
type genCliOptions struct {
	Framework  string
	GoPackage  string
	ModuleType string
}

// cliTmplData represents the data used to render the documentation template.
type cliTmplData struct {
	Opts genCliOptions
	Doc  oclispec.Document
}

//go:embed templates/*
var cliTemplates embed.FS

// getCliTemplate reads the framework-appropariate template file(s) into memory -- a prerequisite for generating the boilerplate code.
func getCliTemplate(framework string) *template.Template {
	// Templates for a specific framework are stored in a subdirectory with the name of the framework nested within `templates/cli/`.
	t, err := template.New("tmpl").Funcs(funcmap()).ParseFS(
		cliTemplates,
		fmt.Sprintf("templates/%s/*", framework),
	)
	if err != nil {
		panic(err)
	}

	return t
}

func genUrfaveCli(tmpl *template.Template, data cliTmplData) ([]GenFile, error) {
	cliInterfaceGenContents := bytes.NewBuffer([]byte{})
	// `cli_interface.gen.go.tmpl` defines the interface that must be implemented to handle all of the CLI command actions.
	err := tmpl.ExecuteTemplate(cliInterfaceGenContents, "cli_interface.gen.go.tmpl", data)
	if err != nil {
		return nil, err
	}

	cliParamsGenContents := bytes.NewBuffer([]byte{})
	// `cli_params.gen.go.tmpl` defines all of the injected argument/flag types.
	err = tmpl.ExecuteTemplate(cliParamsGenContents, "cli_params.gen.go.tmpl", data)
	if err != nil {
		return nil, err
	}

	cliExitCodesGenContents := bytes.NewBuffer([]byte{})
	// `cli_params.gen.go.tmpl` defines all of the injected argument/flag types.
	err = tmpl.ExecuteTemplate(cliExitCodesGenContents, "cli_exitcodes.gen.go.tmpl", data)
	if err != nil {
		return nil, err
	}

	cliGenContents := bytes.NewBuffer([]byte{})
	// `cli.gen.go.tmpl` defines the constructor/entrypoint to the CLI program.
	err = tmpl.ExecuteTemplate(cliGenContents, "cli.gen.go.tmpl", data)
	if err != nil {
		return nil, err
	}

	return []GenFile{
		{
			Name:     "cli_interface.gen.go",
			Contents: cliInterfaceGenContents.Bytes(),
		},
		{
			Name:     "cli_params.gen.go",
			Contents: cliParamsGenContents.Bytes(),
		},
		{
			Name:     "cli_exitcodes.gen.go",
			Contents: cliExitCodesGenContents.Bytes(),
		},
		{
			Name:     "cli.gen.go",
			Contents: cliGenContents.Bytes(),
		},
	}, nil
}

func genYargsCli(tmpl *template.Template, data cliTmplData) ([]GenFile, error) {
	cliInterfaceGenContents := bytes.NewBuffer([]byte{})
	// `cli_interface.gen.js.tmpl` defines the interface that must be implemented to handle all of the CLI command actions.
	err := tmpl.ExecuteTemplate(cliInterfaceGenContents, "cli_interface.gen.js.tmpl", data)
	if err != nil {
		return nil, err
	}

	cliGenContents := bytes.NewBuffer([]byte{})
	// `cli.gen.go.tmpl` defines the constructor/entrypoint to the CLI program.
	err = tmpl.ExecuteTemplate(cliGenContents, "cli.gen.js.tmpl", data)
	if err != nil {
		return nil, err
	}

	return []GenFile{
		{
			Name:     "cli-interface.gen.js",
			Contents: cliInterfaceGenContents.Bytes(),
		},
		{
			Name:     "cli.gen.js",
			Contents: cliGenContents.Bytes(),
		},
	}, nil
}
