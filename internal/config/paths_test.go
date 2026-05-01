package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	path := DefaultConfigPath()
	if path == "" {
		t.Fatal("DefaultConfigPath returned an empty path")
	}

	expectedSuffix := filepath.Join(".config", "twx", "config.yaml")
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Fatalf("DefaultConfigPath() = %q, want suffix %q", path, expectedSuffix)
	}

	if !strings.Contains(path, "twx") {
		t.Fatalf("DefaultConfigPath() = %q, want path to contain twx", path)
	}
}
