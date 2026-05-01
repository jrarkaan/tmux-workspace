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
		"start",
		"windows",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("root help output = %q, want to contain %q", output, want)
		}
	}
}

func TestStartCommandMissingWorkspaceArgReturnsError(t *testing.T) {
	_, err := executeCommand("start")
	if err == nil {
		t.Fatal("start without workspace returned nil error")
	}

	if !strings.Contains(err.Error(), "accepts 1 arg(s), received 0") {
		t.Fatalf("start error = %q, want missing arg detail", err.Error())
	}
}
