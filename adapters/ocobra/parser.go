package ocobra

import (
	"strings"

	"github.com/bcdxn/opencli/spec"
)

// Token types for the lexer.
type tokenType int

const (
	tokWORD     tokenType = iota
	tokLBRACKET           // [
	tokRBRACKET           // ]
	tokLBRACE             // {
	tokRBRACE             // }
	tokPIPE               // |
	tokELLIPSIS           // ...
	tokEOF
)

type token struct {
	typ tokenType
	lit string
}

// tokenize splits a Use string into tokens.
// It handles structural characters ([, ], {, }, |) even when they appear
// directly adjacent to words without whitespace (e.g., "[<arg>]").
func tokenize(input string) []token {
	input = strings.TrimSpace(input)
	var tokens []token
	pos := 0

	for pos < len(input) {
		// Skip whitespace
		for pos < len(input) && input[pos] == ' ' {
			pos++
		}
		if pos >= len(input) {
			break
		}

		ch := input[pos]
		switch ch {
		case '[', ']', '{', '}', '|':
			tokens = append(tokens, token{typ: structCharToType(ch), lit: string(ch)})
			pos++

		default:
			// Scan a word — but stop early if we hit a structural character.
			// Special case: if the word starts with '<', scan until '>' so that
			// spaces inside angle brackets are preserved (e.g., "<task description>").
			start := pos
			if input[pos] == '<' {
				for pos < len(input) && input[pos] != '>' {
					pos++
				}
				// Include the closing '>' if present
				if pos < len(input) {
					pos++
				}
			} else {
				for pos < len(input) && input[pos] != ' ' && !isStructChar(input[pos]) {
					pos++
				}
			}
			word := input[start:pos]

			// Handle "..." appended directly to the word (e.g., "<source>...")
			hasEllipsis := strings.HasSuffix(word, "...")
			if hasEllipsis {
				word = strings.TrimSuffix(word, "...")
			}

			if word == "" {
				// Word was just "..." — emit as ELLIPSIS
				tokens = append(tokens, token{tokELLIPSIS, "..."})
			} else {
				tokens = append(tokens, token{tokWORD, word})
				if hasEllipsis {
					tokens = append(tokens, token{tokELLIPSIS, "..."})
				}
			}
		}
	}

	tokens = append(tokens, token{tokEOF, ""})
	return tokens
}

func isStructChar(ch byte) bool {
	return ch == '[' || ch == ']' || ch == '{' || ch == '}' || ch == '|'
}

func structCharToType(ch byte) tokenType {
	switch ch {
	case '[':
		return tokLBRACKET
	case ']':
		return tokRBRACKET
	case '{':
		return tokLBRACE
	case '}':
		return tokRBRACE
	case '|':
		return tokPIPE
	default:
		return tokWORD
	}
}

// parser walks the token stream and extracts ArgumentItems.
type parser struct {
	tokens        []token
	pos           int
	afterDashDash bool // set when "--" is encountered; POSIX: everything after is positional
}

func (p *parser) cur() token {
	if p.pos >= len(p.tokens) {
		return token{tokEOF, ""}
	}
	return p.tokens[p.pos]
}

func (p *parser) advance() token {
	tok := p.cur()
	if p.pos < len(p.tokens) {
		p.pos++
	}
	return tok
}

// ParseUse parses a Cobra Use string and returns the positional arguments found.
func ParseUse(use string) []spec.ArgumentItem {
	tokens := tokenize(use)
	p := &parser{tokens: tokens}
	return p.parseTopLevel()
}

// parseTopLevel processes tokens at the top level of the Use string.
// The first bare word is treated as the command name and skipped.
func (p *parser) parseTopLevel() []spec.ArgumentItem {
	var args []spec.ArgumentItem
	firstBareWord := true
	prevWasFlag := false

	for p.cur().typ != tokEOF {
		tok := p.cur()

		switch tok.typ {
		case tokWORD:
			p.advance()
			name := stripAngleBrackets(tok.lit)

			if name == "--" {
				p.afterDashDash = true
				continue
			}

			if p.afterDashDash {
				// POSIX: everything after -- is a positional argument regardless of shape.
				arg := spec.ArgumentItem{Name: name, Required: true, Passthrough: true}
				if p.cur().typ == tokELLIPSIS {
					p.advance()
					arg.Variadic = true
				}
				args = append(args, arg)
				continue
			}

			if isFlag(name) {
				prevWasFlag = true
				continue
			}

			if prevWasFlag {
				// Word immediately follows a flag → likely a flag value, skip it.
				prevWasFlag = false
				continue
			}

			if isSkipWord(name) {
				continue
			}

			if firstBareWord {
				firstBareWord = false
				continue
			}

			arg := spec.ArgumentItem{
				Name:     name,
				Required: true,
			}

			if p.cur().typ == tokELLIPSIS {
				p.advance()
				arg.Variadic = true
			}

			args = append(args, arg)

		case tokLBRACKET:
			p.advance() // consume [
			groupArgs := p.parseGroup()
			args = append(args, groupArgs...)
			if p.cur().typ == tokRBRACKET {
				p.advance() // consume ]
			}
			// Variadic after the closing bracket applies to the last parsed arg(s)
			if p.cur().typ == tokELLIPSIS {
				p.advance()
				markLastVariadic(args)
			}

		case tokLBRACE:
			p.advance() // consume {
			groupArgs := p.parseGroup()
			args = append(args, groupArgs...)
			if p.cur().typ == tokRBRACE {
				p.advance() // consume }
			}

		case tokELLIPSIS:
			p.advance()

		default:
			p.advance()
		}
	}

	return args
}

// parseGroup processes tokens inside a bracket or brace group.
// Items parsed here are marked as Required: false because they live inside
// an optional [ ] or mutually-exclusive { } construct.
func (p *parser) parseGroup() []spec.ArgumentItem {
	var args []spec.ArgumentItem
	prevWasFlag := false

	for p.cur().typ != tokEOF && p.cur().typ != tokRBRACKET && p.cur().typ != tokRBRACE {
		tok := p.cur()

		switch tok.typ {
		case tokWORD:
			p.advance()
			name := stripAngleBrackets(tok.lit)

			if name == "--" {
				p.afterDashDash = true
				continue
			}

			if p.afterDashDash {
				// POSIX: everything after -- is a positional argument regardless of shape.
				arg := spec.ArgumentItem{Name: name, Required: false, Passthrough: true}
				if p.cur().typ == tokELLIPSIS {
					p.advance()
					arg.Variadic = true
				}
				args = append(args, arg)
				continue
			}

			if isFlag(name) {
				prevWasFlag = true
				continue
			}

			if prevWasFlag {
				prevWasFlag = false
				continue
			}

			if isSkipWord(name) {
				continue
			}

			arg := spec.ArgumentItem{
				Name:     name,
				Required: false,
			}

			if p.cur().typ == tokELLIPSIS {
				p.advance()
				arg.Variadic = true
			}

			args = append(args, arg)

		case tokLBRACKET:
			p.advance()
			nested := p.parseGroup()
			args = append(args, nested...)
			if p.cur().typ == tokRBRACKET {
				p.advance()
			}

		case tokLBRACE:
			p.advance()
			nested := p.parseGroup()
			args = append(args, nested...)
			if p.cur().typ == tokRBRACE {
				p.advance()
			}

		case tokPIPE:
			p.advance()
			prevWasFlag = false // Reset on mutual-exclusion separator

		case tokELLIPSIS:
			p.advance()
			markLastVariadic(args)

		default:
			p.advance()
		}
	}

	return args
}

// markLastVariadic marks the last argument in the slice as variadic.
func markLastVariadic(args []spec.ArgumentItem) {
	for i := len(args) - 1; i >= 0; i-- {
		args[i].Variadic = true
		break
	}
}

// isFlag returns true if the token looks like a CLI flag (-x, --foo, etc.).
// A lone "-" is also treated as a flag-like sentinel and skipped.
func isFlag(s string) bool {
	return strings.HasPrefix(s, "-")
}

// isSkipWord returns true for conventional placeholder tokens that are not
// real positional arguments (e.g., [command], [flags], etc.).
func isSkipWord(s string) bool {
	lower := strings.ToLower(s)
	switch lower {
	case "command", "commands", "subcommand", "subcommands", "sub-command", "sub-commands", "flag", "flags":
		return true
	}
	return false
}

// stripAngleBrackets removes surrounding < > from a name if present and any constant matter
// before or after the bracket, e.g.: `oci://<image-name>` --> `image-name`
// Spaces within the angle brackets are converted to dashes.
func stripAngleBrackets(s string) string {
	s = strings.TrimSpace(s)
	l := strings.Index(s, "<")
	r := strings.Index(s, ">")

	if l >= 0 && r > l {
		result := s[l+1 : r]
		result = strings.ReplaceAll(result, " ", "-")
		return result
	}
	return s
}
