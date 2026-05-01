package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/system"
	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check local system readiness for twx",
		Run: func(cmd *cobra.Command, args []string) {
			results := system.RunDoctor(configPath)
			hasWarnings := false

			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "twx doctor")
			fmt.Fprintln(out)

			for _, result := range results {
				if result.Warning {
					hasWarnings = true
				}
				fmt.Fprintf(out, "%-16s : %-16s %s\n", result.Name, result.Status, result.Detail)
			}

			fmt.Fprintln(out)
			if hasWarnings {
				fmt.Fprintln(out, "Result: ready with warnings")
			} else {
				fmt.Fprintln(out, "Result: ready")
			}

			fmt.Fprintln(out)
			fmt.Fprintln(out, "Notes:")
			fmt.Fprintln(out, "  - Missing TPM is okay. It is optional.")
			fmt.Fprintln(out, "  - Missing config is okay before running twx init.")
			fmt.Fprintln(out, "  - Missing ~/.tmux.conf is okay, but recommended for better tmux defaults.")
		},
	}
}
