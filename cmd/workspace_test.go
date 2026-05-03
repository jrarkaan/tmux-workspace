package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jrarkaan/tmux-workspace/internal/config"
)

func TestWorkspaceCommandAppearsInRootHelp(t *testing.T) {
	output, err := executeCommand("--help")
	if err != nil {
		t.Fatalf("root help failed: %v", err)
	}

	if !strings.Contains(output, "workspace") {
		t.Fatalf("root help output = %q, want workspace command", output)
	}
}

func TestWorkspaceAddWithTempConfig(t *testing.T) {
	configFile := initCommandConfig(t)

	output, err := executeCommand(
		"--config", configFile,
		"workspace", "add", "backend-dev",
		"--root", "~/App/backend",
		"--windows", "overview,codex-impl,logs",
		"--env", "APP_ENV=development",
		"--command", "overview=clear; pwd; git status",
	)
	if err != nil {
		t.Fatalf("workspace add failed: %v", err)
	}

	for _, want := range []string{
		"Workspace added: backend-dev",
		"Config updated: " + configFile,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("workspace add output = %q, want %q", output, want)
		}
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	workspace := cfg.Workspaces["backend-dev"]
	if workspace.Root != "~/App/backend" {
		t.Fatalf("root = %q, want preserved root", workspace.Root)
	}
	if workspace.Env["APP_ENV"] != "development" {
		t.Fatalf("APP_ENV = %q, want development", workspace.Env["APP_ENV"])
	}
	if workspace.Windows[0].Command != "clear; pwd; git status" {
		t.Fatalf("overview command = %q", workspace.Windows[0].Command)
	}
}

func TestWorkspaceAddFailsIfConfigMissing(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yaml")

	_, err := executeCommand("--config", configFile, "workspace", "add", "backend-dev", "--root", "~/App/backend", "--windows", "overview")
	if err == nil {
		t.Fatal("workspace add returned nil error for missing config")
	}

	for _, want := range []string{
		"Config file does not exist: " + configFile,
		"Run: twx config init",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("workspace add error = %q, want %q", err.Error(), want)
		}
	}
}

func TestWorkspaceAddRejectsDuplicateWithoutForce(t *testing.T) {
	configFile := initCommandConfig(t)
	addWorkspaceForCommandTest(t, configFile, "~/App/backend", "overview,logs")

	_, err := executeCommand("--config", configFile, "workspace", "add", "backend-dev", "--root", "~/Other/backend", "--windows", "overview")
	if err == nil {
		t.Fatal("workspace add duplicate returned nil error")
	}
	if !strings.Contains(err.Error(), "workspace already exists: backend-dev") ||
		!strings.Contains(err.Error(), "Use --force to replace it.") {
		t.Fatalf("workspace add duplicate error = %q", err.Error())
	}
}

func TestWorkspaceAddForceReplacesAndCreatesBackup(t *testing.T) {
	configFile := initCommandConfig(t)
	addWorkspaceForCommandTest(t, configFile, "~/App/backend", "overview,logs")

	output, err := executeCommand("--config", configFile, "workspace", "add", "backend-dev", "--root", "~/Other/backend", "--windows", "overview", "--force")
	if err != nil {
		t.Fatalf("workspace add --force failed: %v", err)
	}

	for _, want := range []string{
		"Workspace replaced: backend-dev",
		"Existing config backed up: " + configFile + ".bak.",
		"Config updated: " + configFile,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("workspace add --force output = %q, want %q", output, want)
		}
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got := cfg.Workspaces["backend-dev"].Root; got != "~/Other/backend" {
		t.Fatalf("root = %q, want replacement root", got)
	}
	assertBackupCount(t, configFile, 2)
}

func TestWorkspaceShowPrintsDetails(t *testing.T) {
	configFile := initCommandConfig(t)
	addWorkspaceForCommandTest(t, configFile, "~/App/backend", "overview,logs")

	output, err := executeCommand("--config", configFile, "workspace", "show", "backend-dev")
	if err != nil {
		t.Fatalf("workspace show failed: %v", err)
	}

	for _, want := range []string{
		"Workspace: backend-dev",
		"Root     : ~/App/backend",
		"Environment:",
		"APP_ENV=development",
		"1. overview",
		"command: clear; pwd; git status",
		"2. logs",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("workspace show output = %q, want %q", output, want)
		}
	}
	if strings.Contains(output, "twx :: declarative tmux workspace manager") {
		t.Fatalf("workspace show output unexpectedly contained banner: %q", output)
	}
}

func TestWorkspaceRemoveWithoutForceDoesNotModifyConfig(t *testing.T) {
	configFile := initCommandConfig(t)
	addWorkspaceForCommandTest(t, configFile, "~/App/backend", "overview,logs")

	output, err := executeCommand("--config", configFile, "workspace", "remove", "backend-dev")
	if err == nil {
		t.Fatal("workspace remove without force returned nil error")
	}
	for _, want := range []string{
		"This will remove workspace from config: backend-dev",
		"Re-run with --force to confirm.",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("workspace remove output = %q, want %q", output, want)
		}
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if _, ok := cfg.Workspaces["backend-dev"]; !ok {
		t.Fatal("workspace was removed without --force")
	}
}

func TestWorkspaceRemoveForceRemovesAndCreatesBackup(t *testing.T) {
	configFile := initCommandConfig(t)
	addWorkspaceForCommandTest(t, configFile, "~/App/backend", "overview,logs")

	output, err := executeCommand("--config", configFile, "workspace", "remove", "backend-dev", "--force")
	if err != nil {
		t.Fatalf("workspace remove --force failed: %v", err)
	}
	for _, want := range []string{
		"Workspace removed: backend-dev",
		"Existing config backed up: " + configFile + ".bak.",
		"Config updated: " + configFile,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("workspace remove --force output = %q, want %q", output, want)
		}
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if _, ok := cfg.Workspaces["backend-dev"]; ok {
		t.Fatal("workspace still exists after remove --force")
	}
	assertBackupCount(t, configFile, 2)
}

func initCommandConfig(t *testing.T) string {
	t.Helper()

	configFile := filepath.Join(t.TempDir(), "config.yaml")
	if _, err := executeCommand("--config", configFile, "config", "init"); err != nil {
		t.Fatalf("config init failed: %v", err)
	}
	return configFile
}

func addWorkspaceForCommandTest(t *testing.T, configFile string, root string, windows string) {
	t.Helper()

	_, err := executeCommand(
		"--config", configFile,
		"workspace", "add", "backend-dev",
		"--root", root,
		"--windows", windows,
		"--env", "APP_ENV=development",
		"--command", "overview=clear; pwd; git status",
	)
	if err != nil {
		t.Fatalf("workspace add failed: %v", err)
	}
}

func assertBackupCount(t *testing.T, configFile string, want int) {
	t.Helper()

	matches, err := filepath.Glob(configFile + ".bak.*")
	if err != nil {
		t.Fatalf("glob backups: %v", err)
	}
	if len(matches) != want {
		entries, _ := os.ReadDir(filepath.Dir(configFile))
		t.Fatalf("backup count = %d, want %d, matches=%#v entries=%#v", len(matches), want, matches, entries)
	}
}
