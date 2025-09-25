// Package config provides configuration structures and constants for gospecify
package config

// Version information
const (
	Version   = "0.1.0"
	Commit    = "dev"
	Date      = "2025-01-01T00:00:00Z"
	UserAgent = "gospecify/" + Version
)

// GitHub repository information
const (
	GitHubOwner = "github"
	GitHubRepo  = "spec-kit"
	GitHubAPI   = "https://api.github.com"
)

// Default paths and directories
const (
	DefaultTemplateDir = "templates"
	DefaultScriptDir   = "scripts"
	DefaultConfigFile  = ".gospecify.yaml"
)

// Script types
const (
	ScriptTypeBash       = "sh"
	ScriptTypePowerShell = "ps"
)

// Script type configurations
var ScriptTypes = map[string]ScriptType{
	"sh": {
		Key:       "sh",
		Name:      "POSIX Shell (bash/zsh)",
		Extension: ".sh",
		Platform:  "unix",
	},
	"ps": {
		Key:       "ps",
		Name:      "PowerShell",
		Extension: ".ps1",
		Platform:  "windows",
	},
}

// ScriptType represents a script execution environment
type ScriptType struct {
	Key       string
	Name      string
	Extension string
	Platform  string
}

// Banner and branding
const Banner = `
███████╗██████╗ ███████╗ ██████╗██╗███████╗██╗   ██╗
██╔════╝██╔══██╗██╔════╝██╔════╝██║██╔════╝╚██╗ ██╔╝
███████╗██████╔╝█████╗  ██║     ██║█████╗   ╚████╔╝
╚════██║██╔═══╝ ██╔══╝  ██║     ██║██╔══╝    ╚██╔╝
███████║██║     ███████╗╚██████╗██║██║        ██║
╚══════╝╚═╝     ╚══════╝ ╚═════╝╚═╝╚═╝        ╚═╝   `

const Tagline = "GitHub Spec Kit - Spec-Driven Development Toolkit"

// Security notice template
const SecurityNoticeTemplate = `Some agents may store credentials, auth tokens, or other identifying and private artifacts in the agent folder within your project.
Consider adding [cyan]%s[/cyan] (or parts of it) to [cyan].gitignore[/cyan] to prevent accidental credential leakage.`

// Next steps template
const NextStepsTemplate = `1. Go to the project folder: [cyan]cd %s[/cyan]
2. Start using slash commands with your AI agent:`
