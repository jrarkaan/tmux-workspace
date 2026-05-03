package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jrarkaan/tmux-workspace/internal/config"
	"github.com/spf13/cobra"
)

func newWorkspaceCommand() *cobra.Command {
	workspaceCmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage workspaces in the twx config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	workspaceCmd.AddCommand(newWorkspaceAddCommand())
	workspaceCmd.AddCommand(newWorkspaceShowCommand())
	workspaceCmd.AddCommand(newWorkspaceRemoveCommand())

	return workspaceCmd
}

func newWorkspaceAddCommand() *cobra.Command {
	var root string
	var windowsCSV string
	var envEntries []string
	var commandEntries []string
	var force bool

	addCmd := &cobra.Command{
		Use:   "add <workspace>",
		Short: "Add a workspace to the twx config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.TrimSpace(args[0])
			path := config.ResolveConfigPath(configPath)

			cfg, err := loadConfigForMutation(path)
			if err != nil {
				return err
			}

			if strings.TrimSpace(root) == "" {
				return fmt.Errorf("--root is required")
			}

			windows, err := parseWindows(windowsCSV)
			if err != nil {
				return err
			}

			env, err := parseEnvEntries(envEntries)
			if err != nil {
				return err
			}

			commands, err := parseCommandEntries(commandEntries, windows)
			if err != nil {
				return err
			}

			workspace := config.Workspace{
				Root:    root,
				Env:     env,
				Windows: windows,
			}
			for i := range workspace.Windows {
				workspace.Windows[i].Command = commands[workspace.Windows[i].Name]
			}

			_, exists := config.GetWorkspace(cfg, name)
			if exists && !force {
				return fmt.Errorf("workspace already exists: %s\nUse --force to replace it.", name)
			}

			if err := config.AddWorkspace(cfg, name, workspace, force); err != nil {
				return err
			}

			backupPath, err := config.SaveWithBackup(path, cfg, time.Now())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if exists {
				fmt.Fprintf(out, "Workspace replaced: %s\n", name)
				fmt.Fprintf(out, "Existing config backed up: %s\n", backupPath)
			} else {
				fmt.Fprintf(out, "Workspace added: %s\n", name)
			}
			fmt.Fprintf(out, "Config updated: %s\n", path)

			return nil
		},
	}

	addCmd.Flags().StringVar(&root, "root", "", "workspace root directory")
	addCmd.Flags().StringVar(&windowsCSV, "windows", "", "comma-separated workspace window names")
	addCmd.Flags().StringArrayVar(&envEntries, "env", nil, "workspace environment variable as KEY=VALUE")
	addCmd.Flags().StringArrayVar(&commandEntries, "command", nil, "window command as WINDOW=COMMAND")
	addCmd.Flags().BoolVar(&force, "force", false, "replace existing workspace")

	return addCmd
}

func newWorkspaceShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <workspace>",
		Short: "Show a workspace from the twx config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.TrimSpace(args[0])
			path := config.ResolveConfigPath(configPath)

			cfg, err := config.Load(path)
			if err != nil {
				return err
			}
			if err := config.Validate(cfg); err != nil {
				return err
			}

			workspace, ok := config.GetWorkspace(cfg, name)
			if !ok {
				return fmt.Errorf("workspace not found: %s", name)
			}

			printWorkspace(cmd, name, workspace)
			return nil
		},
	}
}

func newWorkspaceRemoveCommand() *cobra.Command {
	var force bool

	removeCmd := &cobra.Command{
		Use:   "remove <workspace>",
		Short: "Remove a workspace from the twx config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := strings.TrimSpace(args[0])
			path := config.ResolveConfigPath(configPath)

			cfg, err := loadConfigForMutation(path)
			if err != nil {
				return err
			}
			if _, ok := config.GetWorkspace(cfg, name); !ok {
				return fmt.Errorf("workspace not found: %s", name)
			}

			out := cmd.OutOrStdout()
			if !force {
				fmt.Fprintf(out, "This will remove workspace from config: %s\n", name)
				fmt.Fprintln(out, "Re-run with --force to confirm.")
				return fmt.Errorf("confirmation required: re-run with --force")
			}

			if err := config.RemoveWorkspace(cfg, name); err != nil {
				return err
			}

			backupPath, err := config.SaveWithBackup(path, cfg, time.Now())
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "Workspace removed: %s\n", name)
			fmt.Fprintf(out, "Existing config backed up: %s\n", backupPath)
			fmt.Fprintf(out, "Config updated: %s\n", path)
			return nil
		},
	}

	removeCmd.Flags().BoolVar(&force, "force", false, "confirm workspace removal")

	return removeCmd
}

func loadConfigForMutation(path string) (*config.Config, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("Config file does not exist: %s\nRun: twx config init", path)
		}
		return nil, err
	}

	cfg, err := config.Load(path)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseWindows(windowsCSV string) ([]config.Window, error) {
	if strings.TrimSpace(windowsCSV) == "" {
		return nil, fmt.Errorf("--windows is required")
	}

	seen := map[string]struct{}{}
	parts := strings.Split(windowsCSV, ",")
	windows := make([]config.Window, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			return nil, fmt.Errorf("window name must not be empty")
		}
		if _, ok := seen[name]; ok {
			return nil, fmt.Errorf("duplicate window name: %s", name)
		}
		seen[name] = struct{}{}
		windows = append(windows, config.Window{Name: name})
	}

	return windows, nil
}

func parseEnvEntries(entries []string) (map[string]string, error) {
	if len(entries) == 0 {
		return nil, nil
	}

	env := make(map[string]string, len(entries))
	for _, entry := range entries {
		key, value, ok := strings.Cut(entry, "=")
		key = strings.TrimSpace(key)
		if !ok {
			return nil, fmt.Errorf("invalid --env %q: expected KEY=VALUE", entry)
		}
		if key == "" {
			return nil, fmt.Errorf("invalid --env %q: key must not be empty", entry)
		}
		env[key] = value
	}

	return env, nil
}

func parseCommandEntries(entries []string, windows []config.Window) (map[string]string, error) {
	commands := make(map[string]string, len(entries))
	windowNames := make(map[string]struct{}, len(windows))
	for _, window := range windows {
		windowNames[window.Name] = struct{}{}
	}

	for _, entry := range entries {
		window, command, ok := strings.Cut(entry, "=")
		window = strings.TrimSpace(window)
		if !ok {
			return nil, fmt.Errorf("invalid --command %q: expected WINDOW=COMMAND", entry)
		}
		if window == "" {
			return nil, fmt.Errorf("invalid --command %q: window must not be empty", entry)
		}
		if _, ok := windowNames[window]; !ok {
			return nil, fmt.Errorf("command references unknown window: %s", window)
		}
		commands[window] = command
	}

	return commands, nil
}

func printWorkspace(cmd *cobra.Command, name string, workspace config.Workspace) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "Workspace: %s\n", name)
	fmt.Fprintf(out, "Root     : %s\n", workspace.Root)
	fmt.Fprintln(out)

	if len(workspace.Env) == 0 {
		fmt.Fprintln(out, "Environment: none")
	} else {
		fmt.Fprintln(out, "Environment:")
		keys := make([]string, 0, len(workspace.Env))
		for key := range workspace.Env {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintf(out, "  %s=%s\n", key, workspace.Env[key])
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Windows:")
	for i, window := range workspace.Windows {
		fmt.Fprintf(out, "  %d. %s\n", i+1, window.Name)
		if window.Command != "" {
			fmt.Fprintf(out, "     command: %s\n", window.Command)
		}
	}
}
