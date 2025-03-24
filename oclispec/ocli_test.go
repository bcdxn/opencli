package oclispec

import (
	"testing"
)

func TestCommandTrieInsert(t *testing.T) {
	trie := CommandTrie{}

	expectedCmds := []string{
		"ocli gen cli",
		"ocli gen doc",
		"ocli spec check",
		"ocli spec versions",
	}

	for _, cmd := range expectedCmds {
		trie.Insert(Command{
			Name: cmd,
		})
	}

	actualCmds := []string{}

	var dfsTrie func(node *CommandTrieNode, cmdLine string)
	dfsTrie = func(node *CommandTrieNode, cmdLine string) {
		var line string
		if cmdLine == "" {
			line = node.Name
		} else {
			line = cmdLine + " " + node.Name
		}

		if len(node.Commands) == 0 {
			actualCmds = append(actualCmds, line)
		}

		for _, node := range node.Commands {
			dfsTrie(node, line)
		}
	}

	dfsTrie(trie.Root, "")

	if len(expectedCmds) != len(actualCmds) {
		t.Fatalf("expected commands %v is not equal to actual commands %v", expectedCmds, actualCmds)
	}

	for i := range actualCmds {
		if actualCmds[i] != expectedCmds[i] {
			t.Errorf("expect cmd '%v' does not equal actual cmd '%v'", actualCmds[i], expectedCmds[i])
		}
	}
}
