package cmd

import (
	"fmt"
	"os"

	"github.com/jrarkaan/tmux-workspace/internal/config"
	"github.com/spf13/cobra"
)

var configPath string

func Execute() {
	rootCmd := newRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "twx",
		Short: "Declarative tmux workspace manager for Ubuntu",
		Long: banner() + `

twx manages tmux workspaces declaratively.

It is designed for Ubuntu users who want to replace repetitive tmux
session and window shell scripts with a clear YAML configuration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.PersistentFlags().StringVar(
		&configPath,
		"config",
		config.DefaultConfigPath(),
		"config file path; default is ~/.config/twx/config.yaml",
	)

	rootCmd.AddCommand(newVersionCommand())
	rootCmd.AddCommand(newDoctorCommand())
	rootCmd.AddCommand(newConfigCommand())
	rootCmd.AddCommand(newListCommand())
	rootCmd.AddCommand(newSessionsCommand())
	rootCmd.AddCommand(newWindowsCommand())
	rootCmd.AddCommand(newStartCommand())
	rootCmd.AddCommand(newWorkspaceCommand())
	rootCmd.AddCommand(newWindowCommand())
	rootCmd.AddCommand(newTpmCommand())
	rootCmd.AddCommand(newAttachCommand())
	rootCmd.AddCommand(newKillCommand())
	rootCmd.AddCommand(newRestartCommand())

	return rootCmd
}
