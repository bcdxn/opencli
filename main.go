package main

import (
	"context"
	"os"

	"github.com/bcdxn/opencli/internal/cli"
)

func main() {
	cmd := cli.New(cli.Impl{}, cli.Version)

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}
}
