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

func TestAddWindowAppendsWindow(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := AddWindow(cfg, "backend-dev", Window{Name: "jaeger", Command: "echo jaeger"}, false)
	if err != nil {
		t.Fatalf("AddWindow returned error: %v", err)
	}

	workspace := cfg.Workspaces["backend-dev"]
	if len(workspace.Windows) != 3 {
		t.Fatalf("expected 3 windows, got %d", len(workspace.Windows))
	}

	if workspace.Windows[2].Name != "jaeger" {
		t.Fatalf("expected new window 'jaeger' at end, got %s", workspace.Windows[2].Name)
	}
}

func TestAddWindowRejectsDuplicateWithoutForce(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := AddWindow(cfg, "backend-dev", Window{Name: "logs", Command: "tail"}, false)
	if err == nil {
		t.Fatal("AddWindow returned nil error for duplicate window")
	}
	if !strings.Contains(err.Error(), "window already exists") {
		t.Fatalf("AddWindow error = %q, want duplicate detail", err.Error())
	}
}

func TestAddWindowWithForceReplacesCommandAndPreservesPosition(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := AddWindow(cfg, "backend-dev", Window{Name: "overview", Command: "top"}, true)
	if err != nil {
		t.Fatalf("AddWindow returned error: %v", err)
	}

	workspace := cfg.Workspaces["backend-dev"]
	if len(workspace.Windows) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(workspace.Windows))
	}

	if workspace.Windows[0].Name != "overview" || workspace.Windows[0].Command != "top" {
		t.Fatalf("expected overview command 'top', got %+v", workspace.Windows[0])
	}
}

func TestRemoveWindowRemovesExistingWindow(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := RemoveWindow(cfg, "backend-dev", "logs")
	if err != nil {
		t.Fatalf("RemoveWindow returned error: %v", err)
	}

	workspace := cfg.Workspaces["backend-dev"]
	if len(workspace.Windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(workspace.Windows))
	}

	if workspace.Windows[0].Name != "overview" {
		t.Fatalf("expected overview to remain, got %s", workspace.Windows[0].Name)
	}
}

func TestRemoveWindowErrorsIfWorkspaceMissing(t *testing.T) {
	cfg := DefaultConfig()

	err := RemoveWindow(cfg, "missing-workspace", "logs")
	if err == nil {
		t.Fatal("RemoveWindow returned nil error for missing workspace")
	}
}

func TestRemoveWindowErrorsIfWindowMissing(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := RemoveWindow(cfg, "backend-dev", "missing-window")
	if err == nil {
		t.Fatal("RemoveWindow returned nil error for missing window")
	}
}

func TestRemoveWindowErrorsIfRemovingLastWindow(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	// Remove one window to leave only one
	_ = RemoveWindow(cfg, "backend-dev", "logs")

	err := RemoveWindow(cfg, "backend-dev", "overview")
	if err == nil {
		t.Fatal("RemoveWindow returned nil error when removing last window")
	}
	if !strings.Contains(err.Error(), "cannot remove the last window") {
		t.Fatalf("RemoveWindow error = %q, want 'cannot remove the last window' detail", err.Error())
	}
}

func TestSetWindowCommandUpdatesCommand(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := SetWindowCommand(cfg, "backend-dev", "logs", "kubectl logs")
	if err != nil {
		t.Fatalf("SetWindowCommand returned error: %v", err)
	}

	workspace := cfg.Workspaces["backend-dev"]
	if workspace.Windows[1].Name != "logs" || workspace.Windows[1].Command != "kubectl logs" {
		t.Fatalf("expected logs command 'kubectl logs', got %+v", workspace.Windows[1])
	}
}

func TestSetWindowCommandErrorsIfWorkspaceMissing(t *testing.T) {
	cfg := DefaultConfig()

	err := SetWindowCommand(cfg, "missing-workspace", "logs", "kubectl logs")
	if err == nil {
		t.Fatal("SetWindowCommand returned nil error for missing workspace")
	}
}

func TestSetWindowCommandErrorsIfWindowMissing(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	err := SetWindowCommand(cfg, "backend-dev", "missing-window", "echo")
	if err == nil {
		t.Fatal("SetWindowCommand returned nil error for missing window")
	}
}

func TestGetWindow(t *testing.T) {
	cfg := configWithWorkspace("backend-dev")

	win, ok := GetWindow(cfg, "backend-dev", "logs")
	if !ok || win.Name != "logs" {
		t.Fatalf("GetWindow failed to retrieve existing window")
	}

	_, ok = GetWindow(cfg, "backend-dev", "missing")
	if ok {
		t.Fatalf("GetWindow retrieved non-existent window")
	}

	_, ok = GetWindow(cfg, "missing-workspace", "logs")
	if ok {
		t.Fatalf("GetWindow retrieved from non-existent workspace")
	}
}
