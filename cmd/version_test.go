package cmd

import (
	"bytes"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	rootCmd := newRootCommand()
	var output bytes.Buffer

	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	if got, want := output.String(), "twx version dev\n"; got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}
