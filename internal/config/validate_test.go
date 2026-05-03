package config

import (
	"strings"
	"testing"
)

func TestValidateValidConfig(t *testing.T) {
	if err := Validate(validConfig()); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestValidateDefaultConfig(t *testing.T) {
	if err := Validate(DefaultConfig()); err != nil {
		t.Fatalf("Validate(DefaultConfig()) returned error: %v", err)
	}
}

func TestValidateAllowsEmptyWorkspaces(t *testing.T) {
	cfg := validConfig()
	cfg.Workspaces = map[string]Workspace{}

	if err := Validate(cfg); err != nil {
		t.Fatalf("Validate returned error for empty workspaces: %v", err)
	}
}

func TestValidateInvalidConfigs(t *testing.T) {
	tests := []struct {
		name       string
		cfg        *Config
		wantDetail string
	}{
		{
			name:       "nil config",
			cfg:        nil,
			wantDetail: "config must not be nil",
		},
		{
			name: "unsupported version",
			cfg: mutateValidConfig(func(cfg *Config) {
				cfg.Version = 2
			}),
			wantDetail: "version must be 1",
		},
		{
			name: "missing root",
			cfg: mutateValidConfig(func(cfg *Config) {
				workspace := cfg.Workspaces["backend-dev"]
				workspace.Root = ""
				cfg.Workspaces["backend-dev"] = workspace
			}),
			wantDetail: "root must not be empty",
		},
		{
			name: "empty windows",
			cfg: mutateValidConfig(func(cfg *Config) {
				workspace := cfg.Workspaces["backend-dev"]
				workspace.Windows = nil
				cfg.Workspaces["backend-dev"] = workspace
			}),
			wantDetail: "must have at least one window",
		},
		{
			name: "empty window name",
			cfg: mutateValidConfig(func(cfg *Config) {
				workspace := cfg.Workspaces["backend-dev"]
				workspace.Windows[0].Name = ""
				cfg.Workspaces["backend-dev"] = workspace
			}),
			wantDetail: "name must not be empty",
		},
		{
			name: "duplicate window names",
			cfg: mutateValidConfig(func(cfg *Config) {
				workspace := cfg.Workspaces["backend-dev"]
				workspace.Windows = append(workspace.Windows, Window{Name: "overview"})
				cfg.Workspaces["backend-dev"] = workspace
			}),
			wantDetail: "duplicate window name",
		},
		{
			name: "invalid base index",
			cfg: mutateValidConfig(func(cfg *Config) {
				cfg.Defaults.BaseIndex = 2
			}),
			wantDetail: "defaults.base_index must be 0 or 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.cfg)
			if err == nil {
				t.Fatal("Validate returned nil error")
			}

			if !strings.Contains(err.Error(), tt.wantDetail) {
				t.Fatalf("Validate error = %q, want detail %q", err.Error(), tt.wantDetail)
			}
		})
	}
}

func TestValidateReturnsAllProblems(t *testing.T) {
	cfg := &Config{
		Version: 2,
		Defaults: Defaults{
			BaseIndex: 5,
		},
		Workspaces: map[string]Workspace{
			"backend-dev": {
				Windows: []Window{{Name: ""}},
			},
		},
	}

	err := Validate(cfg)
	if err == nil {
		t.Fatal("Validate returned nil error")
	}

	for _, want := range []string{
		"version must be 1",
		"defaults.base_index must be 0 or 1",
		"root must not be empty",
		"name must not be empty",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("Validate error = %q, want detail %q", err.Error(), want)
		}
	}
}

func validConfig() *Config {
	return &Config{
		Version: 1,
		Defaults: Defaults{
			Attach:          true,
			CreateIfMissing: true,
			BaseIndex:       1,
		},
		Workspaces: map[string]Workspace{
			"backend-dev": {
				Root: "~/App/backend",
				Windows: []Window{
					{Name: "overview", Command: "git status"},
					{Name: "test-watch"},
				},
			},
		},
	}
}

func mutateValidConfig(mutate func(*Config)) *Config {
	cfg := validConfig()
	mutate(cfg)
	return cfg
}
