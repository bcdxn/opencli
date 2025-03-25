package main

import (
	"context"
	"os"

	"github.com/bcdxn/opencli/internal/cli"
)

// Version is set in ldflags during build process
var Version = "DEV"

func main() {
	cmd := cli.New(cli.Impl{}, Version)

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}
}
