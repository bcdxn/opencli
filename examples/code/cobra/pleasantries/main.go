package main

import (
	"context"
	"os"

	"pleasantriescli/internal/cli"
	"pleasantriescli/internal/gencli"
)

// Version is set by goreleaser (ldflags) during build process
var version = "DEV"

func main() {
	actions := cli.NewActions(version)
	code := gencli.Run(context.Background(), actions)
	os.Exit(code)
}
