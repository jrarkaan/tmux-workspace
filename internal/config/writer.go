package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func Save(path string, cfg *Config) error {
	resolvedPath := ExpandHome(path)
	normalizeConfig(cfg)

	if err := Validate(cfg); err != nil {
		return err
	}

	content, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(resolvedPath), 0o755); err != nil {
		return fmt.Errorf("create config directory %s: %w", filepath.Dir(resolvedPath), err)
	}

	if err := os.WriteFile(resolvedPath, content, 0o644); err != nil {
		return fmt.Errorf("write config file %s: %w", resolvedPath, err)
	}

	return nil
}

func Backup(path string, now time.Time) (string, error) {
	resolvedPath := ExpandHome(path)
	backupPath := uniqueBackupPath(backupConfigPath(resolvedPath, now))

	content, err := os.ReadFile(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("read existing config %s: %w", resolvedPath, err)
	}

	if err := os.WriteFile(backupPath, content, 0o644); err != nil {
		return "", fmt.Errorf("write config backup %s: %w", backupPath, err)
	}

	return backupPath, nil
}

func SaveWithBackup(path string, cfg *Config, now time.Time) (string, error) {
	backupPath, err := Backup(path, now)
	if err != nil {
		return "", err
	}

	if err := Save(path, cfg); err != nil {
		return "", err
	}

	return backupPath, nil
}

func normalizeConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	if cfg.Workspaces == nil {
		cfg.Workspaces = map[string]Workspace{}
	}
}

func uniqueBackupPath(path string) string {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	for i := 1; ; i++ {
		candidate := fmt.Sprintf("%s.%d", path, i)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
