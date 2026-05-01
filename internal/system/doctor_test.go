package system

import (
	"path/filepath"
	"testing"
)

func TestRunDoctorReturnsExpectedChecks(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")

	results := RunDoctor(configPath)
	if len(results) == 0 {
		t.Fatal("RunDoctor returned no results")
	}

	byName := make(map[string]CheckResult, len(results))
	for _, result := range results {
		if result.Name == "" {
			t.Fatalf("RunDoctor returned result with empty name: %#v", result)
		}
		if result.Status == "" {
			t.Fatalf("RunDoctor returned result with empty status: %#v", result)
		}
		if result.Detail == "" {
			t.Fatalf("RunDoctor returned result with empty detail: %#v", result)
		}

		byName[result.Name] = result
	}

	expectedNames := []string{
		"OS",
		"Ubuntu version",
		"tmux",
		"git",
		"Shell",
		"Config file",
		"tmux config",
		"TPM directory",
	}

	for _, name := range expectedNames {
		if _, ok := byName[name]; !ok {
			t.Fatalf("RunDoctor missing check %q; got %#v", name, results)
		}
	}

	if got := byName["Config file"].Detail; got != configPath {
		t.Fatalf("Config file detail = %q, want %q", got, configPath)
	}
}
