// Package cmd provides the CLI commands for gospecify
package cmd

import (
	"fmt"

	"github.com/github/spec-kit/gospecify/internal/config"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Show version information for gospecify",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gospecify %s\n", config.Version)
			fmt.Printf("Commit: %s\n", config.Commit)
			fmt.Printf("Built: %s\n", config.Date)
		},
	}
}
