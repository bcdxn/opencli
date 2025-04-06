package cli

import (
	"context"
	"fmt"

	urfavecli "github.com/urfave/cli/v3"
)

type Impl struct{}

func (Impl) PleasantriesFarewell(ctx context.Context, cmd *urfavecli.Command, arguments PleasantriesFarewellArgs, flags PleasantriesFarewellFlags) error {
	if flags.Language == "spanish" {
		fmt.Println("adios", arguments.Name)
	} else {
		fmt.Println("bye", arguments.Name)
	}
	return nil
}

func (Impl) PleasantriesGreet(ctx context.Context, cmd *urfavecli.Command, arguments PleasantriesGreetArgs, flags PleasantriesGreetFlags) error {
	if flags.Language == "spanish" {
		fmt.Println("hola", arguments.Name)
	} else {
		fmt.Println("hello", arguments.Name)
	}
	return nil
}
