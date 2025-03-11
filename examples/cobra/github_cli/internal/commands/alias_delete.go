package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func cmdAliasDelete() *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete (<alias> | --all) [flags]",
		Short: "Delete set aliases",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("gh alias delete")
		},
	}

	return &cmd
}
