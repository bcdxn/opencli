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
						Summary:       "Create command shortcuts",
						CommandFnName: "cmdGhAlias",
						// HandlerFnName: "GhAliasHandler",
						Subcommands: []CLICommand{
							{
								Name:          "delete {<alias> | --all} [flags]",
								Summary:       "Delete set aliases",
								CommandFnName: "cmdGhAliasDelete",
								HandlerFnName: "GhAliasDelete",
								Subcommands:   []CLICommand{},
								Flags: []CLIFlag{
									{
										Name:         "all",
										FlagPropName: "All",
										Summary:      "Delete all aliases",
										Type:         "bool",
										DefaultBool:  false,
									},
								},
								Executable: true,
							},
							// {
							// 	Name:          "import [<filename> | -] [flags]",
							// 	Summary:       "Import aliases from the contents of a YAML file.",
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
							{
								Name:          "set <alias> <expansion> [flags]",
								Summary:       "Define a word that will expand to a full gh command when invoked.",
								CommandFnName: "cmdGhAliasSet",
								HandlerFnName: "GhAliasSet",
								Subcommands:   []CLICommand{},
								Flags: []CLIFlag{
									{
										Name:         "clobber",
										FlagPropName: "Clobber",
										Summary:      "Overwrite existing aliases of the same name",
										Type:         "bool",
										DefaultBool:  false,
									},
									{
										Name:         "shell",
										Alias:        "s",
										FlagPropName: "Shell",
										Summary:      "Declare an alias to be passed through a shell interpreter",
										Type:         "bool",
										DefaultBool:  false,
									},
								},
								Executable: true,
							},
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
	Name            string
	Summary         string
	Description     string
	CommandFnName   string
	HandlerFnName   string
	FlagsStructName string
	Subcommands     []CLICommand
	Flags           []CLIFlag
	Executable      bool
}

type CLIFlag struct {
	FlagPropName  string
	Name          string
	Alias         string
	Summary       string
	Type          string
	DefaultString string
	DefaultBool   bool
}
