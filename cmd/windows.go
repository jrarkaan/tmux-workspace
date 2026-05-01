package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/tmux"
	"github.com/spf13/cobra"
)

func newWindowsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "windows <session>",
		Short: "List windows in a tmux session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			session := args[0]
			client := tmux.NewClient()

			windows, err := client.ListWindows(session)
			if err != nil {
				if tmux.IsTmuxNotInstalled(err) {
					return tmux.ErrTmuxNotInstalled
				}
				if tmux.IsSessionNotFound(err) {
					return fmt.Errorf("%w: %s", tmux.ErrSessionNotFound, session)
				}

				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Windows in session %s:\n", session)
			fmt.Fprintln(out)
			for _, window := range windows {
				fmt.Fprintf(out, "  %s: %s\n", window.Index, window.Name)
			}

			return nil
		},
	}
}
