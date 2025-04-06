package main

import (
	"context"
	"log"
	"os"

	"github.com/bcdxn/opencli/examples/urfavecli/cli"
)

func main() {
	version := "1.0.0"
	cmd := cli.New(cli.Impl{}, version)

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
