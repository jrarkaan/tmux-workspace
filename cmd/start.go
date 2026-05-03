package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/tmux"
	"github.com/spf13/cobra"
)

func newStartCommand() *cobra.Command {
	var noAttach bool
	var force bool

	startCmd := &cobra.Command{
		Use:   "start <workspace>",
		Short: "Create a tmux workspace from config",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaceName := args[0]
			workspace, _, err := resolveWorkspace(workspaceName)
			if err != nil {
				return err
			}

			client := tmux.NewClient()
			_, err = tmux.StartWorkspace(client, workspaceName, workspace, tmux.WorkspaceStartOptions{
				NoAttach: noAttach,
				Force:    force,
				OnMessage: func(message string) {
					fmt.Fprintln(cmd.OutOrStdout(), message)
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

	startCmd.Flags().BoolVar(&noAttach, "no-attach", false, "create the tmux workspace without attaching")
	startCmd.Flags().BoolVar(&force, "force", false, "recreate the tmux workspace if it already exists")

	return startCmd
}
