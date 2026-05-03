package config

import (
	"strings"
	"testing"
)

func TestAddWorkspaceAddsWorkspace(t *testing.T) {
	cfg := DefaultConfig()

	if err := AddWorkspace(cfg, "backend-dev", sampleWorkspace(), false); err != nil {
		t.Fatalf("AddWorkspace returned error: %v", err)
	}

	if _, ok := cfg.Workspaces["backend-dev"]; !ok {
		t.Fatalf("workspace was not added: %#v", cfg.Workspaces)
	}
}

func TestAddWorkspaceRejectsDuplicateWithoutForce(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := AddWorkspace(cfg, "backend-dev", sampleWorkspace(), false)
	if err == nil {
		t.Fatal("AddWorkspace returned nil error for duplicate")
	}
	if !strings.Contains(err.Error(), "workspace already exists") {
		t.Fatalf("AddWorkspace error = %q, want duplicate detail", err.Error())
	}
}

func TestAddWorkspaceReplacesDuplicateWithForce(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")
	workspace := sampleWorkspace()
	workspace.Root = "~/Other/backend"

	if err := AddWorkspace(cfg, "backend-dev", workspace, true); err != nil {
		t.Fatalf("AddWorkspace returned error: %v", err)
	}

	if got := cfg.Workspaces["backend-dev"].Root; got != "~/Other/backend" {
		t.Fatalf("workspace root = %q, want replacement root", got)
	}
}

func TestRemoveWorkspaceRemovesWorkspace(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	if err := RemoveWorkspace(cfg, "backend-dev"); err != nil {
		t.Fatalf("RemoveWorkspace returned error: %v", err)
	}

	if _, ok := cfg.Workspaces["backend-dev"]; ok {
		t.Fatalf("workspace still exists after remove: %#v", cfg.Workspaces)
	}
}

func TestRemoveWorkspaceErrorsWhenMissing(t *testing.T) {
	cfg := DefaultConfig()

	err := RemoveWorkspace(cfg, "backend-dev")
	if err == nil {
		t.Fatal("RemoveWorkspace returned nil error for missing workspace")
	}
	if !strings.Contains(err.Error(), "workspace not found") {
		t.Fatalf("RemoveWorkspace error = %q, want missing detail", err.Error())
	}
}

func TestAddWorkspaceWithEnvAndCommandsValidates(t *testing.T) {
	cfg := DefaultConfig()

	err := AddWorkspace(cfg, "backend-dev", sampleWorkspace(), false)
	if err != nil {
		t.Fatalf("AddWorkspace returned error: %v", err)
	}
	if err := Validate(cfg); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestAddWorkspaceRejectsDuplicateWindowNames(t *testing.T) {
	cfg := DefaultConfig()
	workspace := sampleWorkspace()
	workspace.Windows = append(workspace.Windows, Window{Name: "overview"})

	err := AddWorkspace(cfg, "backend-dev", workspace, false)
	if err == nil {
		t.Fatal("AddWorkspace returned nil error for duplicate window")
	}
	if !strings.Contains(err.Error(), "duplicate window name") {
		t.Fatalf("AddWorkspace error = %q, want duplicate window detail", err.Error())
	}
}

func sampleWorkspace() Workspace {
	return Workspace{
		Root: "~/App/backend",
		Env: map[string]string{
			"APP_ENV": "development",
		},
		Windows: []Window{
			{Name: "overview", Command: "clear; pwd; git status"},
			{Name: "logs"},
		},
	}
}
