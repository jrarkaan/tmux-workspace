package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/tpm"
	"github.com/spf13/cobra"
)

func newTpmCommand() *cobra.Command {
	tpmCmd := &cobra.Command{
		Use:   "tpm",
		Short: "Manage tmux TPM installation and config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	tpmCmd.AddCommand(newTpmStatusCommand())
	tpmCmd.AddCommand(newTpmInstallCommand())

	return tpmCmd
}

func newTpmStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show TPM status and configuration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := tpm.NewManagerDefault()
			status, err := manager.GetStatus()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "TPM status")
			fmt.Fprintln(out)

			fmt.Fprintln(out, "System")
			if status.TmuxInstalled {
				fmt.Fprintf(out, "  tmux : installed, %s\n", status.TmuxVersion)
			} else {
				fmt.Fprintln(out, "  tmux : missing")
			}
			if status.GitInstalled {
				fmt.Fprintf(out, "  git  : installed, %s\n", status.GitVersion)
			} else {
				fmt.Fprintln(out, "  git  : missing")
			}
			fmt.Fprintln(out)

			fmt.Fprintln(out, "Files")
			if status.TPMDirExists {
				fmt.Fprintf(out, "  TPM directory : found, %s\n", status.TPMDir)
			} else {
				fmt.Fprintf(out, "  TPM directory : missing, %s\n", status.TPMDir)
			}
			if status.TmuxConfigExists {
				fmt.Fprintf(out, "  tmux config   : found, %s\n", status.TmuxConfigPath)
			} else {
				fmt.Fprintf(out, "  tmux config   : missing, %s\n", status.TmuxConfigPath)
			}
			fmt.Fprintln(out)

			fmt.Fprintln(out, "Config")
			printConfigState := func(name string, has bool) {
				if has {
					fmt.Fprintf(out, "  %-16s : configured\n", name)
				} else {
					fmt.Fprintf(out, "  %-16s : missing\n", name)
				}
			}
			printConfigState("tpm plugin", status.HasTPMPlugin)
			printConfigState("resurrect plugin", status.HasResurrectPlugin)
			printConfigState("continuum plugin", status.HasContinuumPlugin)
			printConfigState("continuum restore", status.HasContinuumRestore)
			printConfigState("tpm run line", status.HasTPMRunLine)
			fmt.Fprintln(out)

			if status.TPMDirExists && status.HasTPMPlugin && status.HasTPMRunLine {
				fmt.Fprintln(out, "Result: ready")
			} else {
				fmt.Fprintln(out, "Result: TPM is not installed")
				fmt.Fprintln(out, "Next step:")
				fmt.Fprintln(out, "  twx tpm install")
			}

			return nil
		},
	}
}

func newTpmInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install",
		Short: "Install TPM and configure ~/.tmux.conf",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := tpm.NewManagerDefault()

			backupPath, mutated, err := manager.Install()
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if !mutated {
				fmt.Fprintln(out, "TPM config already present. No changes made.")
				return nil
			}

			fmt.Fprintln(out, "TPM installed or already present")

			status, _ := manager.GetStatus()
			fmt.Fprintf(out, "tmux config updated: %s\n", status.TmuxConfigPath)
			if backupPath != "" {
				fmt.Fprintf(out, "Existing tmux config backed up: %s\n", backupPath)
			}
			fmt.Fprintln(out, "Next step:")
			fmt.Fprintln(out, "  Open tmux and press your prefix then I to install plugins.")
			fmt.Fprintln(out, "  Default tmux prefix is Ctrl+b. If you configured twx docs, it may be Ctrl+a.")

			return nil
		},
	}
}
