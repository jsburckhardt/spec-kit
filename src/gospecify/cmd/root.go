// Package cmd provides the CLI commands for gospecify
package cmd

import (
	"fmt"

	"github.com/github/spec-kit/gospecify/internal/config"
	"github.com/spf13/cobra"
)

// Execute runs the root command
func Execute() error {
	return NewRootCmd().Execute()
}

// NewRootCmd creates the root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gospecify",
		Short: "Setup tool for Specify spec-driven development projects",
		Long: fmt.Sprintf(`%s

%s

GitHub Spec Kit - Spec-Driven Development Toolkit`, config.Banner, config.Tagline),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Global setup can be added here
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(config.Banner)
			fmt.Println()
			fmt.Println(config.Tagline)
			fmt.Println()
			fmt.Println("Run 'gospecify --help' for usage information")
		},
	}

	// Add subcommands
	cmd.AddCommand(NewInitCmd())
	cmd.AddCommand(NewCheckCmd())
	cmd.AddCommand(NewVersionCmd())

	return cmd
}
