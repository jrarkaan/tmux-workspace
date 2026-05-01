package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/tmux"
	"github.com/spf13/cobra"
)

func newSessionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sessions",
		Short: "List active tmux sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := tmux.NewClient()
			sessions, err := client.ListSessions()
			if err != nil {
				if tmux.IsTmuxNotInstalled(err) {
					return tmux.ErrTmuxNotInstalled
				}

				return err
			}

			out := cmd.OutOrStdout()
			if len(sessions) == 0 {
				fmt.Fprintln(out, "No active tmux sessions found.")
				return nil
			}

			fmt.Fprintln(out, "Active tmux sessions:")
			fmt.Fprintln(out)
			for _, session := range sessions {
				fmt.Fprintf(out, "  %s\n", session)
			}

			return nil
		},
	}
}
