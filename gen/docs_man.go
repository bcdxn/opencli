package gen

import (
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"unicode/utf8"

	"github.com/bcdxn/opencli/spec"
)

//go:embed templates/man/man.tmpl
var manTemplate []byte

// manTemplateFuncs returns the template.FuncMap used by the man page template.
func manTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"roffEscape":      roffEscape,
		"collectExamples": collectExamples,
	}
}

// genDocsMan executes the man page template in combination with the spec data + options.
// The result is a fully rendered roff/troff man page.
func genDocsMan(data docsTmplData) ([]byte, error) {
	t, err := template.New("man.tmpl").Funcs(manTemplateFuncs()).Parse(string(manTemplate))
	if err != nil {
		return []byte{}, fmt.Errorf("unable to parse man page template: %w", err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = t.ExecuteTemplate(buf, "man.tmpl", data)
	if err != nil {
		return []byte{}, fmt.Errorf("unable to render man page: %w", err)
	}

	return buf.Bytes(), nil
}

// codeBlockRe matches fenced code blocks (3+ backticks).
var codeBlockRe = regexp.MustCompile("```+[\\s\\S]*?```+")

// urlRe matches markdown-style inline links [label](http://...).
// After character escaping, dots in URLs become \., so the pattern
// accepts both escaped and unescaped dots.
var urlRe = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)]*[\.\w][^)]*)\)`)

var paragraphRe = regexp.MustCompile(`\n{2,}`)
var hbRe = regexp.MustCompile(`\n[-~]{3,}\n`)

// listBlockRe matches a run of one or more consecutive unordered list items
// (lines that begin with "- " or "* ").
var listBlockRe = regexp.MustCompile(`(?m)(?:^[*\-] [^\n]*\n?)+`)

// blockQuoteRe matches a run of one or more consecutive markdown block-quote
// lines (lines that begin with "> " or ">").
var blockQuoteRe = regexp.MustCompile(`(?m)(?:^>[\t ]?[^\n]*\n?)+`)

// gfmAlertRe matches a GitHub Flavored Markdown alert type marker, e.g. [!TIP].
var gfmAlertRe = regexp.MustCompile(`(?i)^\[!(TIP|NOTE|WARNING|CAUTION|IMPORTANT)\]$`)

// codeBlock holds a placeholder and the original content for restoration.
type codeBlock struct {
	placeholder string
	content     string
}

// extractCodeBlocks replaces fenced code blocks with null-byte placeholders
// so their contents are not processed by subsequent escaping/transformation steps.
func extractCodeBlocks(s string) (string, []codeBlock) {
	var blocks []codeBlock
	idx := 0

	s = codeBlockRe.ReplaceAllStringFunc(s, func(match string) string {
		inner := match
		if i := strings.Index(inner, "```"); i >= 0 {
			inner = inner[i:]
			if j := strings.LastIndex(inner, "```"); j >= 0 {
				inner = inner[:j+3]
			}
		}
		plh := fmt.Sprintf("\x00CB%d\x00", idx)
		blocks = append(blocks, codeBlock{plh, inner})
		idx++
		return plh
	})

	return s, blocks
}

// escapeRoffChars escapes special roff/troff characters while skipping
// code-block placeholders so they remain intact for restoration.
// Only leading '.' at the start of a line is escaped (to prevent accidental
// macro invocation). Literal backslashes are escaped as '\\' (not '\bK').
func escapeRoffChars(s string, blocks []codeBlock) string {
	var sb strings.Builder
	pos := 0
	atLineStart := true
	for pos < len(s) {
		matched := false
		for _, b := range blocks {
			if strings.HasPrefix(s[pos:], b.placeholder) {
				sb.WriteString(b.placeholder)
				pos += len(b.placeholder)
				matched = true
				break
			}
		}
		if matched {
			continue
		}
		r, size := utf8.DecodeRuneInString(s[pos:])
		switch r {
		case '\\':
			sb.WriteString(`\\`)
			atLineStart = false
		case '.':
			if atLineStart {
				sb.WriteString(`\.`)
			} else {
				sb.WriteRune('.')
			}
			atLineStart = false
		case '\n':
			sb.WriteRune('\n')
			atLineStart = true
		default:
			sb.WriteRune(r)
			atLineStart = false
		}
		pos += size
	}
	return sb.String()
}

// transformParagraphs converts runs of 2+ newlines into .PP roff macros.
func transformParagraphs(s string) string {
	return paragraphRe.ReplaceAllString(s, "\n.PP\n")
}

// transformHorizontalBreaks converts markdown horizontal rules (--- or ~~~) into roff lines.
func transformHorizontalBreaks(s string) string {
	return hbRe.ReplaceAllString(s, "\n\\l'20n'\n")
}

// formatGFMAlertLabel returns the title-cased label for a GFM alert type.
func formatGFMAlertLabel(alertType string) string {
	return strings.ToUpper(alertType[:1]) + strings.ToLower(alertType[1:])
}

// transformBlockQuotes converts markdown block-quote lines and GFM alerts to roff.
func transformBlockQuotes(s string) string {
	return blockQuoteRe.ReplaceAllStringFunc(s, func(block string) string {
		lines := strings.Split(strings.TrimRight(block, "\n"), "\n")
		content := make([]string, 0, len(lines))
		for _, line := range lines {
			switch {
			case strings.HasPrefix(line, "> "):
				content = append(content, line[2:])
			case strings.HasPrefix(line, ">"):
				content = append(content, line[1:])
			default:
				content = append(content, line)
			}
		}

		var buf strings.Builder
		buf.WriteString(".PP\n")
		start := 0

		if len(content) > 0 {
			if m := gfmAlertRe.FindStringSubmatch(content[0]); m != nil {
				label := formatGFMAlertLabel(m[1])
				buf.WriteString("\\fB")
				buf.WriteString(label)
				buf.WriteString(":\\fP\n")
				start = 1
			}
		}

		for _, line := range content[start:] {
			buf.WriteString(line)
			buf.WriteByte('\n')
		}
		return buf.String()
	})
}

// transformUnorderedLists converts markdown unordered list items to .IP roff entries.
func transformUnorderedLists(s string) string {
	return listBlockRe.ReplaceAllStringFunc(s, func(block string) string {
		var buf strings.Builder
		for _, line := range strings.Split(strings.TrimRight(block, "\n"), "\n") {
			if len(line) >= 2 {
				buf.WriteString(".IP \\(bu 2\n")
				buf.WriteString(line[2:])
				buf.WriteByte('\n')
			}
		}
		return buf.String()
	})
}

// transformURLs converts markdown-style inline links to .UR /.UE roff macros.
func transformURLs(s string) string {
	var urlBuf strings.Builder
	prev := 0
	for _, loc := range urlRe.FindAllStringIndex(s, -1) {
		start, end := loc[0], loc[1]
		before := s[prev:start]
		urlBuf.WriteString(before)

		caps := urlRe.FindStringSubmatch(s[start:end])
		if len(caps) >= 3 {
			url := strings.ReplaceAll(caps[2], `\.`, ".")
			label := strings.Join(strings.Fields(caps[1]), " ")
			if len(before) > 0 && before[len(before)-1] != '\n' {
				urlBuf.WriteByte('\n')
			}
			fmt.Fprintf(&urlBuf, ".UR %s\n%s\n.UE \\c", url, label)
		} else {
			urlBuf.WriteString(s[start:end])
		}
		prev = end
	}
	urlBuf.WriteString(s[prev:])
	return urlBuf.String()
}

// restoreCodeBlocks replaces code-block placeholders with the original content,
// wrapped in .nf/.fi roff macros for verbatim rendering.
func restoreCodeBlocks(s string, blocks []codeBlock) string {
	var out strings.Builder
	pos := 0
	for pos < len(s) {
		replaced := false
		for _, b := range blocks {
			if strings.HasPrefix(s[pos:], b.placeholder) {
				out.WriteString(".nf\n")
				out.WriteString(b.content)
				out.WriteString("\n.fi")
				pos += len(b.placeholder)
				replaced = true
				break
			}
		}
		if replaced {
			continue
		}
		r, size := utf8.DecodeRuneInString(s[pos:])
		out.WriteRune(r)
		pos += size
	}
	return out.String()
}

// roffEscape escapes special roff/troff characters and converts certain
// markdown constructs to their roff equivalents.
// - Fenced code blocks are wrapped with .nf / .fi.
// - Markdown links [label](url) become .UR /.UE macros.
func roffEscape(s string) string {
	s, blocks := extractCodeBlocks(s)
	s = escapeRoffChars(s, blocks)
	s = transformParagraphs(s)
	s = transformHorizontalBreaks(s)
	s = transformBlockQuotes(s)
	s = transformUnorderedLists(s)
	s = transformURLs(s)

	return restoreCodeBlocks(s, blocks)
}

// collectExamples recursively collects all examples from the command tree.
func collectExamples(cmd *spec.CommandItem) []spec.Example {
	examples := append([]spec.Example{}, cmd.Examples...)
	for _, c := range cmd.Commands {
		if !c.Hidden {
			examples = append(examples, collectExamples(c)...)
		}
	}
	return examples
}
