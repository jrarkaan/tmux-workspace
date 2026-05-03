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

			var configMissing, tpmMissing, tmuxConfMissing bool
			for _, result := range results {
				if result.Warning {
					hasWarnings = true
				}
				if result.Name == "Config file" && result.Status == "optional missing" {
					configMissing = true
				}
				if result.Name == "TPM directory" && result.Status == "optional missing" {
					tpmMissing = true
				}
				if result.Name == "tmux config" && result.Status == "optional missing" {
					tmuxConfMissing = true
				}
				fmt.Fprintf(out, "%-16s : %-16s %s\n", result.Name, result.Status, result.Detail)
			}

			fmt.Fprintln(out)
			if hasWarnings {
				fmt.Fprintln(out, "Result: ready with warnings")
			} else {
				fmt.Fprintln(out, "Result: ready")
			}

			if configMissing || tpmMissing || tmuxConfMissing {
				fmt.Fprintln(out)
				fmt.Fprintln(out, "Notes:")
				if tpmMissing {
					fmt.Fprintln(out, "  - Missing TPM is okay. It is optional.")
				}
				if configMissing {
					fmt.Fprintln(out, "  - Missing config is okay before running twx init.")
				}
				if tmuxConfMissing {
					fmt.Fprintln(out, "  - Missing ~/.tmux.conf is okay, but recommended for better tmux defaults.")
				}
			}
		},
	}
}
