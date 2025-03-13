package codegen

import (
	"embed"
	"os"
	"text/template"
)

//go:embed all:templates/cobra
var cobraTemplates embed.FS

func (g *Generator) GenerateCLI() {
	f, err := os.OpenFile("cli.gen.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	testCLI := TemplateData{
		Package: "cli",
		Commands: []CLICommand{
			{
				Name:          "gh",
				Summary:       "Work seamlessly with GitHub from the command line.",
				CommandFnName: "cmdGh",
				// HandlerFnName: "GhHandler",
				Subcommands: []CLICommand{
					{
						Name:          "alias",
						Summary:       "",
						CommandFnName: "cmdGhAlias",
						// HandlerFnName: "GhAliasHandler",
						Subcommands: []CLICommand{
							{
								Name:          "delete {<alias> | --all} [flags]",
								Summary:       "Delete set aliases",
								CommandFnName: "cmdGhAliasDelete",
								HandlerFnName: "GhAliasDelete",
								Subcommands:   []CLICommand{},
								Executable:    true,
							},
							// {
							// 	Name:          "import",
							// 	Summary:       "",
							// 	CommandFnName: "cmdGhAliasImport",
							// 	HandlerFnName: "GhAliasImport",
							// 	Subcommands:   []CLICommand{},
							// 	Executable:    true,
							// },
							// {
							// 	Name:          "list",
							// 	Summary:       "",
							// 	CommandFnName: "cmdGhAliasList",
							// 	HandlerFnName: "GhAliasList",
							// 	Subcommands:   []CLICommand{},
							// 	Executable:    true,
							// },
							// {
							// 	Name:          "set",
							// 	Summary:       "",
							// 	CommandFnName: "cmdGhAliasSet",
							// 	HandlerFnName: "GhAliasSet",
							// 	Subcommands:   []CLICommand{},
							// 	Executable:    true,
							// },
						},
					},
					// {
					// 	Name:          "gist",
					// 	Summary:       "",
					// 	CommandFnName: "cmdGhGist",
					// 	HandlerFnName: "GhGistHandler",
					// 	Subcommands:   []CLICommand{},
					// },
				},
			},
		},
	}

	t, err := template.ParseFS(
		cobraTemplates,
		"templates/cobra/*",
	)
	if err != nil {
		panic(err)
	}

	err = t.ExecuteTemplate(f, "cli.gen.go.tmpl", testCLI)
	if err != nil {
		panic(err)
	}
}

type TemplateData struct {
	Package  string
	Commands []CLICommand
}

type CLICommand struct {
	Name          string
	Summary       string
	Description   string
	CommandFnName string
	HandlerFnName string
	Subcommands   []CLICommand
	Executable    bool
}
