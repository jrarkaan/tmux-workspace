package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/jrarkaan/tmux-workspace/internal/config"
	"github.com/spf13/cobra"
)

func newWindowCommand() *cobra.Command {
	windowCmd := &cobra.Command{
		Use:   "window",
		Short: "Manage windows in the twx config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	windowCmd.AddCommand(newWindowAddCommand())
	windowCmd.AddCommand(newWindowRemoveCommand())
	windowCmd.AddCommand(newWindowSetCommandCommand())

	return windowCmd
}

func newWindowAddCommand() *cobra.Command {
	var commandFlag string
	var force bool

	addCmd := &cobra.Command{
		Use:   "add <workspace> <window>",
		Short: "Add a window to a workspace in the twx config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := strings.TrimSpace(args[0])
			windowName := strings.TrimSpace(args[1])
			path := config.ResolveConfigPath(configPath)

			cfg, err := loadConfigForMutation(path)
			if err != nil {
				return err
			}

			if _, ok := config.GetWorkspace(cfg, workspaceName); !ok {
				return fmt.Errorf("workspace not found: %s", workspaceName)
			}

			window := config.Window{
				Name:    windowName,
				Command: commandFlag,
			}

			_, exists := config.GetWindow(cfg, workspaceName, windowName)

			// AddWindow validates implicitly
			if err := config.AddWindow(cfg, workspaceName, window, force); err != nil {
				return err
			}

			backupPath, err := config.SaveWithBackup(path, cfg, time.Now())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()

			if exists {
				fmt.Fprintf(out, "Window replaced: %s\n", windowName)
			} else {
				fmt.Fprintf(out, "Window added: %s\n", windowName)
			}
			fmt.Fprintf(out, "Workspace: %s\n", workspaceName)
			fmt.Fprintf(out, "Existing config backed up: %s\n", backupPath)
			fmt.Fprintf(out, "Config updated: %s\n", path)

			return nil
		},
	}

	addCmd.Flags().StringVar(&commandFlag, "command", "", "initial command for the window")
	addCmd.Flags().BoolVar(&force, "force", false, "replace existing window")

	return addCmd
}

func newWindowRemoveCommand() *cobra.Command {
	var force bool

	removeCmd := &cobra.Command{
		Use:   "remove <workspace> <window>",
		Short: "Remove a window from a workspace in the twx config",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := strings.TrimSpace(args[0])
			windowName := strings.TrimSpace(args[1])
			path := config.ResolveConfigPath(configPath)

			cfg, err := loadConfigForMutation(path)
			if err != nil {
				return err
			}

			if _, ok := config.GetWorkspace(cfg, workspaceName); !ok {
				return fmt.Errorf("workspace not found: %s", workspaceName)
			}

			if _, ok := config.GetWindow(cfg, workspaceName, windowName); !ok {
				return fmt.Errorf("window not found: %s", windowName)
			}

			out := cmd.OutOrStdout()
			if !force {
				fmt.Fprintf(out, "This will remove window from workspace config: %s/%s\n", workspaceName, windowName)
				fmt.Fprintln(out, "Re-run with --force to confirm.")
				return fmt.Errorf("confirmation required: re-run with --force")
			}

			if err := config.RemoveWindow(cfg, workspaceName, windowName); err != nil {
				return err
			}

			backupPath, err := config.SaveWithBackup(path, cfg, time.Now())
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "Window removed: %s\n", windowName)
			fmt.Fprintf(out, "Workspace: %s\n", workspaceName)
			fmt.Fprintf(out, "Existing config backed up: %s\n", backupPath)
			fmt.Fprintf(out, "Config updated: %s\n", path)

			return nil
		},
	}

	removeCmd.Flags().BoolVar(&force, "force", false, "confirm window removal")

	return removeCmd
}

func newWindowSetCommandCommand() *cobra.Command {
	setCommandCmd := &cobra.Command{
		Use:   "set-command <workspace> <window> <command>",
		Short: "Set or update the command for a window in the twx config",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := strings.TrimSpace(args[0])
			windowName := strings.TrimSpace(args[1])
			command := args[2] // Don't trim space for command, user might want leading spaces

			if strings.TrimSpace(command) == "" {
				return fmt.Errorf("command cannot be empty")
			}

			path := config.ResolveConfigPath(configPath)

			cfg, err := loadConfigForMutation(path)
			if err != nil {
				return err
			}

			if _, ok := config.GetWorkspace(cfg, workspaceName); !ok {
				return fmt.Errorf("workspace not found: %s", workspaceName)
			}

			if _, ok := config.GetWindow(cfg, workspaceName, windowName); !ok {
				return fmt.Errorf("window not found: %s", windowName)
			}

			if err := config.SetWindowCommand(cfg, workspaceName, windowName, command); err != nil {
				return err
			}

			backupPath, err := config.SaveWithBackup(path, cfg, time.Now())
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Window command updated: %s\n", windowName)
			fmt.Fprintf(out, "Workspace: %s\n", workspaceName)
			fmt.Fprintf(out, "Existing config backed up: %s\n", backupPath)
			fmt.Fprintf(out, "Config updated: %s\n", path)

			return nil
		},
	}

	return setCommandCmd
}
