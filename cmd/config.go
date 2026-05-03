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
	configCmd.AddCommand(newConfigInitCommand())

	return configCmd
}

func newConfigInitCommand() *cobra.Command {
	var force bool
	var printOnly bool

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create a default twx config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.ResolveConfigPath(configPath)

			result, err := config.InitConfig(path, config.InitOptions{
				Force: force,
				Print: printOnly,
			})
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			if result.Printed {
				content, err := config.DefaultConfigYAML()
				if err != nil {
					return err
				}
				fmt.Fprint(out, string(content))
				return nil
			}

			if result.Existed && !result.Created {
				fmt.Fprintf(out, "Config already exists: %s\n", result.Path)
				fmt.Fprintln(out, "No changes made.")
				return nil
			}

			if result.BackupPath != "" {
				fmt.Fprintf(out, "Existing config backed up: %s\n", result.BackupPath)
			}
			fmt.Fprintf(out, "Config created: %s\n", result.Path)

			return nil
		},
	}

	initCmd.Flags().BoolVar(&force, "force", false, "overwrite existing config after creating a backup")
	initCmd.Flags().BoolVar(&printOnly, "print", false, "print the default config without writing a file")

	return initCmd
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
