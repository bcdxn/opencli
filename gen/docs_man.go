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

// roffEscape escapes special roff/troff characters and converts certain
// markdown constructs to their roff equivalents.
// - Fenced code blocks are wrapped with .nf / .fi.
// - Markdown links [label](url) become .UR /.UE macros.
func roffEscape(s string) string {
	// 1. Extract code blocks so their contents are not escaped or link-processed.
	type block struct{ placeholder, content string }
	var blocks []block
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
		blocks = append(blocks, block{plh, inner})
		idx++
		return plh
	})

	// 2. Character-level escaping (before markdown transformations,
	//    so roff macros inserted later are not re-escaped).
	var sb strings.Builder
	pos := 0
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
			sb.WriteString(`\bK`)
		case '.':
			sb.WriteString(`\.`)
		default:
			sb.WriteRune(r)
		}
		pos += size
	}
	s = sb.String()

	// 3. Convert markdown paragraphs 2+ newlines) into a .PP
	s = paragraphRe.ReplaceAllString(s, "\n.PP\n")

	// 4. Convert orizontal breaks
	s = hbRe.ReplaceAllString(s, "\n\\l'20n'\n")

	// 5. Convert markdown block quotes (">" lines) and GFM alerts to roff.
	//    A GFM alert header ([!TIP], [!NOTE], etc.) becomes a bold label.
	//    All quote content is indented with .RS/.RE; .PP resets the margin
	//    afterwards.
	s = blockQuoteRe.ReplaceAllStringFunc(s, func(block string) string {
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

		// Detect GFM alert header and render as a bold label.
		if len(content) > 0 {
			if m := gfmAlertRe.FindStringSubmatch(content[0]); m != nil {
				label := strings.ToUpper(m[1][:1]) + strings.ToLower(m[1][1:])
				buf.WriteString("\\fB")
				buf.WriteString(label)
				buf.WriteString(":\\fP\n")
				start = 1
			}
		}

		// Wrap remaining content in .RS/.RE for indentation.
		if start < len(content) {
			for _, line := range content[start:] {
				buf.WriteString(line)
				buf.WriteByte('\n')
			}
		}
		return buf.String()
	})

	// 6. Convert markdown unordered list items ("- " / "* ") to .IP roff entries.
	//    roff's fill mode reflows all text, joining adjacent list items onto one
	//    line; .IP begins a fresh indented paragraph for each bullet. .PP after
	//    the final item resets the indent level so subsequent text is not left
	//    hanging at the list indentation.
	s = listBlockRe.ReplaceAllStringFunc(s, func(block string) string {
		var buf strings.Builder
		for _, line := range strings.Split(strings.TrimRight(block, "\n"), "\n") {
			if len(line) >= 2 {
				buf.WriteString(".IP \\(bu 4\n")
				buf.WriteString(line[2:]) // strip "- " / "* " prefix
				buf.WriteByte('\n')
			}
		}
		return buf.String()
	})

	// 7. Convert markdown URLs (roff macros inserted here are the final output,
	//    so they won't be re-escaped).
	//    Two corrections are applied per match:
	//      a) The URL had its dots escaped in step 2 (\. → .) — un-escape them
	//         so .UR receives a valid plain URI.
	//      b) .UR must start at column 1 to be recognised as a macro; if the
	//         preceding text on the same line is non-empty, inject a newline.
	{
		var urlBuf strings.Builder
		prev := 0
		for _, loc := range urlRe.FindAllStringIndex(s, -1) {
			start, end := loc[0], loc[1]
			before := s[prev:start]
			urlBuf.WriteString(before)

			caps := urlRe.FindStringSubmatch(s[start:end])
			if len(caps) >= 3 {
				// Un-escape dots in the URL so the .UR argument is a valid URI.
				url := strings.ReplaceAll(caps[2], `\.`, ".")
				// Normalise the label: collapse newlines and surrounding whitespace.
				label := strings.Join(strings.Fields(caps[1]), " ")
				// Ensure .UR begins on its own line.
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
		s = urlBuf.String()
	}

	// 8. Restore code blocks with .nf/.fi wrappers.
	var out strings.Builder
	pos = 0
	for pos < len(s) {
		for _, b := range blocks {
			if strings.HasPrefix(s[pos:], b.placeholder) {
				out.WriteString(".nf\n")
				out.WriteString(b.content)
				out.WriteString("\n.fi")
				pos += len(b.placeholder)
				break
			}
		}
		if pos >= len(s) {
			break
		}
		r, size := utf8.DecodeRuneInString(s[pos:])
		out.WriteRune(r)
		pos += size
	}

	return out.String()
}

// collectExamples recursively collects all examples from the command tree.
func collectExamples(cmd *spec.CommandItem) []spec.Example {
	var examples []spec.Example
	for _, c := range cmd.Commands {
		if !c.Hidden {
			examples = append(examples, c.Examples...)
			examples = append(examples, collectExamples(c)...)
		}
	}
	return examples
}
