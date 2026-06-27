package main

import (
	"context"
	"os"

	cli "github.com/bcdxn/opencli/internal/cli/app"
	"github.com/bcdxn/opencli/internal/cli/gencli"
)

// Version is set by goreleaser (ldflags) during build process
var version = "DEV"

func main() {
	actions := cli.NewActions(version)
	code := gencli.Run(context.Background(), actions)
	os.Exit(code)
}
