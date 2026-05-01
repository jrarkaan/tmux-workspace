package tmux

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(name string, args ...string) ([]byte, error)
}

type ExecRunner struct{}

func (r ExecRunner) Run(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

type Client struct {
	runner   Runner
	lookPath func(file string) (string, error)
}

type WindowInfo struct {
	Index string
	Name  string
}

func NewClient() *Client {
	return NewClientWithRunner(ExecRunner{})
}

func NewClientWithRunner(r Runner) *Client {
	if r == nil {
		r = ExecRunner{}
	}

	return &Client{
		runner:   r,
		lookPath: exec.LookPath,
	}
}

func (c *Client) IsInstalled() bool {
	_, err := c.lookPath("tmux")
	return err == nil
}

func (c *Client) HasSession(name string) (bool, error) {
	if !c.IsInstalled() {
		return false, ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "has-session", "-t", name)
	if err != nil {
		if isSessionNotFound(output, err) {
			return false, nil
		}

		return false, commandError(fmt.Sprintf("check tmux session %q", name), output, err)
	}

	return true, nil
}

func (c *Client) ListSessions() ([]string, error) {
	if !c.IsInstalled() {
		return nil, ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "list-sessions", "-F", "#S")
	if err != nil {
		if isNoServerOrSessions(output, err) {
			return []string{}, nil
		}

		return nil, commandError("list tmux sessions", output, err)
	}

	return parseLines(output), nil
}

func (c *Client) ListWindows(session string) ([]WindowInfo, error) {
	if !c.IsInstalled() {
		return nil, ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "list-windows", "-t", session, "-F", "#I:#W")
	if err != nil {
		if isSessionNotFound(output, err) {
			return nil, fmt.Errorf("%w: %s", ErrSessionNotFound, session)
		}

		return nil, commandError(fmt.Sprintf("list tmux windows for session %q", session), output, err)
	}

	windows := make([]WindowInfo, 0)
	for _, line := range parseLines(output) {
		index, name, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		windows = append(windows, WindowInfo{
			Index: strings.TrimSpace(index),
			Name:  strings.TrimSpace(name),
		})
	}

	return windows, nil
}

func (c *Client) NewSession(session string, window string, root string) error {
	if !c.IsInstalled() {
		return ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "new-session", "-d", "-s", session, "-n", window, "-c", root)
	if err != nil {
		return commandError(fmt.Sprintf("create tmux session %q", session), output, err)
	}

	return nil
}

func (c *Client) NewWindow(session string, window string, root string) error {
	if !c.IsInstalled() {
		return ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "new-window", "-t", session, "-n", window, "-c", root)
	if err != nil {
		return commandError(fmt.Sprintf("create tmux window %q in session %q", window, session), output, err)
	}

	return nil
}

func (c *Client) SendKeys(target string, command string) error {
	if !c.IsInstalled() {
		return ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "send-keys", "-t", target, command, "C-m")
	if err != nil {
		return commandError(fmt.Sprintf("send command to tmux target %q", target), output, err)
	}

	return nil
}

func (c *Client) SelectWindow(session string, window string) error {
	if !c.IsInstalled() {
		return ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "select-window", "-t", targetWindow(session, window))
	if err != nil {
		return commandError(fmt.Sprintf("select tmux window %q in session %q", window, session), output, err)
	}

	return nil
}

func (c *Client) Attach(session string) error {
	if !c.IsInstalled() {
		return ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "attach", "-t", session)
	if err != nil {
		return commandError(fmt.Sprintf("attach tmux session %q", session), output, err)
	}

	return nil
}

func (c *Client) KillSession(session string) error {
	if !c.IsInstalled() {
		return ErrTmuxNotInstalled
	}

	output, err := c.runner.Run("tmux", "kill-session", "-t", session)
	if err != nil {
		return commandError(fmt.Sprintf("kill tmux session %q", session), output, err)
	}

	return nil
}

func parseLines(output []byte) []string {
	rawLines := strings.Split(string(output), "\n")
	lines := make([]string, 0, len(rawLines))

	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lines = append(lines, line)
	}

	return lines
}

func isNoServerOrSessions(output []byte, err error) bool {
	text := errorText(output, err)
	return strings.Contains(text, "no server running") ||
		strings.Contains(text, "failed to connect") ||
		strings.Contains(text, "error connecting") ||
		strings.Contains(text, "no sessions")
}

func isSessionNotFound(output []byte, err error) bool {
	text := errorText(output, err)
	return strings.Contains(text, "can't find session") ||
		strings.Contains(text, "can’t find session") ||
		strings.Contains(text, "session not found") ||
		isNoServerOrSessions(output, err)
}

func errorText(output []byte, err error) string {
	var parts []string
	if len(output) > 0 {
		parts = append(parts, string(output))
	}
	if err != nil {
		parts = append(parts, err.Error())
	}

	text := strings.ToLower(strings.Join(parts, " "))
	return strings.TrimSpace(text)
}

func commandError(operation string, output []byte, err error) error {
	detail := strings.TrimSpace(string(output))
	if detail == "" {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return fmt.Errorf("%s: %s: %w", operation, detail, err)
}

func IsTmuxNotInstalled(err error) bool {
	return errors.Is(err, ErrTmuxNotInstalled)
}

func IsSessionNotFound(err error) bool {
	return errors.Is(err, ErrSessionNotFound)
}

func targetWindow(session string, window string) string {
	return session + ":" + window
}
