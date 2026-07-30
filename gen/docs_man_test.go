package gen

import (
	"strings"
	"testing"
)

func TestRoffEscape_CodeBlocks(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSub []string // substrings that must appear in output
		notSub  []string // substrings that must NOT appear in output
	}{
		{
			name:    "single fenced code block",
			input:   "before\n```\ncode here\n```\nafter",
			wantSub: []string{".nf\n```\ncode here\n```\n.fi"},
			notSub:  nil,
		},
		{
			name:    "code block with language tag",
			input:   "text\n```go\nfmt.Println(\"hello\")\n```\ndone",
			wantSub: []string{".nf\n```go\nfmt.Println(\"hello\")\n```\n.fi"},
			notSub:  nil,
		},
		{
			name:    "multiple code blocks",
			input:   "```\nblock1\n```\nmiddle\n```\nblock2\n```",
			wantSub: []string{".nf\n```\nblock1\n```\n.fi", ".nf\n```\nblock2\n```\n.fi"},
			notSub:  nil,
		},
		{
			name:    "code block content not escaped",
			input:   "```\nhello.world test\\backslash\n```",
			wantSub: []string{"hello.world test\\backslash"},
			notSub:  []string{`hello\.world`, `\bK`},
		},
		{
			name:    "code block with backticks inside",
			input:   "before\n````\n```nested```\n````\nafter",
			wantSub: []string{".nf"},
			notSub:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("expected output to contain %q, got:\n%s", sub, got)
				}
			}
			for _, sub := range tt.notSub {
				if strings.Contains(got, sub) {
					t.Errorf("expected output to NOT contain %q, got:\n%s", sub, got)
				}
			}
		})
	}
}

func TestRoffEscape_CharacterEscaping(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "backslash escaped",
			input: `hello\world`,
			want:  `hello\bKworld`,
		},
		{
			name:  "dot escaped",
			input: "hello.world",
			want:  `hello\.world`,
		},
		{
			name:  "multiple special chars",
			input: `a.b\c.d`,
			want:  `a\.b\bKc\.d`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected output to contain %q, got:\n%s", tt.want, got)
			}
		})
	}
}

func TestRoffEscape_Paragraphs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "two newlines become .PP",
			input: "para1\n\npara2",
			want:  "\n.PP\n",
		},
		{
			name:  "multiple blank lines collapse to single .PP",
			input: "para1\n\n\n\npara2",
			want:  "\n.PP\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected output to contain %q, got:\n%s", tt.want, got)
			}
		})
	}
}

func TestRoffEscape_HorizontalBreaks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "three dashes become horizontal rule",
			input: "before\n---\nafter",
			want:  "\n\\l'20n'\n",
		},
		{
			name:  "three tildes become horizontal rule",
			input: "before\n~~~\nafter",
			want:  "\n\\l'20n'\n",
		},
		{
			name:  "five dashes also work",
			input: "before\n-----\nafter",
			want:  "\n\\l'20n'\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			if !strings.Contains(got, tt.want) {
				t.Errorf("expected output to contain %q, got:\n%s", tt.want, got)
			}
		})
	}
}

func TestRoffEscape_BlockQuotes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSub []string
	}{
		{
			name:    "simple block quote",
			input:   "> quoted text",
			wantSub: []string{".PP\n"},
		},
		{
			name:    "multi-line block quote",
			input:   "> line one\n> line two",
			wantSub: []string{".PP\n", "line one", "line two"},
		},
		{
			name:    "block quote with > only (no space)",
			input:   ">text",
			wantSub: []string{".PP\n", "text"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("expected output to contain %q, got:\n%s", sub, got)
				}
			}
		})
	}
}

func TestRoffEscape_GFMAlerts(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSub []string
	}{
		{
			name:    "tip alert",
			input:   "> [!TIP]\n> helpful info",
			wantSub: []string{"\\fBTip:\\fP"},
		},
		{
			name:    "note alert",
			input:   "> [!NOTE]\n> some note",
			wantSub: []string{"\\fBNote:\\fP"},
		},
		{
			name:    "warning alert",
			input:   "> [!WARNING]\n> be careful",
			wantSub: []string{"\\fBWarning:\\fP"},
		},
		{
			name:    "caution alert",
			input:   "> [!CAUTION]\n> dangerous",
			wantSub: []string{"\\fBCaution:\\fP"},
		},
		{
			name:    "important alert",
			input:   "> [!IMPORTANT]\n> critical info",
			wantSub: []string{"\\fBImportant:\\fP"},
		},
		{
			name:    "case insensitive alert marker",
			input:   "> [!tip]\n> lowercase tip",
			wantSub: []string{"\\fBTip:\\fP"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("expected output to contain %q, got:\n%s", sub, got)
				}
			}
		})
	}
}

func TestRoffEscape_UnorderedLists(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSub []string
	}{
		{
			name:    "dash list items",
			input:   "- item one\n- item two",
			wantSub: []string{".IP \\(bu 4\nitem one", ".IP \\(bu 4\nitem two"},
		},
		{
			name:    "asterisk list items",
			input:   "* item one\n* item two",
			wantSub: []string{".IP \\(bu 4\nitem one", ".IP \\(bu 4\nitem two"},
		},
		{
			name:    "single list item",
			input:   "- only item",
			wantSub: []string{".IP \\(bu 4\nonly item"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("expected output to contain %q, got:\n%s", sub, got)
				}
			}
		})
	}
}

func TestRoffEscape_URLs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSub []string
	}{
		{
			name:    "basic http url",
			input:   "[label](http://example.com)",
			wantSub: []string{".UR http://example.com\nlabel\n.UE \\c"},
		},
		{
			name:    "https url",
			input:   "[click here](https://example.com/path)",
			wantSub: []string{".UR https://example.com/path\nclick here\n.UE \\c"},
		},
		{
			name:    "url with dots in path",
			input:   "[docs](https://example.com/docs/v1.0/guide.html)",
			wantSub: []string{".UR https://example.com/docs/v1.0/guide.html\ndocs\n.UE \\c"},
		},
		{
			name:    "url preceded by text gets newline",
			input:   "see [link](http://example.com) for more",
			wantSub: []string{"\n.UR http://example.com\nlink\n.UE \\c"},
		},
		{
			name:    "label with newlines is collapsed",
			input:   "[multi\nline\nlabel](http://example.com)",
			wantSub: []string{"multi line label"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("expected output to contain %q, got:\n%s", sub, got)
				}
			}
		})
	}
}

func TestRoffEscape_Integration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantSub []string
	}{
		{
			name:  "mixed content with code block, list, and paragraph",
			input: "Introduction\n\n```\ncode block\n```\n\n- item1\n- item2\n\nFinal paragraph.",
			wantSub: []string{
				".nf\n```\ncode block\n```\n.fi",
				".IP \\(bu 4\nitem1",
				".IP \\(bu 4\nitem2",
				"\n.PP\n",
			},
		},
		{
			name:    "empty string",
			input:   "",
			wantSub: nil,
		},
		{
			name:    "plain text no special chars",
			input:   "just plain text here",
			wantSub: []string{"just plain text here"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roffEscape(tt.input)
			if tt.wantSub == nil && tt.input == "" {
				if got != "" {
					t.Errorf("expected empty output for empty input, got: %q", got)
				}
				return
			}
			for _, sub := range tt.wantSub {
				if !strings.Contains(got, sub) {
					t.Errorf("expected output to contain %q, got:\n%s", sub, got)
				}
			}
		})
	}
}
