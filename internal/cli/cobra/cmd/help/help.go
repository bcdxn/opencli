package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func HelpFunc(w io.Writer, cmd *cobra.Command, args map[string]string, useLine string) error {
	desc := cmd.Short

	if len(cmd.Long) > 0 {
		desc = heredoc.Doc(cmd.Long)
	}
	fmt.Fprintf(w, "%s\n\n", desc)
	fmt.Fprint(w, bold("USAGE:"))
	fmt.Fprintf(w, "\n  %s", useLine)

	var subcommands []*cobra.Command
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() {
			continue
		}
		subcommands = append(subcommands, c)
	}

	if len(subcommands) > 0 {
		fmt.Fprint(w, bold("\n\nAVAILABLE COMMANDS:"))
		col1 := []string{}
		col2 := []string{}
		for _, c := range subcommands {
			col1 = append(col1, c.Name())
			col2 = append(col2, c.Short)
		}
		fmt.Fprint(w, columns(col1, col2, ": "))
	}

	if len(args) > 0 {
		fmt.Fprint(w, bold("\n\nARGUMENTS:"))
		col1 := []string{}
		col2 := []string{}
		for name, desc := range args {
			col1 = append(col1, fmt.Sprintf("<%s>", name))
			col2 = append(col2, desc)
		}
		fmt.Fprint(w, columns(col1, col2, " "))
	}

	flagUsages := cmd.LocalFlags().FlagUsages()
	if flagUsages != "" {
		fmt.Fprintln(w, bold("\n\nFLAGS:"))
		fmt.Fprint(w, flagUsages)
	} else {
		fmt.Fprint(w, "\n")
	}

	return nil
}

func UsageFunc(w io.Writer, cmd *cobra.Command, args map[string]string, useLine string) error {
	fmt.Fprint(w, bold("\nUSAGE:"))
	fmt.Fprintf(w, "\n  %s", useLine)

	var subcommands []*cobra.Command
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() {
			continue
		}
		subcommands = append(subcommands, c)
	}

	if len(subcommands) > 0 {
		fmt.Fprint(w, bold("\n\nAVAILABLE COMMANDS:"))
		col1 := []string{}
		col2 := []string{}
		for _, c := range subcommands {
			col1 = append(col1, c.Name())
			col2 = append(col2, c.Short)
		}
		fmt.Fprint(w, columns(col1, col2, ": "))
	}

	if len(args) > 0 {
		fmt.Fprint(w, bold("\n\nARGUMENTS:"))
		col1 := []string{}
		col2 := []string{}
		for name, desc := range args {
			col1 = append(col1, fmt.Sprintf("<%s>", name))
			col2 = append(col2, desc)
		}
		fmt.Fprint(w, columns(col1, col2, " "))
	}

	flagUsages := cmd.LocalFlags().FlagUsages()
	if flagUsages != "" {
		fmt.Fprintln(w, bold("\n\nFLAGS:"))
		fmt.Fprint(w, flagUsages)
	} else {
		fmt.Fprint(w, "\n")
	}

	fmt.Fprint(w, "\n\n")

	return nil
}

func bold(str string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", str)
}

func columns(col1, col2 []string, delimiter string) string {
	width := columnWidth(col1)

	var sb strings.Builder
	for i := range col1 {
		sb.WriteString(fmt.Sprintf("\n  %s%s%s%s", col1[i], delimiter, pad(col1[i], " ", width), col2[i]))
	}

	return sb.String()
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

func pad(str string, c string, l int) string {
	padLen := l - len(str)

	var sb strings.Builder
	for range padLen {
		sb.WriteString(c)
	}

	return sb.String()
}
