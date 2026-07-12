// Package gen contains funtions to generate documentation and framework-specific boilerplate code
// from a valid OpenCLI specification document.
//
// Given a [spec.Document], gen can produce:
//
//  1. Documentation — Call [Docs] to generate Markdown (or other formats) suitable
//     for publishing CLI reference documentation.
//
//  2. CLI code — Call [CLI] with a target [CLIFramework] to generate boilerplate
//     code for one of the supported frameworks:
//
//     • [CobraFramework] — github.com/spf13/cobra (Go)
//     • [UrfaveCliFramework] — github.com/urfave/cli/v3 (Go)
//     • [YargsFramework] — yargs (Node.js)
//
// Generated code is returned as a map of relative file paths to their contents,
// ready to be written to disk. All files belong to a single gencli package so no
// cross-package import paths are required.
//
// Both [CLI] and [Docs] support functional options for configuring output format,
// framework, badges, footers, and other generation parameters.
package gen

import (
	"fmt"
	"runtime/debug"

	"github.com/bcdxn/opencli/spec"
)

// CLIFramework is the target framework for CLI code generation.
//
// Use with [GenCLIWithFramework] to select a supported framework:
//
//	files, err := gen.CLI(doc, gen.GenCLIWithFramework(gen.CobraFramework))
type CLIFramework string

const (
	// CobraFramework generates code for the github.com/spf13/cobra framework.
	CobraFramework CLIFramework = "COBRA"
	// YargsFramework generates code for the yargs Node.js framework.
	YargsFramework CLIFramework = "YARGS"
	// UrfaveCliFramework generates code for the github.com/urfave/cli/v3 framework.
	UrfaveCliFramework CLIFramework = "URFAVECLI"
)

// IsValid reports whether f is a supported CLIFramework.
func (f CLIFramework) IsValid() bool {
	return f == CobraFramework || f == YargsFramework || f == UrfaveCliFramework
}

// CLI generates framework-specific boilerplate code from an OpenCLI document.
// It returns a map of relative file paths (under gencli/) to their generated content,
// ready to be written to disk. All files belong to the gencli package.
//
// Example:
//
//	files, err := gen.CLI(doc, gen.GenCLIWithFramework(gen.YargsFramework))
//	if err != nil { return err }
//	for path, content := range files {
//	  fs.WriteFile(path, content, 0644)
//	}
func CLI(doc *spec.Document, options ...GenCLIOption) (map[string][]byte, error) {
	moduleVersion := "unknown version"

	if bi, ok := debug.ReadBuildInfo(); ok {
		if bi.Main.Version != "" {
			moduleVersion = bi.Main.Version
		}
	}

	opts := &genCLIOptions{
		Framework:     CobraFramework,
		ModuleVersion: moduleVersion,
	}
	for _, option := range options {
		option(opts)
	}

	if !opts.Framework.IsValid() {
		return nil, fmt.Errorf("invalid CLI framework: %s", opts.Framework)
	}

	switch opts.Framework {
	case CobraFramework:
		return genCLICobra(doc, opts)
	case YargsFramework:
		return genCLIYargs(doc, opts)
	case UrfaveCliFramework:
		return genCLIUrfaveCli(doc, opts)
	}

	return nil, fmt.Errorf("unsupported CLI framework: %s", opts.Framework)
}

/* Functional Options
------------------------------------------------------------------------------------------------- */

// GenCLIOption is a functional option to configure the CLI function.
type GenCLIOption func(*genCLIOptions)

// GenCLIWithFramework sets the target CLI framework.
func GenCLIWithFramework(f CLIFramework) GenCLIOption {
	return func(opts *genCLIOptions) {
		opts.Framework = f
	}
}

// genCLIOptions holds the configuration for CLI generation.
type genCLIOptions struct {
	Framework     CLIFramework
	ModuleVersion string
}
