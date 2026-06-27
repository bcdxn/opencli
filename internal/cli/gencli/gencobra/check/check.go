package check

import (
	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/spec"
	"github.com/spf13/cobra"
)

func NewCmdCheck(a gencli.ActionsInterface) *cobra.Command {
	// Flag vars
	var failOnErr bool
	// Command def
	command := &cobra.Command{
		Use:   "check",
		Short: "Check a given document for validity",
		Long:  "Validate the given file against the OpenCLI specification. This includes static schema checks as well as dynamic checks.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return gencli.BadUserInput(c, "required arg <spec-file> is missing")
			}

			checkArgs := gencli.OcliCheckArgs{
				PathToSpec: args[0],
			}

			checkFlags := gencli.OcliCheckFlags{
				FailOnErr: failOnErr,
			}
			return a.OcliCheck(checkArgs, checkFlags)
		},
	}
	// Command config
	command.SilenceErrors = true
	command.SilenceUsage = true
	// Flag bindings
	command.Flags().BoolVar(&failOnErr, "fail-on-err", true, "exit with non-zero code if validation errors are found")
	// Help
	command.SetHelpFunc(func(c *cobra.Command, args []string) {
		a.HelpFunc(&spec.CommandItem{})
	})

	return command
}
