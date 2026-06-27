package gencobra

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bcdxn/opencli/internal/cli/gencli"
	"github.com/bcdxn/opencli/internal/cli/gencli/gencobra/root"
)

func Run(actions gencli.ActionsInterface, ctx context.Context) int {
	rootCmd, err := root.NewCmdRoot(actions)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create root command: %s\n", err)
		return gencli.ExitCodeInternalErr
	}

	if _, err := rootCmd.ExecuteContextC(ctx); err != nil {
		fmt.Fprintf(actions.IOStreams().Out(), "error: %v\n\n", err.Error())

		if cliErr, ok := errors.AsType[*gencli.CLIError](err); ok {
			cliErr.Command.UsageFunc()(cliErr.Command)
			return cliErr.Code
		}

		return gencli.ExitCodeInternalErr
	}

	return gencli.ExitCodeOK
}

// func getRootSpecCmd() *spec.CommandItem {
// 	return &spec.CommandItem{
// 		Segment:     "ocli",
// 		CommandLine: "ocli {commands} <arguments> [flags]",
// 		Flags: []spec.FlagItem{
// 			{
// 				Name:    "help",
// 				Summary: "Show contextual help menu",
// 			},
// 		},
// 		Commands: []*spec.CommandItem{
// 			{
// 				Segment:     "check",
// 				CommandLine: "ocli check <arguments> [flags]",
// 				Summary:     "Check a given document for validity",
// 			},
// 			{
// 				Segment:     "gen",
// 				CommandLine: "ocli gen {commands} <arguments> [flags]",
// 				Summary:     "Commands used to generate code/docs from an OpenCLI Spec document",
// 			},
// 		},
// 	}
// }
