package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSaveWritesLoadableValidYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	cfg := configWithWorkspace("backend-dev")

	if err := Save(path, cfg); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if err := Validate(loaded); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestBackupCreatesTimestampedBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("version: 1\n"), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	now := time.Date(2026, 5, 3, 8, 30, 0, 0, time.UTC)
	backupPath, err := Backup(path, now)
	if err != nil {
		t.Fatalf("Backup returned error: %v", err)
	}

	if !strings.HasSuffix(backupPath, "config.yaml.bak.20260503-083000") {
		t.Fatalf("Backup path = %q, want timestamp suffix", backupPath)
	}
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("backup file was not created: %v", err)
	}
}

func TestSaveWithBackupBacksUpThenWrites(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	original := []byte("version: 1\nworkspaces: {}\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	backupPath, err := SaveWithBackup(path, configWithWorkspace("backend-dev"), time.Date(2026, 5, 3, 8, 31, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("SaveWithBackup returned error: %v", err)
	}

	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(backupContent) != string(original) {
		t.Fatalf("backup content = %q, want original %q", string(backupContent), string(original))
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if _, ok := loaded.Workspaces["backend-dev"]; !ok {
		t.Fatalf("workspace missing after SaveWithBackup: %#v", loaded.Workspaces)
	}
}

func configWithWorkspace(name string) *Config {
	cfg := DefaultConfig()
	cfg.Workspaces[name] = Workspace{
		Root: "~/App/backend",
		Env: map[string]string{
			"APP_ENV": "development",
		},
		Windows: []Window{
			{Name: "overview", Command: "git status"},
			{Name: "logs"},
		},
	}
	return cfg
}
