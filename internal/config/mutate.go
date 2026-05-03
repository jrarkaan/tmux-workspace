package config

import (
	"fmt"
	"strings"
)

func AddWorkspace(cfg *Config, name string, workspace Workspace, force bool) error {
	if cfg == nil {
		return fmt.Errorf("config must not be nil")
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("workspace name must not be empty")
	}
	if cfg.Workspaces == nil {
		cfg.Workspaces = map[string]Workspace{}
	}
	if _, exists := cfg.Workspaces[name]; exists && !force {
		return fmt.Errorf("workspace already exists: %s", name)
	}

	cfg.Workspaces[name] = workspace
	return Validate(cfg)
}

func RemoveWorkspace(cfg *Config, name string) error {
	if cfg == nil {
		return fmt.Errorf("config must not be nil")
	}
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("workspace name must not be empty")
	}
	if _, exists := cfg.Workspaces[name]; !exists {
		return fmt.Errorf("workspace not found: %s", name)
	}

	delete(cfg.Workspaces, name)
	return Validate(cfg)
}

func GetWorkspace(cfg *Config, name string) (Workspace, bool) {
	if cfg == nil || cfg.Workspaces == nil {
		return Workspace{}, false
	}

	workspace, ok := cfg.Workspaces[name]
	return workspace, ok
}

func AddWindow(cfg *Config, workspaceName string, window Window, force bool) error {
	workspace, ok := GetWorkspace(cfg, workspaceName)
	if !ok {
		return fmt.Errorf("workspace not found: %s", workspaceName)
	}

	windowName := strings.TrimSpace(window.Name)
	if windowName == "" {
		return fmt.Errorf("window name must not be empty")
	}

	existingIdx := -1
	for i, w := range workspace.Windows {
		if w.Name == windowName {
			existingIdx = i
			break
		}
	}

	if existingIdx != -1 {
		if !force {
			return fmt.Errorf("window already exists: %s\nUse --force to replace it.", windowName)
		}
		workspace.Windows[existingIdx].Command = window.Command
	} else {
		workspace.Windows = append(workspace.Windows, window)
	}

	cfg.Workspaces[workspaceName] = workspace
	return Validate(cfg)
}

func RemoveWindow(cfg *Config, workspaceName string, windowName string) error {
	workspace, ok := GetWorkspace(cfg, workspaceName)
	if !ok {
		return fmt.Errorf("workspace not found: %s", workspaceName)
	}

	windowName = strings.TrimSpace(windowName)

	existingIdx := -1
	for i, w := range workspace.Windows {
		if w.Name == windowName {
			existingIdx = i
			break
		}
	}

	if existingIdx == -1 {
		return fmt.Errorf("window not found: %s", windowName)
	}

	if len(workspace.Windows) <= 1 {
		return fmt.Errorf("cannot remove the last window from workspace: %s", workspaceName)
	}

	workspace.Windows = append(workspace.Windows[:existingIdx], workspace.Windows[existingIdx+1:]...)
	cfg.Workspaces[workspaceName] = workspace
	return Validate(cfg)
}

func SetWindowCommand(cfg *Config, workspaceName string, windowName string, command string) error {
	workspace, ok := GetWorkspace(cfg, workspaceName)
	if !ok {
		return fmt.Errorf("workspace not found: %s", workspaceName)
	}

	windowName = strings.TrimSpace(windowName)

	existingIdx := -1
	for i, w := range workspace.Windows {
		if w.Name == windowName {
			existingIdx = i
			break
		}
	}

	if existingIdx == -1 {
		return fmt.Errorf("window not found: %s", windowName)
	}

	workspace.Windows[existingIdx].Command = command
	cfg.Workspaces[workspaceName] = workspace
	return Validate(cfg)
}

func GetWindow(cfg *Config, workspaceName string, windowName string) (Window, bool) {
	workspace, ok := GetWorkspace(cfg, workspaceName)
	if !ok {
		return Window{}, false
	}

	windowName = strings.TrimSpace(windowName)
	for _, w := range workspace.Windows {
		if w.Name == windowName {
			return w, true
		}
	}

	return Window{}, false
}
