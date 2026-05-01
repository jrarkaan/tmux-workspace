package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadValidYAML(t *testing.T) {
	path := writeTempConfig(t, `version: 1
defaults:
  attach: true
  create_if_missing: true
  base_index: 1
workspaces:
  backend-dev:
    root: ~/App/backend
    env:
      APP_ENV: development
    windows:
      - name: overview
        command: clear; pwd; git status
      - name: test-watch
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Version != 1 {
		t.Fatalf("Version = %d, want 1", cfg.Version)
	}
	if len(cfg.Workspaces) != 1 {
		t.Fatalf("workspace count = %d, want 1", len(cfg.Workspaces))
	}

	workspace := cfg.Workspaces["backend-dev"]
	if workspace.Root != "~/App/backend" {
		t.Fatalf("workspace root = %q, want ~/App/backend", workspace.Root)
	}
	if workspace.Env["APP_ENV"] != "development" {
		t.Fatalf("APP_ENV = %q, want development", workspace.Env["APP_ENV"])
	}
	if len(workspace.Windows) != 2 {
		t.Fatalf("window count = %d, want 2", len(workspace.Windows))
	}
}

func TestLoadInvalidYAMLReturnsError(t *testing.T) {
	path := writeTempConfig(t, "version: [")

	if _, err := Load(path); err == nil {
		t.Fatal("Load returned nil error for invalid YAML")
	}
}

func TestLoadMissingFileReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")

	if _, err := Load(path); err == nil {
		t.Fatal("Load returned nil error for missing file")
	}
}

func TestResolveConfigPath(t *testing.T) {
	override := filepath.Join(t.TempDir(), "config.yaml")
	if got := ResolveConfigPath(override); got != override {
		t.Fatalf("ResolveConfigPath(%q) = %q, want %q", override, got, override)
	}

	if got := ResolveConfigPath(""); got != DefaultConfigPath() {
		t.Fatalf("ResolveConfigPath(\"\") = %q, want %q", got, DefaultConfigPath())
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("os.UserHomeDir unavailable: %v", err)
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "home only", path: "~", want: home},
		{name: "home child", path: "~/abc", want: filepath.Join(home, "abc")},
		{name: "unchanged", path: "/tmp/twx/config.yaml", want: "/tmp/twx/config.yaml"},
		{name: "tilde prefix unchanged", path: "~abc/config.yaml", want: "~abc/config.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExpandHome(tt.path); got != tt.want {
				t.Fatalf("ExpandHome(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	return path
}
