package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	cmd "github.com/bcdxn/opencli/internal/cli/cmd/factory"
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/help"
	"github.com/bcdxn/opencli/internal/cli/cobra/cmd/root"
	cliutils "github.com/bcdxn/opencli/internal/cli/utils"
)

func Main(version string) int {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load executable path: %s\n", err)
		return cliutils.ExitCodeInternalErr
	}
	ioStreams := cliutils.System()
	f := cmd.NewCmdFactory(version, ioStreams, executablePath)

	rootCmd, err := root.NewCmdRoot(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create root command: %s\n", err)
		return cliutils.ExitCodeInternalErr
	}

	ctx := context.Background()

	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		fmt.Fprintf(ioStreams.Out, "%v\n", err.Error())

		if cliErr, ok := errors.AsType[*cliutils.CLIError](err); ok {
			if cliErr.Code == cliutils.ExitCodeInternalErr {
				return cliErr.Code
			}

			help.UsageFunc(ioStreams.Out, cliErr.Command, cliErr.ArgsDef, cliErr.UseLine)
			return cliErr.Code
		}

		return cliutils.ExitCodeInternalErr
	}

	return cliutils.ExitCodeOK
}
