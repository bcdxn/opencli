package alias

import (
	"context"
	"fmt"

	"github.com/bcdxn/openclispec/poc/examples/cobra/github_cli/internal/cli"
	"github.com/spf13/cobra"
)

type AliasHandlers struct{}

func (h AliasHandlers) GhAliasRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("gh alias")
	return nil
}

func (h AliasHandlers) GhAliasDelete(ctx context.Context, flags cli.GhAliasDeleteFlags) error {
	fmt.Println("gh alias delete")
	return nil
}

func (h AliasHandlers) GhAliasSet(ctx context.Context, flags cli.GhAliasSetFlags) error {
	fmt.Println("gh alias delete")
	return nil
}
