package main

import (
	"context"
	"log"
	"os"

	"github.com/bcdxn/opencli/examples/urfavecli/cli"
)

var version = "DEV"

func main() {
	cmd := cli.New(cli.Impl{}, version)

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
