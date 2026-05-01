package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	resolvedPath := ExpandHome(path)

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s", resolvedPath)
		}

		return nil, fmt.Errorf("read config file %s: %w", resolvedPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("invalid YAML in config file %s: %w", resolvedPath, err)
	}

	return &cfg, nil
}

func ResolveConfigPath(override string) string {
	if override != "" {
		return ExpandHome(override)
	}

	return DefaultConfigPath()
}
