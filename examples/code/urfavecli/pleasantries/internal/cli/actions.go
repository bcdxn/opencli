package cli

import (
	"context"
	"fmt"
	"pleasantriescli/internal/gencli"

	"github.com/bcdxn/opencli/spec"
)

// Ensure we conform to the generated ActionsInterface
var _ gencli.ActionsInterface = (*Actions)(nil)

func NewActions(version string) *Actions {
	return &Actions{
		IOS:     gencli.DefaultIOS(),
		version: version,
	}
}

// Actions implements the gencli Actions interface and can be passed via the gencli.Factory
type Actions struct {
	IOS     gencli.IOStreams
	version string
}

func (a Actions) PleasantriesGreet(ctx context.Context, args gencli.PleasantriesGreetArgs, flags gencli.PleasantriesGreetFlags) error {
	if flags.Language == gencli.PleasantriesGreetLanguageEnglish {
		fmt.Println("Hello", args.Name)
	} else {
		fmt.Println("Hola", args.Name)
	}
	return nil
}
func (a Actions) PleasantriesFarewell(ctx context.Context, args gencli.PleasantriesFarewellArgs, flags gencli.PleasantriesFarewellFlags) error {
	if flags.Language == gencli.PleasantriesFarewellLanguageEnglish {
		fmt.Println("Good bye", args.Name)
	} else {
		fmt.Println("Adios", args.Name)
	}
	return nil
}

func (a Actions) HelpFunc(cmd *spec.CommandItem) {
	gencli.DefaultHelpFunc(a, cmd)
}
func (a Actions) UsageFunc(cmd *spec.CommandItem) error {
	return gencli.DefaultUsageFunc(a, cmd)
}
func (a Actions) IOStreams() gencli.IOStreams {
	return gencli.DefaultIOS()
}
func (a Actions) Version() string {
	return a.version
}
