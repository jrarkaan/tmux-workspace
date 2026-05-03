package tpm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Status struct {
	TmuxInstalled       bool
	TmuxVersion         string
	GitInstalled        bool
	GitVersion          string
	TPMDir              string
	TPMDirExists        bool
	TmuxConfigPath      string
	TmuxConfigExists    bool
	HasTPMPlugin        bool
	HasResurrectPlugin  bool
	HasContinuumPlugin  bool
	HasContinuumRestore bool
	HasTPMRunLine       bool
}

type Runner interface {
	Run(name string, args ...string) ([]byte, error)
}

type ExecRunner struct{}

func (r ExecRunner) Run(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

type Manager struct {
	runner   Runner
	homeDir  string
	now      func() time.Time
	lookPath func(file string) (string, error)
}

func NewManagerDefault() *Manager {
	home, _ := os.UserHomeDir()
	if home == "" {
		home = "."
	}
	return NewManager(ExecRunner{}, home)
}

func NewManager(runner Runner, homeDir string) *Manager {
	return &Manager{
		runner:   runner,
		homeDir:  homeDir,
		now:      time.Now,
		lookPath: exec.LookPath,
	}
}

func (m *Manager) GetStatus() (Status, error) {
	s := Status{
		TPMDir:         filepath.Join(m.homeDir, ".tmux", "plugins", "tpm"),
		TmuxConfigPath: filepath.Join(m.homeDir, ".tmux.conf"),
	}

	if _, err := m.lookPath("tmux"); err == nil {
		s.TmuxInstalled = true
		out, _ := m.runner.Run("tmux", "-V")
		s.TmuxVersion = strings.TrimSpace(string(out))
	}

	if _, err := m.lookPath("git"); err == nil {
		s.GitInstalled = true
		out, _ := m.runner.Run("git", "--version")
		s.GitVersion = strings.TrimSpace(string(out))
	}

	if info, err := os.Stat(s.TPMDir); err == nil && info.IsDir() {
		s.TPMDirExists = true
	}

	if info, err := os.Stat(s.TmuxConfigPath); err == nil && !info.IsDir() {
		s.TmuxConfigExists = true
		content, _ := os.ReadFile(s.TmuxConfigPath)
		s.parseConfig(string(content))
	}

	return s, nil
}

func (s *Status) parseConfig(content string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "tmux-plugins/tpm") && strings.HasPrefix(line, "set") {
			s.HasTPMPlugin = true
		}
		if strings.Contains(line, "tmux-plugins/tmux-resurrect") && strings.HasPrefix(line, "set") {
			s.HasResurrectPlugin = true
		}
		if strings.Contains(line, "tmux-plugins/tmux-continuum") && strings.HasPrefix(line, "set") {
			s.HasContinuumPlugin = true
		}
		if strings.Contains(line, "@continuum-restore") && strings.HasPrefix(line, "set") {
			s.HasContinuumRestore = true
		}
		if strings.Contains(line, "run") && strings.Contains(line, "tpm") && !strings.HasPrefix(line, "#") {
			s.HasTPMRunLine = true
		}
	}
}

const managedBlock = `
# twx TPM managed block
set -g @plugin 'tmux-plugins/tpm'
set -g @plugin 'tmux-plugins/tmux-resurrect'
set -g @plugin 'tmux-plugins/tmux-continuum'
set -g @continuum-restore 'on'
set -g @continuum-save-interval '15'
run '~/.tmux/plugins/tpm/tpm'
# end twx TPM managed block
`

func (m *Manager) Install() (string, bool, error) {
	if _, err := m.lookPath("git"); err != nil {
		return "", false, fmt.Errorf("git is not installed or not in PATH")
	}

	tpmDir := filepath.Join(m.homeDir, ".tmux", "plugins", "tpm")
	pluginsDir := filepath.Join(m.homeDir, ".tmux", "plugins")

	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return "", false, fmt.Errorf("create plugins directory: %w", err)
	}

	if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
		if _, err := m.runner.Run("git", "clone", "https://github.com/tmux-plugins/tpm", tpmDir); err != nil {
			return "", false, fmt.Errorf("clone TPM: %w", err)
		}
	}

	configPath := filepath.Join(m.homeDir, ".tmux.conf")
	var content string

	info, err := os.Stat(configPath)
	if err == nil && !info.IsDir() {
		b, _ := os.ReadFile(configPath)
		content = string(b)
	}

	if strings.Contains(content, "# twx TPM managed block") {
		return "", false, nil
	}

	var backupPath string
	if content != "" {
		timestamp := m.now().Format("20060102-150405")
		backupPath = fmt.Sprintf("%s.bak.%s", configPath, timestamp)
		if err := os.WriteFile(backupPath, []byte(content), 0644); err != nil {
			return "", false, fmt.Errorf("backup existing config: %w", err)
		}
	}

	newContent := content
	if newContent != "" && !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += managedBlock

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return "", false, fmt.Errorf("write config: %w", err)
	}

	return backupPath, true, nil
}
