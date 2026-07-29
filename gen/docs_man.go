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

// wsRe matches all new line whitespaces which _can_ occur within a URL label
var wsRe = regexp.MustCompile(`\r?\n+`)

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

	// 3. Convert markdown URLs (roff macros inserted here are the final output,
	//    so they won't be re-escaped).
	s = urlRe.ReplaceAllStringFunc(s, func(match string) string {
		caps := urlRe.FindStringSubmatch(match)
		if len(caps) >= 3 {
			return fmt.Sprintf(".UR %s\n%s\n.UE \\c", caps[2], wsRe.ReplaceAllString(caps[1], " "))
		}
		return match
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
