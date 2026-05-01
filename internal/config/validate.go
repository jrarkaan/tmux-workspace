package config

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Problems []string
}

func (e ValidationError) Error() string {
	if len(e.Problems) == 0 {
		return "config validation failed"
	}

	return "config validation failed: " + strings.Join(e.Problems, "; ")
}

func Validate(cfg *Config) error {
	var problems []string

	if cfg == nil {
		return ValidationError{Problems: []string{"config must not be nil"}}
	}

	if cfg.Version != 1 {
		problems = append(problems, fmt.Sprintf("version must be 1, got %d", cfg.Version))
	}

	if cfg.Defaults.BaseIndex != 0 && cfg.Defaults.BaseIndex != 1 {
		problems = append(problems, fmt.Sprintf("defaults.base_index must be 0 or 1, got %d", cfg.Defaults.BaseIndex))
	}

	if len(cfg.Workspaces) == 0 {
		problems = append(problems, "workspaces must not be empty")
	}

	for name, workspace := range cfg.Workspaces {
		if strings.TrimSpace(name) == "" {
			problems = append(problems, "workspace name must not be empty")
		}

		displayName := name
		if displayName == "" {
			displayName = "<empty>"
		}

		if strings.TrimSpace(workspace.Root) == "" {
			problems = append(problems, fmt.Sprintf("workspace %q root must not be empty", displayName))
		}

		if len(workspace.Windows) == 0 {
			problems = append(problems, fmt.Sprintf("workspace %q must have at least one window", displayName))
		}

		windowNames := make(map[string]struct{}, len(workspace.Windows))
		for i, window := range workspace.Windows {
			windowName := strings.TrimSpace(window.Name)
			if windowName == "" {
				problems = append(problems, fmt.Sprintf("workspace %q window %d name must not be empty", displayName, i+1))
				continue
			}

			if _, ok := windowNames[windowName]; ok {
				problems = append(problems, fmt.Sprintf("workspace %q has duplicate window name %q", displayName, windowName))
				continue
			}

			windowNames[windowName] = struct{}{}
		}
	}

	if len(problems) > 0 {
		return ValidationError{Problems: problems}
	}

	return nil
}
