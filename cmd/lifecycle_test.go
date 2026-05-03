package cmd

import (
	"strings"
	"testing"
)

func TestLifecycleCommandsAppearInRootHelp(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("root help failed: %v", err)
	}

	for _, want := range []string{"attach", "kill", "restart"} {
		if !strings.Contains(output, want) {
			t.Fatalf("root help output = %q, want %q", output, want)
		}
	}
}

func TestLifecycleCommandsMissingWorkspaceArgReturnError(t *testing.T) {
	for _, command := range []string{"attach", "kill", "restart"} {
		t.Run(command, func(t *testing.T) {
			_, err := executeCommand(command)
			if err == nil {
				t.Fatalf("%s without workspace returned nil error", command)
			}
			if !strings.Contains(err.Error(), "accepts 1 arg(s), received 0") {
				t.Fatalf("%s error = %q, want missing arg detail", command, err.Error())
			}
		})
	}
}

func TestLifecycleCommandsMissingWorkspaceReturnWorkspaceNotFound(t *testing.T) {
	configFile := initCommandConfig(t)

	for _, command := range []string{"attach", "kill", "restart"} {
		t.Run(command, func(t *testing.T) {
			output, err := executeCommand("--config", configFile, command, "missing-workspace")
			if err == nil {
				t.Fatalf("%s missing workspace returned nil error", command)
			}
			if !strings.Contains(err.Error(), "workspace not found: missing-workspace") {
				t.Fatalf("%s error = %q, want workspace not found", command, err.Error())
			}
			if strings.Contains(output, "twx :: declarative tmux workspace manager") {
				t.Fatalf("%s output unexpectedly contained banner: %q", command, output)
			}
		})
	}
}
