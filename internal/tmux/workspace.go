package tmux

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jrarkaan/tmux-workspace/internal/config"
)

type WorkspaceStartOptions struct {
	NoAttach  bool
	Force     bool
	OnMessage func(string)
}

type WorkspaceStartResult struct {
	Messages []string
}

type WorkspaceRestartOptions struct {
	NoAttach  bool
	OnMessage func(string)
}

func StartWorkspace(client *Client, name string, workspace config.Workspace, opts WorkspaceStartOptions) (*WorkspaceStartResult, error) {
	result := &WorkspaceStartResult{}

	exists, err := client.HasSession(name)
	if err != nil {
		return result, err
	}

	if exists {
		if !opts.Force {
			result.addMessage(opts, "Session already exists: "+name)
			if opts.NoAttach {
				result.addMessage(opts, "Skipping attach because --no-attach was provided.")
				return result, nil
			}

			return result, client.Attach(name)
		}

		if err := client.KillSession(name); err != nil {
			return result, err
		}
	}

	root, messages := resolveWorkspaceRoot(workspace.Root)
	for _, message := range messages {
		result.addMessage(opts, message)
	}

	firstWindow := workspace.Windows[0]
	if err := client.NewSession(name, firstWindow.Name, root); err != nil {
		return result, err
	}
	if err := sendWindowCommand(client, name, firstWindow, workspace.Env); err != nil {
		return result, err
	}

	for _, window := range workspace.Windows[1:] {
		if err := client.NewWindow(name, window.Name, root); err != nil {
			return result, err
		}
		if err := sendWindowCommand(client, name, window, workspace.Env); err != nil {
			return result, err
		}
	}

	if err := client.SelectWindow(name, firstWindow.Name); err != nil {
		return result, err
	}

	if !opts.NoAttach {
		if err := client.Attach(name); err != nil {
			return result, err
		}
	}

	return result, nil
}

func RestartWorkspace(client *Client, name string, workspace config.Workspace, opts WorkspaceRestartOptions) (*WorkspaceStartResult, error) {
	result := &WorkspaceStartResult{}
	result.addMessage(WorkspaceStartOptions{OnMessage: opts.OnMessage}, "Restarting workspace: "+name)

	startResult, err := StartWorkspace(client, name, workspace, WorkspaceStartOptions{
		NoAttach:  true,
		Force:     true,
		OnMessage: opts.OnMessage,
	})
	result.Messages = append(result.Messages, startResult.Messages...)
	if err != nil {
		return result, err
	}

	result.addMessage(WorkspaceStartOptions{OnMessage: opts.OnMessage}, "Session recreated: "+name)
	if opts.NoAttach {
		result.addMessage(WorkspaceStartOptions{OnMessage: opts.OnMessage}, "Skipping attach because --no-attach was provided.")
		return result, nil
	}

	return result, client.Attach(name)
}

func (r *WorkspaceStartResult) addMessage(opts WorkspaceStartOptions, message string) {
	r.Messages = append(r.Messages, message)
	if opts.OnMessage != nil {
		opts.OnMessage(message)
	}
}

func resolveWorkspaceRoot(root string) (string, []string) {
	expanded := config.ExpandHome(root)
	if info, err := os.Stat(expanded); err == nil && info.IsDir() {
		return expanded, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return home, []string{fmt.Sprintf("Workspace root does not exist, using home: %s", root)}
}

func sendWindowCommand(client *Client, session string, window config.Window, env map[string]string) error {
	if strings.TrimSpace(window.Command) == "" {
		return nil
	}

	return client.SendKeys(targetWindow(session, window.Name), commandWithEnv(window.Command, env))
}

func commandWithEnv(command string, env map[string]string) string {
	if len(env) == 0 {
		return command
	}

	keys := make([]string, 0, len(env))
	for key := range env {
		if strings.TrimSpace(key) == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	if len(keys) == 0 {
		return command
	}

	parts := make([]string, 0, len(keys)+1)
	for _, key := range keys {
		// TODO: move env handling to tmux environment primitives when lifecycle commands mature.
		parts = append(parts, fmt.Sprintf("export %s='%s'", key, shellSingleQuote(env[key])))
	}
	parts = append(parts, command)

	return strings.Join(parts, "; ")
}

func shellSingleQuote(value string) string {
	return strings.ReplaceAll(value, "'", `'\''`)
}
