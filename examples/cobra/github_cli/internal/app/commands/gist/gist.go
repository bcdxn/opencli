package gist

import (
	"context"
	"fmt"

	"github.com/bcdxn/openclispec/examples/cobra/github_cli/internal/cli"
	"github.com/spf13/cobra"
)

type GistHandlers struct{}

func (h GistHandlers) GhGistRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("gh gist")
	return nil
}

func (h GistHandlers) GhGistClone(cmd *cobra.Command, args []string) error {
	fmt.Println("gh gist clone")
	return nil
}

func (h GistHandlers) GhGistCreate(
	ctx context.Context,
	args cli.GhGistCreateArgs,
	flags cli.GhGistCreateFlags,
) error {
	fmt.Printf("gh gist create --desc='%s' --public=%t\n", flags.Desc, flags.Public)
	return nil
}
