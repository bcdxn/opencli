// generated code
package cli

import (
	"context"

	"github.com/spf13/cobra"
)

type CLIHandlersImplementation interface {
	GhRoot(cmd *cobra.Command, args []string) error
	// alias grouping
	GhAliasRoot(cmd *cobra.Command, args []string) error
	GhAliasDelete(cmd *cobra.Command, args []string) error
	// gist grouping
	GhGistRoot(cmd *cobra.Command, args []string) error
	GhGistClone(cmd *cobra.Command, args []string) error
	GhGistCreate(tx context.Context, args GhGistCreateArgs, flags GhGistCreateFlags) error
}

func FromCobra(handlers CLIHandlersImplementation) *cobra.Command {
	return cmdRoot(handlers)
}

func cmdRoot(handlers CLIHandlersImplementation) *cobra.Command {
	root := cobra.Command{
		Use:   "gh",
		Short: "Work seamlessly with GitHub from the command line.",
		RunE:  handlers.GhRoot,
	}

	root.AddCommand(cmdAliasRoot(handlers))
	root.AddCommand(cmdGistRoot(handlers))

	return &root
}

func cmdAliasRoot(handlers CLIHandlersImplementation) *cobra.Command {
	cmd := cobra.Command{
		Use:   "alias {command} <arguments> [flags]",
		Short: "Aliases can be used to make shortcuts for gh commands or to compose multiple commands.",
		RunE:  handlers.GhAliasRoot,
	}

	cmd.AddCommand(cmdAliasDelete(handlers))

	return &cmd
}

func cmdAliasDelete(handlers CLIHandlersImplementation) *cobra.Command {
	cmd := cobra.Command{
		Use:   "delete (<alias> | --all) [flags]",
		Short: "Delete set aliases",
		RunE:  handlers.GhAliasDelete,
	}

	return &cmd
}

func cmdGistRoot(handlers CLIHandlersImplementation) *cobra.Command {
	cmd := cobra.Command{
		Use:   "gist {command} <arguments> [flags]",
		Short: "Work with GitHub gists.",
		RunE:  handlers.GhGistRoot,
	}

	cmd.AddCommand(cmdGistClone(handlers))
	cmd.AddCommand(cmdGistCreate(handlers))

	return &cmd
}

func cmdGistClone(handlers CLIHandlersImplementation) *cobra.Command {
	cmd := cobra.Command{
		Use:   "clone <gist> [<directory>] [-- <gitflags>...]",
		Short: "Clone a GitHub gist locally",
		RunE:  handlers.GhGistClone,
	}

	return &cmd
}

// named positional arguments
type GhGistCreateArgs struct {
	Filename []string
}

// flags
type GhGistCreateFlags struct {
	Desc   string
	Public bool
}

func cmdGistCreate(handlers CLIHandlersImplementation) *cobra.Command {
	var desc string
	var public bool

	cmd := cobra.Command{
		Use:   "create [<filename>... | -] [flags]",
		Short: "Create a new GitHub gist with given contents.",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := GhGistCreateFlags{
				Desc:   desc,
				Public: public,
			}
			return handlers.GhGistCreate(cmd.Context(), GhGistCreateArgs{}, flags)
		},
	}

	cmd.Flags().StringVarP(&desc, "desc", "d", "", "A description for this gist")
	cmd.Flags().BoolVarP(&public, "public", "p", false, `List the gist publicly (default "secret")`)

	return &cmd
}
