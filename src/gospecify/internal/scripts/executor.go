// Package scripts provides cross-platform script execution
package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/jsburckhardt/spec-kit/gospecify/internal/config"
	"github.com/jsburckhardt/spec-kit/gospecify/pkg/errors"
)

// Executor handles cross-platform script execution
type Executor struct {
	projectPath string
	scriptType  string
}

// NewExecutor creates a new script executor
func NewExecutor(projectPath, scriptType string) *Executor {
	return &Executor{
		projectPath: projectPath,
		scriptType:  scriptType,
	}
}

// ExecuteScript executes a script by name with optional arguments
func (e *Executor) ExecuteScript(scriptName string, args ...string) error {
	script, err := e.getScriptContent(scriptName)
	if err != nil {
		return err
	}

	// Create temporary script file
	tempFile, err := e.createTempScript(script)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(tempFile) }()

	// Execute the script
	return e.executeScriptFile(tempFile, args...)
}

// getScriptContent retrieves script content from embedded assets
func (e *Executor) getScriptContent(scriptName string) ([]byte, error) {
	// This will be implemented when we integrate with the embedded assets
	// For now, return a placeholder error
	return nil, errors.NewAssetNotFound(fmt.Sprintf("script %s", scriptName))
}

// createTempScript creates a temporary script file with proper permissions
func (e *Executor) createTempScript(content []byte) (string, error) {
	var extension string
	switch e.scriptType {
	case config.ScriptTypeBash:
		extension = ".sh"
	case config.ScriptTypePowerShell:
		extension = ".ps1"
	default:
		return "", errors.NewValidationError(fmt.Sprintf("unsupported script type: %s", e.scriptType))
	}

	tempFile, err := os.CreateTemp("", fmt.Sprintf("gospecify-script-*%s", extension))
	if err != nil {
		return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to create temp script", err)
	}
	defer func() { _ = tempFile.Close() }()

	if _, err := tempFile.Write(content); err != nil {
		_ = os.Remove(tempFile.Name())
		return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to write temp script", err)
	}

	// Make executable on Unix systems
	if e.scriptType == config.ScriptTypeBash && runtime.GOOS != "windows" {
		if err := os.Chmod(tempFile.Name(), 0755); err != nil {
			_ = os.Remove(tempFile.Name())
			return "", errors.Wrap(errors.ErrCodeFileSystemError, "failed to make script executable", err)
		}
	}

	return tempFile.Name(), nil
}

// executeScriptFile executes a script file with the appropriate interpreter
func (e *Executor) executeScriptFile(scriptPath string, args ...string) error {
	var cmd *exec.Cmd

	switch e.scriptType {
	case config.ScriptTypeBash:
		cmd = e.createBashCommand(scriptPath, args...)
	case config.ScriptTypePowerShell:
		cmd = e.createPowerShellCommand(scriptPath, args...)
	default:
		return errors.NewValidationError(fmt.Sprintf("unsupported script type: %s", e.scriptType))
	}

	// Set working directory
	cmd.Dir = e.projectPath

	// Set environment
	cmd.Env = os.Environ()

	// Execute command
	if err := cmd.Run(); err != nil {
		return errors.Wrap(errors.ErrCodeScriptError,
			fmt.Sprintf("script execution failed: %s", scriptPath), err)
	}

	return nil
}

// createBashCommand creates a bash command for script execution
func (e *Executor) createBashCommand(scriptPath string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		// On Windows, try to find bash (Git Bash, WSL, etc.)
		if bashPath := e.findBashOnWindows(); bashPath != "" {
			return exec.Command(bashPath, append([]string{scriptPath}, args...)...)
		}
		// Fallback to cmd.exe with bash
		return exec.Command("cmd", append([]string{"/c", "bash", scriptPath}, args...)...)
	}

	// Unix-like systems
	return exec.Command("/bin/bash", append([]string{scriptPath}, args...)...)
}

// createPowerShellCommand creates a PowerShell command for script execution
func (e *Executor) createPowerShellCommand(scriptPath string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		// Windows PowerShell
		return exec.Command("powershell",
			append([]string{"-ExecutionPolicy", "Bypass", "-File", scriptPath}, args...)...)
	}

	// Unix-like systems with PowerShell Core (pwsh)
	return exec.Command("pwsh",
		append([]string{"-File", scriptPath}, args...)...)
}

// findBashOnWindows finds bash executable on Windows systems
func (e *Executor) findBashOnWindows() string {
	possiblePaths := []string{
		"C:\\Program Files\\Git\\bin\\bash.exe",
		"C:\\Program Files (x86)\\Git\\bin\\bash.exe",
		"C:\\Windows\\System32\\bash.exe", // WSL
		"C:\\msys64\\usr\\bin\\bash.exe",
		"C:\\cygwin64\\bin\\bash.exe",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Check if bash is in PATH
	if path, err := exec.LookPath("bash"); err == nil {
		return path
	}

	return ""
}

// ValidateScriptType validates that the script type is supported
func ValidateScriptType(scriptType string) error {
	switch scriptType {
	case config.ScriptTypeBash, config.ScriptTypePowerShell:
		return nil
	default:
		return errors.NewValidationError(fmt.Sprintf("unsupported script type: %s", scriptType))
	}
}

// GetScriptExtension returns the file extension for a script type
func GetScriptExtension(scriptType string) string {
	switch scriptType {
	case config.ScriptTypeBash:
		return ".sh"
	case config.ScriptTypePowerShell:
		return ".ps1"
	default:
		return ""
	}
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsUnix returns true if running on Unix-like systems
func IsUnix() bool {
	return runtime.GOOS != "windows"
}
