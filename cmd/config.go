package cmd

import (
	"fmt"

	"github.com/jrarkaan/tmux-workspace/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect and validate twx config",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	configCmd.AddCommand(newConfigPathCommand())
	configCmd.AddCommand(newConfigValidateCommand())

	return configCmd
}

func newConfigPathCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the resolved config path",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), config.ResolveConfigPath(configPath))
		},
	}
}

func newConfigValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the twx config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.ResolveConfigPath(configPath)

			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			if err := config.Validate(cfg); err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Config is valid: %s\n", path)
			fmt.Fprintf(out, "Workspaces: %d\n", len(cfg.Workspaces))

			return nil
		},
	}
}
