// Package config provides configuration structures and constants for gospecify
package config

// AIAssistant represents an AI coding assistant configuration
type AIAssistant struct {
	Key        string     `json:"key"`
	Name       string     `json:"name"`
	Directory  string     `json:"directory"`
	Format     FileFormat `json:"format"`
	CLITool    string     `json:"cli_tool,omitempty"`
	ArgFormat  string     `json:"arg_format"`
	IsIDEBased bool       `json:"is_ide_based"`
	Website    string     `json:"website"`
}

// FileFormat represents the file format used by an AI assistant
type FileFormat string

const (
	FormatMarkdown FileFormat = "md"
	FormatTOML     FileFormat = "toml"
	FormatPrompt   FileFormat = "prompt.md"
)

// AIAssistants contains all supported AI assistants with their configurations
var AIAssistants = map[string]AIAssistant{
	"copilot": {
		Key:        "copilot",
		Name:       "GitHub Copilot",
		Directory:  ".github/prompts/",
		Format:     FormatPrompt,
		ArgFormat:  "$ARGUMENTS",
		IsIDEBased: true,
		Website:    "https://github.com/features/copilot",
	},
	"claude": {
		Key:       "claude",
		Name:      "Claude Code",
		Directory: ".claude/commands/",
		Format:    FormatMarkdown,
		CLITool:   "claude",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://docs.anthropic.com/en/docs/claude-code/setup",
	},
	"gemini": {
		Key:       "gemini",
		Name:      "Gemini CLI",
		Directory: ".gemini/commands/",
		Format:    FormatTOML,
		CLITool:   "gemini",
		ArgFormat: "{{args}}",
		Website:   "https://github.com/google-gemini/gemini-cli",
	},
	"cursor": {
		Key:       "cursor",
		Name:      "Cursor",
		Directory: ".cursor/commands/",
		Format:    FormatMarkdown,
		CLITool:   "cursor-agent",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://cursor.sh/",
	},
	"qwen": {
		Key:       "qwen",
		Name:      "Qwen Code",
		Directory: ".qwen/commands/",
		Format:    FormatTOML,
		CLITool:   "qwen",
		ArgFormat: "{{args}}",
		Website:   "https://github.com/QwenLM/Qwen2.5-Coder",
	},
	"opencode": {
		Key:       "opencode",
		Name:      "opencode",
		Directory: ".opencode/command/",
		Format:    FormatMarkdown,
		CLITool:   "opencode",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://opencode.ai",
	},
	"codex": {
		Key:       "codex",
		Name:      "Codex CLI",
		Directory: ".codex/",
		Format:    FormatMarkdown,
		CLITool:   "codex",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://github.com/microsoft/codex-cli",
	},
	"windsurf": {
		Key:        "windsurf",
		Name:       "Windsurf",
		Directory:  ".windsurf/workflows/",
		Format:     FormatMarkdown,
		ArgFormat:  "$ARGUMENTS",
		IsIDEBased: true,
		Website:    "https://codeium.com/windsurf",
	},
	"kilocode": {
		Key:       "kilocode",
		Name:      "Kilo Code",
		Directory: ".kilocode/",
		Format:    FormatMarkdown,
		CLITool:   "kilocode",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://kilocode.com/",
	},
	"auggie": {
		Key:       "auggie",
		Name:      "Auggie CLI",
		Directory: ".augment/",
		Format:    FormatMarkdown,
		CLITool:   "auggie",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://augmentcode.com/",
	},
	"roo": {
		Key:       "roo",
		Name:      "Roo Code",
		Directory: ".roo/",
		Format:    FormatMarkdown,
		CLITool:   "roo",
		ArgFormat: "$ARGUMENTS",
		Website:   "https://roocode.com/",
	},
}

// AIChoices provides a mapping of assistant keys to display names for CLI selection
var AIChoices = map[string]string{
	"copilot":  "GitHub Copilot",
	"claude":   "Claude Code",
	"gemini":   "Gemini CLI",
	"cursor":   "Cursor",
	"qwen":     "Qwen Code",
	"opencode": "opencode",
	"codex":    "Codex CLI",
	"windsurf": "Windsurf",
	"kilocode": "Kilo Code",
	"auggie":   "Auggie CLI",
	"roo":      "Roo Code",
}

// AgentFolderMap provides security folder mappings for agents
var AgentFolderMap = map[string]string{
	"claude":   ".claude/",
	"gemini":   ".gemini/",
	"cursor":   ".cursor/",
	"qwen":     ".qwen/",
	"opencode": ".opencode/",
	"codex":    ".codex/",
	"windsurf": ".windsurf/",
	"kilocode": ".kilocode/",
	"auggie":   ".augment/",
	"copilot":  ".github/",
	"roo":      ".roo/",
}
