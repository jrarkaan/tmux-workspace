package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type InitOptions struct {
	Force bool
	Print bool
	Now   func() time.Time
}

type InitResult struct {
	Path       string
	BackupPath string
	Created    bool
	Existed    bool
	Printed    bool
}

func DefaultConfig() *Config {
	return &Config{
		Version: 1,
		Defaults: Defaults{
			Attach:          true,
			CreateIfMissing: true,
			BaseIndex:       1,
		},
		Workspaces: map[string]Workspace{},
	}
}

func DefaultConfigYAML() ([]byte, error) {
	return []byte(`version: 1

defaults:
  attach: true
  create_if_missing: true
  base_index: 1

workspaces: {}
`), nil
}

func InitConfig(path string, opts InitOptions) (InitResult, error) {
	resolvedPath := ExpandHome(path)
	result := InitResult{Path: resolvedPath}

	if opts.Print {
		result.Printed = true
		return result, nil
	}

	content, err := DefaultConfigYAML()
	if err != nil {
		return result, err
	}

	info, err := os.Stat(resolvedPath)
	if err != nil && !os.IsNotExist(err) {
		return result, fmt.Errorf("check config file %s: %w", resolvedPath, err)
	}

	exists := err == nil
	if exists && info.IsDir() {
		return result, fmt.Errorf("config path is a directory: %s", resolvedPath)
	}

	if exists {
		result.Existed = true
		if !opts.Force {
			return result, nil
		}

		backupPath, err := Backup(resolvedPath, now(opts))
		if err != nil {
			return result, err
		}
		result.BackupPath = backupPath
	}

	if err := os.MkdirAll(filepath.Dir(resolvedPath), 0o755); err != nil {
		return result, fmt.Errorf("create config directory %s: %w", filepath.Dir(resolvedPath), err)
	}

	if err := os.WriteFile(resolvedPath, content, 0o644); err != nil {
		return result, fmt.Errorf("write config file %s: %w", resolvedPath, err)
	}

	result.Created = true
	return result, nil
}

func backupConfigPath(path string, now time.Time) string {
	return fmt.Sprintf("%s.bak.%s", path, now.Format("20060102-150405"))
}

func now(opts InitOptions) time.Time {
	if opts.Now != nil {
		return opts.Now()
	}

	return time.Now()
}
