package main

import (
	"context"
	"fmt"

	"github.com/bcdxn/openclispec/poc/examples/cobra/github_cli/internal/app"
	"github.com/bcdxn/openclispec/poc/examples/cobra/github_cli/internal/cli"
)

func main() {
	handlers := app.Handlers{} // implementation
	a := cli.New(handlers)     // generated wrapper

	a.ExecuteContext(context.Background())
	fmt.Println()
}
