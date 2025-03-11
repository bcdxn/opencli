package main

import (
	"context"
	"fmt"

	"github.com/bcdxn/openclispec/examples/cobra/github_cli/internal/commands"
)

func main() {
	commands.Execute(context.Background())
	fmt.Println()
}
