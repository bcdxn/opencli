package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cmdAliasRoot() *cobra.Command {
	cmd := cobra.Command{
		Use:   "alias {command} <arguments> [flags]",
		Short: "Aliases can be used to make shortcuts for gh commands or to compose multiple commands.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("gh alias")
		},
	}

	cmd.AddCommand(cmdAliasDelete())

	return &cmd
}
