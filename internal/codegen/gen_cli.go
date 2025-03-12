package codegen

import (
	_ "embed"
	"os"
	"text/template"
)

//go:embed templates/cobra.tmpl
var tmpl string

func (g *Generator) GenerateCLI() {
	f, err := os.OpenFile("test.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	err = template.Must(template.New("cli").Parse(tmpl)).Execute(f, TemplateDate{Package: "test"})
	if err != nil {
		panic(err)
	}
}

type TemplateDate struct {
	Package string
	Spec    CLISpec
}

type CLISpec struct {
	Info CLISpecInfo
}

type CLISpecInfo struct {
	Title       string
	Summary     string
	Description string
}
