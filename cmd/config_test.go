package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigPathCommandUsesOverride(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")
	output, err := executeCommand("config", "path", "--config", configFile)
	if err != nil {
		t.Fatalf("config path failed: %v", err)
	}

	if got, want := strings.TrimSpace(output), configFile; got != want {
		t.Fatalf("config path output = %q, want %q", got, want)
	}
}

func TestConfigValidateCommand(t *testing.T) {
	configFile := writeCommandTestConfig(t)

	output, err := executeCommand("--config", configFile, "config", "validate")
	if err != nil {
		t.Fatalf("config validate failed: %v", err)
	}

	if !strings.Contains(output, "Config is valid: "+configFile) {
		t.Fatalf("config validate output missing valid message: %q", output)
	}
	if !strings.Contains(output, "Workspaces: 2") {
		t.Fatalf("config validate output missing workspace count: %q", output)
	}
}

func executeCommand(args ...string) (string, error) {
	rootCmd := newRootCommand()
	var output bytes.Buffer

	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)
	rootCmd.SetArgs(args)

	err := rootCmd.Execute()
	return output.String(), err
}

func writeCommandTestConfig(t *testing.T) string {
	t.Helper()

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	content := `version: 1
defaults:
  attach: true
  create_if_missing: true
  base_index: 1
workspaces:
  frontend-dev:
    root: ~/App/frontend
    windows:
      - name: overview
      - name: dev-server
  backend-dev:
    root: ~/App/backend
    env:
      APP_ENV: development
    windows:
      - name: overview
        command: clear; pwd; git status
      - name: test-watch
`

	if err := os.WriteFile(configFile, []byte(content), 0o600); err != nil {
		t.Fatalf("write command test config: %v", err)
	}

	return configFile
}
