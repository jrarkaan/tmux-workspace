package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the twx version",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "twx version %s\n", version)
			fmt.Fprintf(out, "commit: %s\n", commit)
			fmt.Fprintf(out, "date: %s\n", date)
		},
	}
}
