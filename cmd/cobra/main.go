package main

import (
	"os"

	oclicmd "github.com/bcdxn/opencli/internal/cli/cobra/ocli"
)

// Version is set by goreleaser (ldflags) during build process
var version = "DEV"

func main() {
	code := oclicmd.Main(version)
	os.Exit(code)
}
