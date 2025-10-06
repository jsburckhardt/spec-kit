// Package cmd provides the CLI commands for gospecify
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jsburckhardt/spec-kit/gospecify/internal/config"
	"github.com/jsburckhardt/spec-kit/gospecify/internal/scripts"
	"github.com/jsburckhardt/spec-kit/gospecify/internal/templates"
	"github.com/jsburckhardt/spec-kit/gospecify/internal/ui"
	"github.com/jsburckhardt/spec-kit/gospecify/pkg/errors"
	"github.com/spf13/cobra"
)

// NewInitCmd creates the init command
func NewInitCmd() *cobra.Command {
	var cfg config.ProjectConfig

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Specify project from the latest template",
		Long: `Initialize a new Specify project from the latest template.

This command will:
1. Check that required tools are installed (git is optional)
2. Let you choose your AI assistant (Claude Code, Gemini CLI, GitHub Copilot, etc.)
3. Download the appropriate template from GitHub
4. Extract the template to a new project directory or current directory
5. Initialize a fresh git repository (if not --no-git and no existing repo)
6. Optionally set up AI assistant commands

Examples:
  gospecify init my-project
  gospecify init my-project --ai claude
  gospecify init --here --ai claude
  gospecify init --here --force`,
		Args: func(cmd *cobra.Command, args []string) error {
			if cfg.Here && len(args) > 0 {
				return fmt.Errorf("cannot specify both project name and --here flag")
			}
			if !cfg.Here && len(args) == 0 {
				return fmt.Errorf("must specify either a project name or use --here flag")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				cfg.Name = args[0]
			}
			cfg.CreatedAt = time.Now()
			return runInit(&cfg)
		},
	}

	// Flags
	cmd.Flags().StringVar(&cfg.AIAssistant, "ai", "",
		"AI assistant to use: claude, gemini, copilot, cursor, qwen, opencode, windsurf, kilocode, auggie, or roo")
	cmd.Flags().StringVar(&cfg.ScriptType, "script", "",
		"Script type to use: sh or ps")
	cmd.Flags().BoolVar(&cfg.IgnoreTools, "ignore-agent-tools", false,
		"Skip checks for AI agent tools like Claude Code")
	cmd.Flags().BoolVar(&cfg.NoGit, "no-git", false,
		"Skip git repository initialization")
	cmd.Flags().BoolVar(&cfg.Here, "here", false,
		"Initialize project in the current directory instead of creating a new one")
	cmd.Flags().BoolVar(&cfg.Force, "force", false,
		"Force merge/overwrite when using --here (skip confirmation)")
	cmd.Flags().BoolVar(&cfg.SkipTLS, "skip-tls", false,
		"Skip SSL/TLS verification (not recommended)")
	cmd.Flags().BoolVar(&cfg.Debug, "debug", false,
		"Show verbose diagnostic output for network and extraction failures")
	cmd.Flags().StringVar(&cfg.GitHubToken, "github-token", "",
		"GitHub token to use for API requests (or set GH_TOKEN or GITHUB_TOKEN environment variable)")

	return cmd
}

// runInit executes the init command
func runInit(cfg *config.ProjectConfig) error {

	// Initialize progress tracker
	tracker := &config.StepTracker{
		Title: "Initializing Specify Project",
	}
	tracker.Add("validate", "Validate configuration")
	tracker.Add("assistant", "Select AI assistant")
	tracker.Add("script", "Select script type")
	tracker.Add("tools", "Check required tools")
	tracker.Add("download", "Prepare project directory")
	tracker.Add("extract", "Setup embedded assets")
	tracker.Add("process", "Process templates")
	tracker.Add("scripts", "Generate scripts")
	tracker.Add("git", "Initialize git repository")

	// Set up live progress display
	progress := ui.NewLiveProgress(tracker)
	fmt.Println(progress.Render())

	// Step 1: Validate configuration
	tracker.Start("validate", "")
	if err := validateConfig(cfg); err != nil {
		tracker.Error("validate", err.Error())
		return err
	}
	tracker.Complete("validate", "Configuration valid")

	// Step 2: Select AI assistant
	tracker.Start("assistant", "")
	assistant, err := selectAssistant(cfg)
	if err != nil {
		tracker.Error("assistant", err.Error())
		return err
	}
	cfg.AIAssistant = assistant.Key
	tracker.Complete("assistant", fmt.Sprintf("Selected %s", assistant.Name))

	// Step 3: Select script type
	tracker.Start("script", "")
	scriptType, err := selectScriptType(cfg)
	if err != nil {
		tracker.Error("script", err.Error())
		return err
	}
	cfg.ScriptType = scriptType
	tracker.Complete("script", fmt.Sprintf("Selected %s", config.ScriptTypes[scriptType].Name))

	// Step 4: Check required tools
	tracker.Start("tools", "")
	if err := checkRequiredTools(assistant, cfg.IgnoreTools); err != nil {
		tracker.Error("tools", err.Error())
		return err
	}
	tracker.Complete("tools", "All tools available")

	// Step 5: Prepare project directory (using embedded assets)
	tracker.Start("download", "")
	projectPath, err := prepareProjectDirectory(cfg)
	if err != nil {
		tracker.Error("download", err.Error())
		return err
	}
	cfg.Path = projectPath
	tracker.Complete("download", "Project directory prepared")

	// Step 6: Skip extract (using embedded assets only)
	tracker.Start("extract", "")
	tracker.Complete("extract", "Using embedded assets")

	// Step 7: Process templates
	tracker.Start("process", "")
	if err := processTemplates(projectPath, assistant, scriptType); err != nil {
		tracker.Error("process", err.Error())
		return err
	}
	tracker.Complete("process", "Templates processed")

	// Step 8: Generate scripts
	tracker.Start("scripts", "")
	if err := generateScripts(projectPath, assistant, scriptType); err != nil {
		tracker.Error("scripts", err.Error())
		return err
	}
	tracker.Complete("scripts", "Scripts generated")

	// Step 9: Initialize git repository
	tracker.Start("git", "")
	if err := initializeGit(projectPath, cfg.NoGit); err != nil {
		tracker.Error("git", err.Error())
		return err
	}
	if !cfg.NoGit {
		tracker.Complete("git", "Git repository initialized")
	} else {
		tracker.Skip("git", "Skipped")
	}

	// Show success message and next steps
	if err := showSuccessMessage(cfg, assistant); err != nil {
		return err
	}

	return nil
}

// validateConfig validates the initial configuration
func validateConfig(cfg *config.ProjectConfig) error {
	if cfg.Here {
		cwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to get current directory", err)
		}
		cfg.Path = cwd
		cfg.Name = filepath.Base(cwd)
	} else {
		var err error
		cfg.Path, err = filepath.Abs(cfg.Name)
		if err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to resolve project path", err)
		}
	}

	// Check if directory exists - only relevant when creating new project directory
	if !cfg.Here {
		// When not using --here, we're creating a new directory that shouldn't exist
		if _, err := os.Stat(cfg.Path); err == nil {
			return errors.NewValidationError(
				fmt.Sprintf("Directory %s already exists", cfg.Path))
		}
	}
	// When using --here, the current directory should exist and we don't need to check

	return nil
}

// selectAssistant selects the AI assistant to use
func selectAssistant(cfg *config.ProjectConfig) (*config.AIAssistant, error) {
	if cfg.AIAssistant != "" {
		assistant, exists := config.AIAssistants[cfg.AIAssistant]
		if !exists {
			return nil, errors.NewValidationError(
				fmt.Sprintf("Unknown AI assistant: %s", cfg.AIAssistant))
		}
		return &assistant, nil
	}

	// Interactive selection
	selector := ui.NewSelector("Select your AI assistant", config.AIChoices, "claude")
	selected, err := selector.Run()
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeValidationError, "assistant selection failed", err)
	}

	assistant := config.AIAssistants[selected]
	return &assistant, nil
}

// selectScriptType selects the script type to use
func selectScriptType(cfg *config.ProjectConfig) (string, error) {
	if cfg.ScriptType != "" {
		if _, exists := config.ScriptTypes[cfg.ScriptType]; !exists {
			return "", errors.NewValidationError(
				fmt.Sprintf("Unknown script type: %s", cfg.ScriptType))
		}
		return cfg.ScriptType, nil
	}

	// Interactive selection
	scriptChoices := make(map[string]string)
	for key, scriptType := range config.ScriptTypes {
		scriptChoices[key] = scriptType.Name
	}

	selector := ui.NewSelector("Select your script type", scriptChoices, "sh")
	selected, err := selector.Run()
	if err != nil {
		return "", errors.Wrap(errors.ErrCodeValidationError, "script type selection failed", err)
	}

	return selected, nil
}

// checkRequiredTools checks that required tools are available
func checkRequiredTools(assistant *config.AIAssistant, ignoreTools bool) error {
	if ignoreTools {
		return nil
	}

	// Check for git (optional)
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Warning: git not found. Consider installing git for version control.")
	}

	// Check for AI assistant CLI tool (if required)
	if assistant.CLITool != "" {
		if _, err := exec.LookPath(assistant.CLITool); err != nil {
			return errors.NewToolNotFound(assistant.CLITool)
		}
	}

	return nil
}

// processTemplates processes templates from embedded assets and creates project structure
func processTemplates(projectPath string, assistant *config.AIAssistant, scriptType string) error {
	// Load embedded assets
	assets, err := templates.LoadEmbeddedAssets()
	if err != nil {
		return errors.Wrap(errors.ErrCodeAssetNotFound, "failed to load embedded assets", err)
	}

	// Create template processor
	processor := templates.NewProcessor(assets, assistant, scriptType)

	// Create base project structure
	dirs := []string{
		".specify/templates",
		".specify/templates/commands",
		assistant.Directory,
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create directory", err)
		}
	}

	// Process and write all templates to project directory
	processedTemplates, err := processor.ProcessAllTemplates()
	if err != nil {
		return err
	}

	for templateName, content := range processedTemplates {
		templatePath := filepath.Join(projectPath, ".specify", "templates", templateName)
		if err := os.MkdirAll(filepath.Dir(templatePath), 0755); err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create template directory", err)
		}

		if err := os.WriteFile(templatePath, content, 0644); err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to write template", err)
		}
	}

	// Copy command templates to assistant folder
	for templateName, content := range processedTemplates {
		if strings.HasPrefix(templateName, "commands/") {
			commandName := strings.TrimPrefix(templateName, "commands/")
			commandPath := filepath.Join(projectPath, assistant.Directory, commandName)

			if err := os.WriteFile(commandPath, content, 0644); err != nil {
				return errors.Wrap(errors.ErrCodeFileSystemError, "failed to write command template", err)
			}
		}
	}

	return nil
}

// generateScripts generates the setup scripts
func generateScripts(projectPath string, assistant *config.AIAssistant, scriptType string) error {
	// Load embedded assets
	assets, err := templates.LoadEmbeddedAssets()
	if err != nil {
		return errors.Wrap(errors.ErrCodeAssetNotFound, "failed to load embedded assets", err)
	}

	// Create script generator
	generator := scripts.NewGenerator(assets, assistant, scriptType)

	// Generate all scripts
	generatedScripts, err := generator.GenerateAllScripts()
	if err != nil {
		return err
	}

	// Write scripts to project directory
	scriptsDir := filepath.Join(projectPath, ".specify", "scripts")
	for scriptName, content := range generatedScripts {
		scriptPath := filepath.Join(scriptsDir, scriptName+scripts.GetScriptExtension(scriptType))
		if err := os.MkdirAll(scriptsDir, 0755); err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create scripts directory", err)
		}

		if err := os.WriteFile(scriptPath, content, 0755); err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to write script", err)
		}
	}

	return nil
}

// initializeGit initializes a git repository if requested
func initializeGit(projectPath string, noGit bool) error {
	if noGit {
		return nil
	}

	// Check if already a git repository
	if _, err := os.Stat(filepath.Join(projectPath, ".git")); err == nil {
		return nil // Already a git repo
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrCodeGitError, "failed to initialize git repository", err)
	}

	// Create initial commit
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrCodeGitError, "failed to add files to git", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit - Specify project setup")
	cmd.Dir = projectPath
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrCodeGitError, "failed to create initial commit", err)
	}

	return nil
}

// showSuccessMessage displays success message and next steps
func showSuccessMessage(cfg *config.ProjectConfig, assistant *config.AIAssistant) error {
	fmt.Println()
	fmt.Println(ui.InfoPanel.Render(fmt.Sprintf("✅ Successfully initialized Specify project in %s", cfg.Path)))
	fmt.Println()

	// Show security notice
	if folder, exists := config.AgentFolderMap[assistant.Key]; exists {
		securityMessage := fmt.Sprintf(
			"Some agents may store credentials, auth tokens, or other identifying and private artifacts in the agent folder within your project.\nConsider adding %s (or parts of it) to %s to prevent accidental credential leakage.",
			ui.CyanStyle.Render(folder),
			ui.CyanStyle.Render(".gitignore"))
		fmt.Println(ui.WarningPanel.Render(securityMessage))
		fmt.Println()
	}

	// Show next steps
	steps := []string{
		fmt.Sprintf("1. Go to the project folder: %s", ui.CyanStyle.Render(fmt.Sprintf("cd %s", cfg.Name))),
		"2. Start using slash commands with your AI agent:",
		fmt.Sprintf("   - %s - Analyze codebase and requirements", ui.CyanStyle.Render("/analyze")),
		fmt.Sprintf("   - %s - Get clarification on requirements", ui.CyanStyle.Render("/clarify")),
		fmt.Sprintf("   - %s - Implement features from specifications", ui.CyanStyle.Render("/implement")),
		fmt.Sprintf("   - %s - Create development plans", ui.CyanStyle.Render("/plan")),
		fmt.Sprintf("   - %s - Generate detailed specifications", ui.CyanStyle.Render("/specify")),
		fmt.Sprintf("   - %s - Break down work into tasks", ui.CyanStyle.Render("/tasks")),
	}

	fmt.Println(ui.SuccessPanel.Render(strings.Join(steps, "\n")))

	return nil
}

// prepareProjectDirectory creates the project directory structure without GitHub download
func prepareProjectDirectory(cfg *config.ProjectConfig) (string, error) {
	var projectPath string

	if cfg.Here {
		// Use current directory
		cwd, err := os.Getwd()
		if err != nil {
			return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to get current directory", err)
		}
		projectPath = cwd

		// Check if directory is empty or force flag is set
		entries, err := os.ReadDir(projectPath)
		if err != nil {
			return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to read current directory", err)
		}

		if len(entries) > 0 && !cfg.Force {
			return "", errors.New(errors.ErrCodeValidationError, "directory is not empty (use --force to override)")
		}
	} else {
		// Create new project directory
		projectPath = filepath.Join(".", cfg.Name)

		// Check if directory already exists
		if _, err := os.Stat(projectPath); err == nil {
			return "", errors.New(errors.ErrCodeValidationError, fmt.Sprintf("Directory %s already exists", projectPath))
		}

		// Create the project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to create project directory", err)
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to get absolute path", err)
	}

	return absPath, nil
}
