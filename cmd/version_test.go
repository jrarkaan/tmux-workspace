package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootHelpIncludesBanner(t *testing.T) {
	rootCmd := newRootCommand()
	var output bytes.Buffer

	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("root help failed: %v", err)
	}

	if !strings.Contains(output.String(), "twx :: declarative tmux workspace manager") {
		t.Fatalf("root help did not contain banner tagline: %q", output.String())
	}
}

func TestVersionCommand(t *testing.T) {
	output, err := executeCommand("version")
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if got, want := output, "twx version dev\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestDoctorCommandDoesNotIncludeBanner(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	output, err := executeCommand("--config", configPath, "doctor")
	if err != nil {
		t.Fatalf("doctor command failed: %v", err)
	}

	if strings.Contains(output, "twx :: declarative tmux workspace manager") {
		t.Fatalf("doctor output unexpectedly contained banner: %q", output)
	}
}
