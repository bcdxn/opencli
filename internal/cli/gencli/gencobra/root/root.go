package root

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/internal/cli/gencli/gencobra/check"
	"github.com/bcdxn/opencli/internal/cli/gencli/gencobra/gen"
	"github.com/bcdxn/opencli/spec"
	"github.com/spf13/cobra"
)

func NewCmdRoot(a gencli.ActionsInterface) (*cobra.Command, error) {
	// Command def
	command := &cobra.Command{
		Use:   "ocli",
		Short: heredoc.Docf("OpenCLI CLI"),
		Long:  heredoc.Docf("ocli is a command line interface designed to work with OpenCLI Spec documents"),
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return gencli.BadUserInput(c, "subcommand is required")
		},
	}
	// Command config
	command.SilenceErrors = true
	command.SilenceUsage = true
	// Flag bindings
	command.PersistentFlags().Bool("help", false, "Show help for command")
	command.Flags().Bool("version", false, "Show ocli version")
	// Sub commands
	command.AddCommand(gen.NewCmdGen(a))
	command.AddCommand(check.NewCmdCheck(a))
	command.AddCommand(&cobra.Command{
		Use:    "completion",
		Hidden: true, // hide the 'auto completion' built in command from the help menu
	})
	// Help
	command.SetHelpFunc(func(_ *cobra.Command, args []string) {
		a.HelpFunc(getRootSpecCmd())
	})
	// Usage
	command.SetUsageFunc(func(_ *cobra.Command) error {
		return a.UsageFunc(getRootSpecCmd())
	})

	return command, nil
}

func getRootSpecCmd() *spec.CommandItem {
	return &spec.CommandItem{
		Segment:         "ocli",
		CommandLine:     "ocli {commands} <arguments> [flags]",
		Summary:         "A CLI for working with OpenCLI Specs",
		Description:     "`ocli` is a command line interface designed to make working with [OpenCLI Spec documents](https://github.com/bcdxn/opencli/tree/main) easier. It provides a number of capabilities, including:\n- validating OpenCLI Spec documents\n- generating Documentation from OpenCLI Spec documents\n\nThe commands are documented below. You can also find out more about each command using the contextual `--help` flag. e.g.:\n\n```sh\n$ ocli gen --help\n```",
		VisibleChildren: true,
		VisibleArgs:     false,
		VisibleFlags:    true,
		Flags: []spec.FlagItem{
			{
				Name:    "help",
				Summary: "Show contextual help menu",
			},
		},
		Commands: []*spec.CommandItem{
			{
				Segment:     "check",
				CommandLine: "ocli check <arguments> [flags]",
				Summary:     "Check a given document for validity",
			},
			{
				Segment:     "gen",
				CommandLine: "ocli gen {commands} <arguments> [flags]",
				Summary:     "Commands used to generate code/docs from an OpenCLI Spec document",
			},
		},
	}
}
