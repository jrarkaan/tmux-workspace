package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the twx version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "twx version %s\n", version)
		},
	}
}
