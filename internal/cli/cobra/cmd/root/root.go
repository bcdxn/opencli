package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/check"
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/gen"
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/help"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
	"github.com/spf13/cobra"
)

func NewCmdRoot(f *cliutils.Factory) (*cobra.Command, error) {
	var argsDef map[string]string
	useLine := "ocli {command} <arguments> [flags]"
	command := &cobra.Command{
		Use:   useLine,
		Short: heredoc.Docf("OpenCLI CLI"),
		Long:  heredoc.Docf("ocli is a command line interface designed to work with OpenCLI Spec documents"),
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return cliutils.BadUserInput(c, argsDef, useLine, "subcommand is required")
		},
	}

	command.SilenceErrors = true
	command.SilenceUsage = true

	// Flags
	command.PersistentFlags().Bool("help", false, "Show help for command")
	command.Flags().Bool("version", false, "Show ocli version")
	// Sub Commands
	command.AddCommand(gen.NewCmdGen(f))
	command.AddCommand(check.NewCmdCheck(f))
	command.AddCommand(&cobra.Command{
		Use:    "completion",
		Hidden: true, // hide the 'auto completion' built in command from the help menu
	})
	//
	w := f.IOStreams.Out
	command.SetHelpFunc(func(c *cobra.Command, args []string) {
		help.HelpFunc(w, command, argsDef, useLine)
	})

	return command, nil
}
