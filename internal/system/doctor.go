package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type CheckResult struct {
	Name    string
	Status  string
	Detail  string
	Warning bool
}

func RunDoctor(configPath string) []CheckResult {
	home := userHome()

	return []CheckResult{
		runtimeCheck(),
		ubuntuCheck(),
		commandCheck("tmux", "tmux", "-V", true),
		commandCheck("git", "git", "--version", true),
		shellCheck(),
		fileCheck("Config file", configPath, true),
		fileCheck("tmux config", filepath.Join(home, ".tmux.conf"), true),
		dirCheck("TPM directory", filepath.Join(home, ".tmux", "plugins", "tpm"), true),
	}
}

func runtimeCheck() CheckResult {
	return CheckResult{
		Name:   "OS",
		Status: "info",
		Detail: runtime.GOOS + "/" + runtime.GOARCH,
	}
}

func ubuntuCheck() CheckResult {
	values, err := readOSRelease("/etc/os-release")
	if err != nil {
		return CheckResult{
			Name:    "Ubuntu version",
			Status:  "warning",
			Detail:  "unable to read /etc/os-release",
			Warning: true,
		}
	}

	prettyName := values["PRETTY_NAME"]
	if prettyName == "" {
		prettyName = "unknown version"
	}

	if values["ID"] != "ubuntu" {
		return CheckResult{
			Name:    "Ubuntu version",
			Status:  "warning",
			Detail:  prettyName,
			Warning: true,
		}
	}

	return CheckResult{
		Name:   "Ubuntu version",
		Status: "ok",
		Detail: prettyName,
	}
}

func commandCheck(name, binary, versionArg string, warning bool) CheckResult {
	if _, err := exec.LookPath(binary); err != nil {
		return CheckResult{
			Name:    name,
			Status:  "warning missing",
			Detail:  binary + " not found in PATH",
			Warning: warning,
		}
	}

	output, err := exec.Command(binary, versionArg).CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:    name,
			Status:  "warning",
			Detail:  strings.TrimSpace(string(output)),
			Warning: warning,
		}
	}

	return CheckResult{
		Name:   name,
		Status: "ok",
		Detail: strings.TrimSpace(string(output)),
	}
}

func shellCheck() CheckResult {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return CheckResult{
			Name:    "Shell",
			Status:  "warning",
			Detail:  "SHELL environment variable is not set",
			Warning: true,
		}
	}

	return CheckResult{
		Name:   "Shell",
		Status: "ok",
		Detail: shell,
	}
}

func fileCheck(name, path string, optional bool) CheckResult {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			status := "missing"
			if optional {
				status = "optional missing"
			}

			return CheckResult{
				Name:    name,
				Status:  status,
				Detail:  path,
				Warning: optional,
			}
		}

		return CheckResult{
			Name:    name,
			Status:  "warning",
			Detail:  err.Error(),
			Warning: true,
		}
	}

	if info.IsDir() {
		return CheckResult{
			Name:    name,
			Status:  "warning",
			Detail:  path + " is a directory",
			Warning: true,
		}
	}

	return CheckResult{Name: name, Status: "ok", Detail: path}
}

func dirCheck(name, path string, optional bool) CheckResult {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			status := "missing"
			if optional {
				status = "optional missing"
			}

			return CheckResult{
				Name:    name,
				Status:  status,
				Detail:  path,
				Warning: optional,
			}
		}

		return CheckResult{
			Name:    name,
			Status:  "warning",
			Detail:  err.Error(),
			Warning: true,
		}
	}

	if !info.IsDir() {
		return CheckResult{
			Name:    name,
			Status:  "warning",
			Detail:  path + " is not a directory",
			Warning: true,
		}
	}

	return CheckResult{Name: name, Status: "ok", Detail: path}
}

func readOSRelease(path string) (map[string]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	values := make(map[string]string)
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		values[key] = strings.Trim(strings.TrimSpace(value), `"`)
	}

	return values, nil
}

func userHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}

	return home
}
