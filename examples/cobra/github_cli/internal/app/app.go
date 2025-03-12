package app

import (
	"fmt"

	"github.com/bcdxn/openclispec/examples/cobra/github_cli/internal/app/commands/alias"
	"github.com/bcdxn/openclispec/examples/cobra/github_cli/internal/app/commands/gist"
	"github.com/spf13/cobra"
)

type Handlers struct {
	alias.AliasHandlers
	gist.GistHandlers
}

func (h Handlers) GhRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("gh")
	return nil
}
