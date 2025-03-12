package main

import "github.com/bcdxn/openclispec/internal/codegen"

func main() {
	g := codegen.Generator{}
	g.GenerateCLI()
}
