package ocobra

import (
	"testing"

	"github.com/bcdxn/opencli/spec"
)

func TestParseUse(t *testing.T) {
	tests := []struct {
		name     string
		use      string
		expected []spec.ArgumentItem
	}{
		{
			name: "official example: add [-F file | -D dir]... [-f format] profile",
			use:  "add [-F file | -D dir]... [-f format] profile",
			expected: []spec.ArgumentItem{
				{Name: "profile", Required: true},
			},
		},
		{
			name: "required argument with angle brackets",
			use:  "ls <sub-directory>",
			expected: []spec.ArgumentItem{
				{Name: "sub-directory", Required: true},
			},
		},
		{
			name: "optional argument with angle brackets",
			use:  "ls [<sub-directory>]",
			expected: []spec.ArgumentItem{
				{Name: "sub-directory", Required: false},
			},
		},
		{
			name: "required argument without angle brackets",
			use:  "ls sub-directory",
			expected: []spec.ArgumentItem{
				{Name: "sub-directory", Required: true},
			},
		},
		{
			name: "optional argument without angle brackets",
			use:  "ls [sub-directory]",
			expected: []spec.ArgumentItem{
				{Name: "sub-directory", Required: false},
			},
		},
		{
			name: "variadic argument",
			use:  "cp <source>... <dest>",
			expected: []spec.ArgumentItem{
				{Name: "source", Required: true, Variadic: true},
				{Name: "dest", Required: true},
			},
		},
		{
			name: "mutually exclusive required group {arg1 | arg2}",
			use:  "cmd {arg1 | arg2}",
			expected: []spec.ArgumentItem{
				{Name: "arg1", Required: false},
				{Name: "arg2", Required: false},
			},
		},
		{
			name: "delete {<alias> | --all}",
			use:  "delete {<alias> | --all}",
			expected: []spec.ArgumentItem{
				{Name: "alias", Required: false},
			},
		},
		{
			name: "import [<filename> | -]",
			use:  "import [<filename> | -]",
			expected: []spec.ArgumentItem{
				{Name: "filename", Required: false},
			},
		},
		{
			name: "ssh with special -- separator",
			use:  "ssh [<flags>...] [-- <ssh-flags>...] [<command>]",
			expected: []spec.ArgumentItem{
				{Name: "ssh-flags", Required: false, Variadic: true, Passthrough: true},
				{Name: "command", Required: false, Passthrough: true},
			},
		},
		{
			name: "multiple required arguments",
			use:  "mv <source> <dest>",
			expected: []spec.ArgumentItem{
				{Name: "source", Required: true},
				{Name: "dest", Required: true},
			},
		},
		{
			name: "mixed optional and required",
			use:  "cp [<source>...] <dest>",
			expected: []spec.ArgumentItem{
				{Name: "source", Required: false, Variadic: true},
				{Name: "dest", Required: true},
			},
		},
		{
			name:     "command placeholder skipped",
			use:      "kubectl [command] [flags]",
			expected: []spec.ArgumentItem{},
		},
		{
			name: "sub-command placeholder skipped",
			use:  "git [sub-command] <repo>",
			expected: []spec.ArgumentItem{
				{Name: "repo", Required: true},
			},
		},
		{
			name:     "variadic after bracket group",
			use:      "cmd [-F file | -D dir]...",
			expected: []spec.ArgumentItem{},
		},
		{
			name: "optional mutually exclusive with brackets [{arg1 | arg2}]",
			use:  "cmd [{arg1 | arg2}]",
			expected: []spec.ArgumentItem{
				{Name: "arg1", Required: false},
				{Name: "arg2", Required: false},
			},
		},
		{
			name: "complex nested example",
			use:  "deploy [--image <img>] [<context>] {<app> | --all}",
			expected: []spec.ArgumentItem{
				{Name: "context", Required: false},
				{Name: "app", Required: false},
			},
		},
		{
			name: "constant prefix example",
			use:  "verify [<file-path> | oci://<image-uri>] [--owner | --repo]",
			expected: []spec.ArgumentItem{
				{Name: "file-path", Required: false},
				{Name: "image-uri", Required: false},
			},
		},
		{
			name: "optional argument with angle brackets and flags",
			use:  "create [<task description>] [flags]",
			expected: []spec.ArgumentItem{
				{Name: "task-description", Required: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseUse(tt.use)
			if len(got) != len(tt.expected) {
				t.Fatalf("length mismatch: got %d, want %d\n  got:    %+v\n  want:   %+v", len(got), len(tt.expected), got, tt.expected)
			}
			for i, want := range tt.expected {
				got := got[i]
				if got.Name != want.Name {
					t.Errorf("arg[%d].Name: got %q, want %q", i, got.Name, want.Name)
				}
				if got.Required != want.Required {
					t.Errorf("arg[%d].Required: got %v, want %v", i, got.Required, want.Required)
				}
				if got.Variadic != want.Variadic {
					t.Errorf("arg[%d].Variadic: got %v, want %v", i, got.Variadic, want.Variadic)
				}
			}
		})
	}
}
