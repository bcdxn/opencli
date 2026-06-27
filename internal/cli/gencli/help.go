package gencli

import (
	"fmt"
	"strings"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
	"github.com/bcdxn/opencli/spec"
)

// func HelpFunc(w io.Writer, cmd *cobra.Command, args map[string]string, useLine string) error {
// 	desc := cmd.Short

// 	if len(cmd.Long) > 0 {
// 		desc = heredoc.Doc(cmd.Long)
// 	}
// 	fmt.Fprintf(w, "%s\n\n", desc)
// 	fmt.Fprint(w, bold("USAGE:"))
// 	fmt.Fprintf(w, "\n  %s", useLine)

// 	var subcommands []*cobra.Command
// 	for _, c := range cmd.Commands() {
// 		if !c.IsAvailableCommand() {
// 			continue
// 		}
// 		subcommands = append(subcommands, c)
// 	}

// 	if len(subcommands) > 0 {
// 		fmt.Fprint(w, bold("\n\nAVAILABLE COMMANDS:"))
// 		col1 := []string{}
// 		col2 := []string{}
// 		for _, c := range subcommands {
// 			col1 = append(col1, c.Name())
// 			col2 = append(col2, c.Short)
// 		}
// 		fmt.Fprint(w, columns(col1, col2, ": "))
// 	}

// 	if len(args) > 0 {
// 		fmt.Fprint(w, bold("\n\nARGUMENTS:"))
// 		col1 := []string{}
// 		col2 := []string{}
// 		for name, desc := range args {
// 			col1 = append(col1, fmt.Sprintf("<%s>", name))
// 			col2 = append(col2, desc)
// 		}
// 		fmt.Fprint(w, columns(col1, col2, " "))
// 	}

// 	flagUsages := cmd.LocalFlags().FlagUsages()
// 	if flagUsages != "" {
// 		fmt.Fprintln(w, bold("\n\nFLAGS:"))
// 		fmt.Fprint(w, flagUsages)
// 	} else {
// 		fmt.Fprint(w, "\n")
// 	}

// 	return nil
// }

// func UsageFunc(w io.Writer, cmd *cobra.Command, args map[string]string, useLine string) error {
// 	fmt.Fprint(w, bold("\nUSAGE:"))
// 	fmt.Fprintf(w, "\n  %s", useLine)

// 	var subcommands []*cobra.Command
// 	for _, c := range cmd.Commands() {
// 		if !c.IsAvailableCommand() {
// 			continue
// 		}
// 		subcommands = append(subcommands, c)
// 	}

// 	if len(subcommands) > 0 {
// 		fmt.Fprint(w, bold("\n\nAVAILABLE COMMANDS:"))
// 		col1 := []string{}
// 		col2 := []string{}
// 		for _, c := range subcommands {
// 			col1 = append(col1, c.Name())
// 			col2 = append(col2, c.Short)
// 		}
// 		fmt.Fprint(w, columns(col1, col2, ": "))
// 	}

// 	if len(args) > 0 {
// 		fmt.Fprint(w, bold("\n\nARGUMENTS:"))
// 		col1 := []string{}
// 		col2 := []string{}
// 		for name, desc := range args {
// 			col1 = append(col1, fmt.Sprintf("<%s>", name))
// 			col2 = append(col2, desc)
// 		}
// 		fmt.Fprint(w, columns(col1, col2, " "))
// 	}

// 	flagUsages := cmd.LocalFlags().FlagUsages()
// 	if flagUsages != "" {
// 		fmt.Fprintln(w, bold("\n\nFLAGS:"))
// 		fmt.Fprint(w, flagUsages)
// 	} else {
// 		fmt.Fprint(w, "\n")
// 	}

// 	fmt.Fprint(w, "\n\n")

// 	return nil
// }

var mdFormatting = []byte(`{
	"document": {
		"block_prefix": "\n",
		"block_suffix": "\n",
		"margin": 0
	},
	"heading": {
		"block_suffix": "\n",
		"bold": true
	},
	"h1": {
		"prefix": " ",
		"suffix": " ",
		"bold": true
	},
	"h2": {
		"prefix": "## ",
		"bold": true
	},
	"emph": {
		"underline": true
	},
	"strong": {
		"bold": true
	},
	"link": {
		"underline": true
	},
	"code_block": {
		"margin": 2
	},
	"list": {
		"indent": 2
	},
	"item": {
		"prefix": "• "
	}
}`)

var lightTheme = []byte(`{
	"document": {
		"color": "236"
	},
	"heading": {
		"color": "239"
	},
	"h1": {
		"color": "236",
		"background_color": "252"
	},
	"h2": {
		"color": "238"
	},
	"link": {
		"color": "31"
	},
	"code": {
		"color": "167",
		"background_color": "254"
	},
	"code_block": {
		"color": "244"
	}
}`)

var darkTheme = []byte(`{
	"document": {
		"color": "251"
	},
	"heading": {
		"color": "250"
	},
	"h1": {
		"color": "252",
		"background_color": "238"
	},
	"h2": {
		"color": "250"
	},
	"link": {
		"color": "110"
	},
	"code": {
		"color": "180",
		"background_color": "237"
	},
	"code_block": {
		"color": "246"
	}
}`)

var noPadding = []byte(`{
	"document": {
		"block_prefix": "",
		"block_suffix": "",
		"margin": 0
	}
}`)

var bold = lipgloss.NewStyle().Bold(true)

func markdownTheme(a ActionsInterface) []byte {
	if a.IOStreams().TerminalTheme() == "dark" {
		return darkTheme
	}

	return lightTheme
}

func DefaultHelpFunc(a ActionsInterface, cmd *spec.CommandItem) {
	stdout := a.IOStreams().Out()
	r, _ := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(mdFormatting),
		glamour.WithStylesFromJSONBytes(markdownTheme(a)),
		// wrap output at specific width (default is 80)
		glamour.WithWordWrap(80),
	)

	desc := []string{cmd.Summary}
	if cmd.Description != "" {
		desc = append(desc, cmd.Description)
	}

	formattedDesc, err := r.Render(strings.Join(desc, "\n\n"))
	if err != nil {
		panic(err)
	}
	// Description
	lipgloss.Fprint(stdout, formattedDesc)
	// Usage
	lipgloss.Fprint(stdout, bold.Render("USAGE:"))
	lipgloss.Fprint(stdout, useLine(cmd))

	if cmd.VisibleChildren {
		lipgloss.Fprintf(stdout, "\n%s\n", bold.Render("AVAILABLE COMMANDS"))
		lipgloss.Fprint(stdout, availableCommands(a, cmd))
	}

	if cmd.VisibleArgs {
		lipgloss.Fprintf(stdout, "\n%s\n", bold.Render("ARGUMENTS"))
		lipgloss.Fprint(a.IOStreams().Out(), availableArgs(a, cmd))
	}

	if cmd.VisibleFlags {
		lipgloss.Fprintf(stdout, "\n\n%s\n", bold.Render("FLAGS"))
		lipgloss.Fprint(stdout, availableFlags(a, cmd))
	}
}

func DefaultUsageFunc(a ActionsInterface, cmd *spec.CommandItem) error {
	stdout := a.IOStreams().Out()
	// Usage
	lipgloss.Fprint(stdout, bold.Render("USAGE:"))
	lipgloss.Fprint(stdout, useLine(cmd))

	if cmd.VisibleChildren {
		lipgloss.Fprintf(stdout, "\n%s\n", bold.Render("AVAILABLE COMMANDS"))
		lipgloss.Fprint(stdout, availableCommands(a, cmd))
	}

	if cmd.VisibleArgs {
		lipgloss.Fprintf(stdout, "\n\n%s\n", bold.Render("ARGUMENTS"))
		lipgloss.Fprint(a.IOStreams().Out(), availableArgs(a, cmd))
	}

	if cmd.VisibleFlags {
		lipgloss.Fprintf(stdout, "\n\n%s\n", bold.Render("FLAGS"))
		lipgloss.Fprint(stdout, availableFlags(a, cmd))
		lipgloss.Fprint(stdout, "\n")
	}

	return nil
}

func useLine(cmd *spec.CommandItem) string {
	line := []string{cmd.CommandLine}

	if len(cmd.CommandModifiers) > 0 {
		line = append(line, strings.Join(cmd.CommandModifiers, " "))
	}
	if len(cmd.ArgsModifiers) > 0 {
		line = append(line, strings.Join(cmd.ArgsModifiers, " "))
	}
	if len(cmd.FlagsModifiers) > 0 {
		line = append(line, strings.Join(cmd.FlagsModifiers, " "))
	}

	return fmt.Sprintf("\n  %s\n", strings.Join(line, " "))
}

func availableCommands(a ActionsInterface, cmd *spec.CommandItem) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(mdFormatting),
		glamour.WithStylesFromJSONBytes(markdownTheme(a)),
		glamour.WithStylesFromJSONBytes(noPadding),
		glamour.WithWordWrap(40),
	)
	if err != nil {
		panic(err)
	}

	names := []string{}
	for _, subcmd := range cmd.Commands {
		names = append(names, subcmd.Segment)
	}

	colWidth := columnWidth(names)

	leftStyle := lipgloss.NewStyle().
		Width(colWidth + 3). // includes the width of padding
		PaddingRight(1).
		PaddingLeft(2)

	rightStyle := lipgloss.NewStyle().
		Width(40)

	rows := []string{}
	for _, subcmd := range cmd.Commands {
		formattedName := leftStyle.Render(subcmd.Segment)
		formattedSummary, err := r.Render(subcmd.Summary)
		if err != nil {
			panic(err)
		}

		formattedSummary = strings.TrimSuffix(formattedSummary, "\n")
		formattedSummary = rightStyle.Render(formattedSummary)
		rows = append(rows, lipgloss.JoinHorizontal(
			lipgloss.Top,
			formattedName,
			formattedSummary,
		))
	}

	ret := lipgloss.JoinVertical(lipgloss.Top, rows...)

	return ret
}

func availableArgs(a ActionsInterface, cmd *spec.CommandItem) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(mdFormatting),
		glamour.WithStylesFromJSONBytes(markdownTheme(a)),
		glamour.WithStylesFromJSONBytes(noPadding),
		glamour.WithWordWrap(40),
	)
	if err != nil {
		panic(err)
	}

	names := []string{}
	for _, arg := range cmd.Args {
		names = append(names, arg.Name)
	}

	colWidth := columnWidth(names)

	leftStyle := lipgloss.NewStyle().
		Width(colWidth + 5). // includes the width of padding + angle brackets styling, < >
		PaddingRight(1).
		PaddingLeft(2)

	rightStyle := lipgloss.NewStyle().
		Width(40)

	rows := []string{}
	for _, arg := range cmd.Args {
		formattedName := leftStyle.Render(fmt.Sprintf("<%>", arg.Name))
		formattedSummary, err := r.Render(arg.Summary)
		if err != nil {
			panic(err)
		}

		formattedSummary = strings.TrimSuffix(formattedSummary, "\n")
		formattedSummary = rightStyle.Render(formattedSummary)
		rows = append(rows, lipgloss.JoinHorizontal(
			lipgloss.Top,
			formattedName,
			formattedSummary,
		))
	}

	ret := lipgloss.JoinVertical(lipgloss.Top, rows...)

	return ret
}

func availableFlags(a ActionsInterface, cmd *spec.CommandItem) string {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(mdFormatting),
		glamour.WithStylesFromJSONBytes(markdownTheme(a)),
		glamour.WithStylesFromJSONBytes(noPadding),
		glamour.WithWordWrap(40),
	)
	if err != nil {
		panic(err)
	}

	names := []string{}
	for _, flag := range cmd.Flags {
		names = append(names, flagNameWithAliases(flag))
	}

	colWidth := columnWidth(names)

	leftStyle := lipgloss.NewStyle().
		Width(colWidth + 3). // includes the width of padding
		PaddingRight(1).
		PaddingLeft(2)

	rightStyle := lipgloss.NewStyle().
		Width(40)

	rows := []string{}
	for _, flag := range cmd.Flags {
		formattedName := leftStyle.Render(flagNameWithAliases(flag))
		formattedSummary, err := r.Render(flag.Summary)
		if err != nil {
			panic(err)
		}

		formattedSummary = strings.TrimSuffix(formattedSummary, "\n")
		formattedSummary = rightStyle.Render(formattedSummary)
		rows = append(rows, lipgloss.JoinHorizontal(
			lipgloss.Top,
			formattedName,
			formattedSummary,
		))
	}

	ret := lipgloss.JoinVertical(lipgloss.Top, rows...)

	return ret
}

func flagNameWithAliases(flag spec.FlagItem) string {
	flagWithAliases := []string{fmt.Sprintf("--%s", flag.Name)}
	for _, alias := range flag.Aliases {
		flagWithAliases = append(flagWithAliases, fmt.Sprintf("-%s", alias))
	}
	return strings.Join(flagWithAliases, " ")
}

func columnWidth(rows []string) int {
	max := 0

	for _, row := range rows {
		if len(row) > max {
			max = len(row)
		}
	}

	return max
}
