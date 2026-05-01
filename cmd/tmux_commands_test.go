package cmd

import (
	"strings"
	"testing"
)

func TestRootHelpIncludesTmuxInspectionCommands(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("root help failed: %v", err)
	}

	for _, want := range []string{
		"sessions",
		"windows",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("root help output = %q, want to contain %q", output, want)
		}
	}
}
