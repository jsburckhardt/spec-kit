// Package cmd provides the CLI commands for gospecify
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/jsburckhardt/spec-kit/gospecify/internal/config"
	"github.com/jsburckhardt/spec-kit/gospecify/internal/ui"
	"github.com/spf13/cobra"
)

// NewCheckCmd creates the check command
func NewCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Check that required tools are installed",
		Long: `Check that required tools are installed for gospecify to work properly.

This command verifies that all necessary tools are available:
- Git (optional, but recommended)
- AI assistant CLI tools (optional, checked per assistant)

Examples:
  gospecify check`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCheck()
		},
	}
}

// runCheck executes the check command
func runCheck() error {
	fmt.Println(ui.InfoPanel.Render("🔍 Checking system for required tools..."))
	fmt.Println()

	// Check results
	results := make(map[string]bool)

	// Check git
	results["git"] = checkTool("git", "Version control system")

	// Check AI assistant tools
	assistantTools := []string{"claude", "gemini", "cursor", "qwen", "opencode", "codex", "kilocode", "auggie", "roo"}
	for _, tool := range assistantTools {
		if assistant, exists := config.AIAssistants[tool]; exists && assistant.CLITool != "" {
			results[assistant.CLITool] = checkTool(assistant.CLITool, fmt.Sprintf("CLI for %s", assistant.Name))
		}
	}

	// Display results
	fmt.Println("📋 Tool Check Results:")
	fmt.Println()

	allGood := true
	for tool, available := range results {
		if available {
			fmt.Printf("✅ %s - Available\n", tool)
		} else {
			fmt.Printf("❌ %s - Not found\n", tool)
			allGood = false
		}
	}

	fmt.Println()

	if allGood {
		fmt.Println(ui.SuccessPanel.Render("🎉 All tools are properly installed!"))
	} else {
		fmt.Println(ui.WarningPanel.Render("⚠️  Some tools are missing. You can still use gospecify, but some AI assistants may not work."))
		fmt.Println()
		fmt.Println("💡 To install missing tools:")
		fmt.Println("   - Git: https://git-scm.com/downloads")
		fmt.Println("   - Claude Code: https://docs.anthropic.com/en/docs/claude-code/setup")
		fmt.Println("   - Other AI tools: Check their respective documentation")
	}

	return nil
}

// checkTool checks if a tool is available on the system
func checkTool(tool, description string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}
