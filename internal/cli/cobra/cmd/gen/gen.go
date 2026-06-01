package gen

import (
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/gen/docs"
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/help"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
	"github.com/spf13/cobra"
)

func NewCmdGen(f *cliutils.Factory) *cobra.Command {
	var argsDef map[string]string
	useLine := "ocli gen {command} [flags]"

	command := &cobra.Command{
		Use:   "gen",
		Short: "Commands used to generate code/docs from an OpenCLI Spec document",
		RunE: func(c *cobra.Command, args []string) error {
			return cliutils.BadUserInput(c, argsDef, useLine, "subcommand is required")
		},
	}

	command.AddCommand(docs.NewCmdDocs(f))

	command.SilenceErrors = true
	command.SilenceUsage = true

	w := f.IOStreams.Out

	command.SetHelpFunc(func(c *cobra.Command, args []string) {
		help.HelpFunc(w, command, argsDef, useLine)
	})

	return command
}
