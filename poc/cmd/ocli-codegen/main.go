package main

import (
	"github.com/bcdxn/openclispec/poc/internal/codegen"
)

func main() {
	g := codegen.Generator{}
	g.GenerateCLI()
}
