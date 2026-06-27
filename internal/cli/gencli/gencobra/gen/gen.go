package gen

import (
	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/internal/cli/gencli/gencobra/gen/docs"
	"github.com/bcdxn/opencli/spec"
	"github.com/spf13/cobra"
)

func NewCmdGen(a gencli.ActionsInterface) *cobra.Command {
	command := &cobra.Command{
		Use:   "gen",
		Short: "Commands used to generate code/docs from an OpenCLI Spec document",
		RunE: func(c *cobra.Command, args []string) error {
			return gencli.BadUserInput(c, "subcommand is required")
		},
	}

	command.AddCommand(docs.NewCmdDocs(a))

	command.SilenceErrors = true
	command.SilenceUsage = true

	// Help
	command.SetHelpFunc(func(c *cobra.Command, args []string) {
		a.HelpFunc(&spec.CommandItem{})
	})

	return command
}
