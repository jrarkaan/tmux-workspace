package cmd

import (
	"strings"
	"testing"
)

func TestListCommandPrintsWorkspacesAndWindows(t *testing.T) {
	configFile := writeCommandTestConfig(t)

	output, err := executeCommand("--config", configFile, "list")
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}

	for _, want := range []string{
		"Available workspaces:",
		"backend-dev",
		"frontend-dev",
		"windows : overview, test-watch",
		"windows : overview, dev-server",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("list output = %q, want to contain %q", output, want)
		}
	}

	if strings.Index(output, "backend-dev") > strings.Index(output, "frontend-dev") {
		t.Fatalf("list output is not sorted alphabetically: %q", output)
	}

	if strings.Contains(output, "twx :: declarative tmux workspace manager") {
		t.Fatalf("list output unexpectedly contained banner: %q", output)
	}
}
