package config

import (
	"os"
	"path/filepath"
)

func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".twx", "config.yaml")
	}

	return filepath.Join(home, ".config", "twx", "config.yaml")
}
