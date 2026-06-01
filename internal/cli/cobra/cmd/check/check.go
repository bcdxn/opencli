package check

import (
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/help"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
	"github.com/spf13/cobra"
)

func NewCmdCheck(f *cliutils.Factory) *cobra.Command {
	var argsDef map[string]string
	useLine := "ocli check <spec-file> [flags]"
	var failOnErr bool

	command := &cobra.Command{
		Use:   "check",
		Short: "Check a given document for validity",
		Long:  "Validate the given file against the OpenCLI specification. This includes static schema checks as well as dynamic checks.",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cliutils.BadUserInput(c, argsDef, useLine, "required arg <spec-file> is missing")
			}

			checkArgs := cliutils.OcliCheckArgs{
				PathToSpec: args[0],
			}

			checkFlags := cliutils.OcliCheckFlags{
				FailOnErr: failOnErr,
			}
			return f.Actions.OcliCheck(checkArgs, checkFlags)
		},
	}

	w := f.IOStreams.Out
	// Arguments
	argsDef = make(map[string]string)
	argsDef["spec-file"] = "path to spec file"

	// Flags
	command.Flags().BoolVar(&failOnErr, "fail-on-err", true, "exit with non-zero code if validation errors are found")

	command.SetHelpFunc(func(c *cobra.Command, args []string) {
		help.HelpFunc(w, command, argsDef, useLine)
	})

	return command
}
