package tpm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type fakeRunner struct {
	commands [][]string
	err      error
}

func (f *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := append([]string{name}, args...)
	f.commands = append(f.commands, cmd)
	if name == "git" && args[0] == "--version" {
		return []byte("git version 2.43.0\n"), f.err
	}
	if name == "tmux" && args[0] == "-V" {
		return []byte("tmux 3.4\n"), f.err
	}
	return nil, f.err
}

func newTestManager(homeDir string) *Manager {
	m := NewManager(&fakeRunner{}, homeDir)
	m.now = func() time.Time {
		return time.Date(2026, 5, 3, 10, 0, 0, 0, time.UTC)
	}
	m.lookPath = func(file string) (string, error) {
		return "/usr/bin/" + file, nil
	}
	return m
}

func TestStatusMissingTPM(t *testing.T) {
	tempDir := t.TempDir()
	m := newTestManager(tempDir)

	s, err := m.GetStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.TPMDirExists {
		t.Fatal("expected TPMDirExists to be false")
	}
	if s.TmuxConfigExists {
		t.Fatal("expected TmuxConfigExists to be false")
	}
}

func TestStatusDetectsTPMAndConfig(t *testing.T) {
	tempDir := t.TempDir()
	os.MkdirAll(filepath.Join(tempDir, ".tmux", "plugins", "tpm"), 0755)

	configContent := `
set -g @plugin 'tmux-plugins/tpm'
run '~/.tmux/plugins/tpm/tpm'
	`
	os.WriteFile(filepath.Join(tempDir, ".tmux.conf"), []byte(configContent), 0644)

	m := newTestManager(tempDir)
	s, err := m.GetStatus()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !s.TPMDirExists {
		t.Fatal("expected TPMDirExists to be true")
	}
	if !s.TmuxConfigExists {
		t.Fatal("expected TmuxConfigExists to be true")
	}
	if !s.HasTPMPlugin {
		t.Fatal("expected HasTPMPlugin to be true")
	}
	if !s.HasTPMRunLine {
		t.Fatal("expected HasTPMRunLine to be true")
	}
}

func TestInstallCreatesConfigAndClonesTPM(t *testing.T) {
	tempDir := t.TempDir()
	m := newTestManager(tempDir)

	backupPath, mutated, err := m.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mutated {
		t.Fatal("expected mutated to be true")
	}
	if backupPath != "" {
		t.Fatalf("expected empty backupPath because config didn't exist, got %s", backupPath)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, ".tmux.conf"))
	if err != nil {
		t.Fatalf("failed to read .tmux.conf: %v", err)
	}
	if !strings.Contains(string(content), "# twx TPM managed block") {
		t.Fatal("expected config to contain managed block")
	}

	runner := m.runner.(*fakeRunner)
	if len(runner.commands) != 1 || runner.commands[0][0] != "git" || runner.commands[0][1] != "clone" {
		t.Fatalf("expected git clone command, got %v", runner.commands)
	}
}

func TestInstallBacksUpExistingConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".tmux.conf")
	os.WriteFile(configPath, []byte("set -g mouse on\n"), 0644)

	m := newTestManager(tempDir)

	backupPath, mutated, err := m.Install()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mutated {
		t.Fatal("expected mutated to be true")
	}
	if backupPath == "" {
		t.Fatal("expected backupPath, got empty string")
	}
	if !strings.HasPrefix(backupPath, configPath+".bak.") {
		t.Fatalf("unexpected backupPath format: %s", backupPath)
	}

	backupContent, _ := os.ReadFile(backupPath)
	if string(backupContent) != "set -g mouse on\n" {
		t.Fatalf("unexpected backup content: %q", string(backupContent))
	}
}

func TestInstallIsIdempotent(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".tmux.conf")
	os.MkdirAll(filepath.Join(tempDir, ".tmux", "plugins", "tpm"), 0755)

	m := newTestManager(tempDir)

	// First install
	_, mutated1, err := m.Install()
	if err != nil {
		t.Fatalf("unexpected error on first install: %v", err)
	}
	if !mutated1 {
		t.Fatal("expected mutated on first install")
	}

	// Second install
	backupPath, mutated2, err := m.Install()
	if err != nil {
		t.Fatalf("unexpected error on second install: %v", err)
	}
	if mutated2 {
		t.Fatal("expected not mutated on second install")
	}
	if backupPath != "" {
		t.Fatal("expected no backup path on second install")
	}

	content, _ := os.ReadFile(configPath)
	if strings.Count(string(content), "# twx TPM managed block") != 1 {
		t.Fatal("expected only one managed block")
	}
}

func TestInstallFailsIfGitMissing(t *testing.T) {
	tempDir := t.TempDir()
	m := newTestManager(tempDir)
	m.lookPath = func(file string) (string, error) {
		if file == "git" {
			return "", fmt.Errorf("not found")
		}
		return "/usr/bin/" + file, nil
	}

	_, _, err := m.Install()
	if err == nil {
		t.Fatal("expected error when git is missing")
	}
}
