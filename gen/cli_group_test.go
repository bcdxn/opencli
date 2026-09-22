package gen

import (
	"strings"
	"testing"

	"github.com/bcdxn/opencli/codec"
	"github.com/bcdxn/opencli/spec"
)

// groupTestDoc builds a minimal spec document whose root command is the given
// item. Used to exercise generation-time rejection of args/flags on groups,
// which validation forbids but callers may skip (e.g., ocli gen cli).
func groupTestDoc(root *spec.CommandItem) *spec.Document {
	return &spec.Document{
		OpenCLIVersion: spec.SchemaVersion,
		Info:           spec.Info{Title: "Group Test CLI", Binary: "grouptest"},
		Commands:       root,
	}
}

func groupTestFlag(name string) []spec.FlagItem {
	return []spec.FlagItem{{Name: name}}
}

// TestCLI_GroupWithFlagsRejected ensures a group command that declares flags is
// rejected at generation time for every framework. The cobra template used to
// emit flag bindings referencing undeclared variables in this case, producing
// code that did not compile; the guard now fails cleanly instead.
func TestCLI_GroupWithFlagsRejected(t *testing.T) {
	// Explicit kind: group with a flag on an explicit-kind root command.
	explicit := groupTestDoc(&spec.CommandItem{
		Segment:  "grouptest",
		Kind:     spec.CommandKindGroup,
		Commands: []*spec.CommandItem{{Segment: "leaf"}},
	})
	explicit.Commands.Flags = groupTestFlag("count")

	// Implicit group (no kind) with a flag on the parent that has children.
	implicit := groupTestDoc(&spec.CommandItem{
		Segment:  "grouptest",
		Commands: []*spec.CommandItem{{Segment: "leaf"}},
	})
	implicit.Commands.Flags = groupTestFlag("count")

	for _, framework := range []CLIFramework{CobraFramework, UrfaveCliFramework, YargsFramework} {
		t.Run(string(framework), func(t *testing.T) {
			for name, doc := range map[string]*spec.Document{"explicit-kind": explicit, "implicit-group": implicit} {
				if _, err := CLI(doc, GenCLIWithFramework(framework)); err == nil {
					t.Errorf("%s: expected error for group command with flags, got none", name)
				} else if !strings.Contains(err.Error(), "group commands cannot have arguments or flags") {
					t.Errorf("%s: unexpected error message %q", name, err)
				}
			}
		})
	}
}

// TestCLI_GroupWithArgsRejected ensures a group command that declares positional
// arguments is rejected at generation time, mirroring the validation rule.
func TestCLI_GroupWithArgsRejected(t *testing.T) {
	doc := groupTestDoc(&spec.CommandItem{
		Segment:  "grouptest",
		Kind:     spec.CommandKindGroup,
		Commands: []*spec.CommandItem{{Segment: "leaf"}},
	})
	doc.Commands.Args = []spec.ArgumentItem{{Name: "name"}}

	for _, framework := range []CLIFramework{CobraFramework, UrfaveCliFramework, YargsFramework} {
		if _, err := CLI(doc, GenCLIWithFramework(framework)); err == nil {
			t.Errorf("expected error for group command with arguments (framework %s), got none", framework)
		} else if !strings.Contains(err.Error(), "group commands cannot have arguments or flags") {
			t.Errorf("unexpected error message: %q", err)
		}
	}
}

// TestCLI_GroupWithNestedGroupRejected ensures the guard recurses into nested
// groups, not just the root command.
func TestCLI_GroupWithNestedGroupRejected(t *testing.T) {
	doc := groupTestDoc(&spec.CommandItem{
		Segment: "grouptest",
		Commands: []*spec.CommandItem{{
			Segment:  "mid",
			Kind:     spec.CommandKindGroup,
			Flags:    groupTestFlag("count"),
			Commands: []*spec.CommandItem{{Segment: "leaf"}},
		}},
	})

	if _, err := CLI(doc, GenCLIWithFramework(CobraFramework)); err == nil {
		t.Fatal("expected error for nested group command with flags, got none")
	} else if !strings.Contains(err.Error(), `"mid"`) && !strings.Contains(err.Error(), "grouptest mid") {
		t.Errorf("error should name the offending group: %q", err)
	}
}

// TestCLI_ValidGroupStillGenerates ensures flag-free groups (the common case, e.g.
// petstore's bare `kind: group` commands) still generate successfully and are not
// over-rejected by the guard.
func TestCLI_ValidGroupStillGenerates(t *testing.T) {
	doc := groupTestDoc(&spec.CommandItem{
		Segment: "grouptest",
		Kind:    spec.CommandKindGroup,
		Commands: []*spec.CommandItem{{
			Segment: "leaf",
			Flags:   groupTestFlag("count"), // flags on the LEAF are fine
		}},
	})

	for _, framework := range []CLIFramework{CobraFramework, UrfaveCliFramework, YargsFramework} {
		files, err := CLI(doc, GenCLIWithFramework(framework))
		if err != nil {
			t.Errorf("framework %s: unexpected error for valid group spec: %v", framework, err)
			continue
		}
		if len(files) == 0 {
			t.Errorf("framework %s: expected generated files but got none", framework)
		}
	}

	// The embedded petstore example (the golden source) has bare groups and must
	// still generate without error.
	petstore, err := codec.UnmarshalYAML(exampleYAML)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling example OpenCLI doc: %v", err)
	}
	for _, framework := range []CLIFramework{CobraFramework, UrfaveCliFramework, YargsFramework} {
		files, err := CLI(petstore, GenCLIWithFramework(framework))
		if err != nil {
			t.Errorf("framework %s: unexpected error for petstore spec: %v", framework, err)
		} else if len(files) == 0 {
			t.Errorf("framework %s: expected generated files but got none", framework)
		}
	}
}
