package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/config"
	"github.com/jrarkaan/tmux-workspace/internal/tmux"
	"github.com/spf13/cobra"
)

func newAttachCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "attach <workspace>",
		Short: "Attach to a running tmux workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := args[0]
			if _, _, err := resolveWorkspace(workspaceName); err != nil {
				return err
			}

			client := tmux.NewClient()
			exists, err := client.HasSession(workspaceName)
			if err != nil {
				if tmux.IsTmuxNotInstalled(err) {
					return tmux.ErrTmuxNotInstalled
				}
				return err
			}
			if !exists {
				return fmt.Errorf("tmux session not found: %s\nRun: twx start %s", workspaceName, workspaceName)
			}

			return client.Attach(workspaceName)
		},
	}
}

func newKillCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "kill <workspace>",
		Short: "Kill a running tmux workspace session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := args[0]
			if _, _, err := resolveWorkspace(workspaceName); err != nil {
				return err
			}

			client := tmux.NewClient()
			exists, err := client.HasSession(workspaceName)
			if err != nil {
				if tmux.IsTmuxNotInstalled(err) {
					return tmux.ErrTmuxNotInstalled
				}
				return err
			}
			if !exists {
				return fmt.Errorf("tmux session not found: %s", workspaceName)
			}
			if err := client.KillSession(workspaceName); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Session killed: %s\n", workspaceName)
			return nil
		},
	}
}

func newRestartCommand() *cobra.Command {
	var noAttach bool

	restartCmd := &cobra.Command{
		Use:   "restart <workspace>",
		Short: "Recreate a tmux workspace session from config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := args[0]
			workspace, _, err := resolveWorkspace(workspaceName)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			client := tmux.NewClient()
			_, err = tmux.RestartWorkspace(client, workspaceName, workspace, tmux.WorkspaceRestartOptions{
				NoAttach: noAttach,
				OnMessage: func(message string) {
					fmt.Fprintln(out, message)
				},
			})
			if err != nil {
				if tmux.IsTmuxNotInstalled(err) {
					return tmux.ErrTmuxNotInstalled
				}
				return err
			}

			return nil
		},
	}

	restartCmd.Flags().BoolVar(&noAttach, "no-attach", false, "restart the tmux workspace without attaching")

	return restartCmd
}

func resolveWorkspace(workspaceName string) (config.Workspace, string, error) {
	path := config.ResolveConfigPath(configPath)

	cfg, err := config.Load(path)
	if err != nil {
		return config.Workspace{}, path, err
	}
	if err := config.Validate(cfg); err != nil {
		return config.Workspace{}, path, err
	}

	workspace, ok := config.GetWorkspace(cfg, workspaceName)
	if !ok {
		return config.Workspace{}, path, fmt.Errorf("workspace not found: %s", workspaceName)
	}

	return workspace, path, nil
}
