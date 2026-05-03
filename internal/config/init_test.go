package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultConfigYAMLLoadsAndValidates(t *testing.T) {
	content, err := DefaultConfigYAML()
	if err != nil {
		t.Fatalf("DefaultConfigYAML returned error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write default config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if err := Validate(cfg); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestInitConfigPrintDoesNotCreateFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config.yaml")

	result, err := InitConfig(path, InitOptions{Print: true})
	if err != nil {
		t.Fatalf("InitConfig returned error: %v", err)
	}

	if !result.Printed {
		t.Fatalf("Printed = false, want true")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("config path exists after print-only init, stat err: %v", err)
	}
}

func TestInitConfigCreatesParentDirectoriesAndConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "twx", "config.yaml")

	result, err := InitConfig(path, InitOptions{})
	if err != nil {
		t.Fatalf("InitConfig returned error: %v", err)
	}

	if !result.Created {
		t.Fatalf("Created = false, want true")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("config file was not created: %v", err)
	}
}

func TestInitConfigDoesNotOverwriteExistingWithoutForce(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	original := []byte("existing: true\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	result, err := InitConfig(path, InitOptions{})
	if err != nil {
		t.Fatalf("InitConfig returned error: %v", err)
	}

	if !result.Existed || result.Created {
		t.Fatalf("result = %#v, want existed without created", result)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if string(content) != string(original) {
		t.Fatalf("config content = %q, want original %q", string(content), string(original))
	}
}

func TestInitConfigForceBacksUpAndWritesDefaultConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	original := []byte("existing: true\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("write existing config: %v", err)
	}

	now := time.Date(2026, 5, 3, 12, 34, 56, 0, time.UTC)
	result, err := InitConfig(path, InitOptions{
		Force: true,
		Now: func() time.Time {
			return now
		},
	})
	if err != nil {
		t.Fatalf("InitConfig returned error: %v", err)
	}

	wantBackup := filepath.Join(dir, "config.yaml.bak.20260503-123456")
	if result.BackupPath != wantBackup {
		t.Fatalf("BackupPath = %q, want %q", result.BackupPath, wantBackup)
	}
	if !strings.Contains(result.BackupPath, "20260503-123456") {
		t.Fatalf("backup path %q missing timestamp", result.BackupPath)
	}

	backupContent, err := os.ReadFile(result.BackupPath)
	if err != nil {
		t.Fatalf("read backup: %v", err)
	}
	if string(backupContent) != string(original) {
		t.Fatalf("backup content = %q, want original %q", string(backupContent), string(original))
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}
