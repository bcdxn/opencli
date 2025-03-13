package alias

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type AliasHandlers struct{}

func (h AliasHandlers) GhAliasRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("gh alias")
	return nil
}

func (h AliasHandlers) GhAliasDelete(ctx context.Context) error {
	fmt.Println("gh alias delete")
	return nil
}
