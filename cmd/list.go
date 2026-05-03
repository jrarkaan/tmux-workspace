package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jrarkaan/tmux-workspace/internal/config"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.ResolveConfigPath(configPath)

			cfg, err := config.Load(path)
			if err != nil {
				return err
			}

			if err := config.Validate(cfg); err != nil {
				return err
			}

			printWorkspaceList(cmd, cfg)

			return nil
		},
	}
}

func printWorkspaceList(cmd *cobra.Command, cfg *config.Config) {
	names := make([]string, 0, len(cfg.Workspaces))
	for name := range cfg.Workspaces {
		names = append(names, name)
	}
	sort.Strings(names)

	out := cmd.OutOrStdout()
	if len(names) == 0 {
		fmt.Fprintln(out, "No workspaces configured.")
		fmt.Fprintln(out, "Add one by editing your config or using a future workspace command.")
		return
	}

	fmt.Fprintln(out, "Available workspaces:")
	fmt.Fprintln(out)

	for _, name := range names {
		workspace := cfg.Workspaces[name]
		windowNames := make([]string, 0, len(workspace.Windows))
		for _, window := range workspace.Windows {
			windowNames = append(windowNames, window.Name)
		}

		fmt.Fprintf(out, "  %s\n", name)
		fmt.Fprintf(out, "    root    : %s\n", workspace.Root)
		fmt.Fprintf(out, "    windows : %s\n", strings.Join(windowNames, ", "))
		fmt.Fprintln(out)
	}
}
