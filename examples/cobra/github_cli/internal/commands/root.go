package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) error {
	return Root().ExecuteContext(ctx)
}

func Root() *cobra.Command {
	root := cobra.Command{
		Use:   "gh",
		Short: "Work seamlessly with GitHub from the command line.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("gh ")
		},
	}

	root.AddCommand(cmdAliasRoot())

	return &root
}
