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
