// Package scripts provides cross-platform script execution
package scripts

import (
	"fmt"
	"strings"

	"github.com/github/spec-kit/gospecify/internal/config"
	"github.com/github/spec-kit/gospecify/internal/templates"
	"github.com/github/spec-kit/gospecify/pkg/errors"
)

// Generator creates dynamic scripts from embedded templates
type Generator struct {
	assets     *templates.EmbeddedAssets
	assistant  *config.AIAssistant
	scriptType string
}

// NewGenerator creates a new script generator
func NewGenerator(assets *templates.EmbeddedAssets, assistant *config.AIAssistant, scriptType string) *Generator {
	return &Generator{
		assets:     assets,
		assistant:  assistant,
		scriptType: scriptType,
	}
}

// GenerateScript generates a script from an embedded template
func (g *Generator) GenerateScript(scriptName string) ([]byte, error) {
	script, exists := g.assets.GetScript(g.getScriptPath(scriptName))
	if !exists {
		return nil, errors.NewAssetNotFound(fmt.Sprintf("script template %s", scriptName))
	}

	content := string(script)

	// Apply replacements
	content = g.applyReplacements(content)

	return []byte(content), nil
}

// getScriptPath returns the embedded path for a script
func (g *Generator) getScriptPath(scriptName string) string {
	var extension string
	switch g.scriptType {
	case config.ScriptTypeBash:
		extension = ".sh"
	case config.ScriptTypePowerShell:
		extension = ".ps1"
	default:
		extension = ".sh"
	}

	return fmt.Sprintf("%s/%s%s", g.scriptType, scriptName, extension)
}

// applyReplacements applies placeholder replacements to script content
func (g *Generator) applyReplacements(content string) string {
	replacements := map[string]string{
		"__AGENT__": g.assistant.Key,
		"{SCRIPT}":  g.getScriptReference(),
	}

	for placeholder, replacement := range replacements {
		content = strings.ReplaceAll(content, placeholder, replacement)
	}

	return content
}

// getScriptReference returns the appropriate script reference for the assistant
func (g *Generator) getScriptReference() string {
	switch g.scriptType {
	case config.ScriptTypeBash:
		return "./setup.sh"
	case config.ScriptTypePowerShell:
		return ".\\setup.ps1"
	default:
		return "./setup.sh"
	}
}

// GenerateAllScripts generates all scripts for the current configuration
func (g *Generator) GenerateAllScripts() (map[string][]byte, error) {
	scripts := make(map[string][]byte)

	scriptNames := []string{
		"check-prerequisites",
		"create-new-feature",
		"setup-plan",
		"update-agent-context",
	}

	for _, scriptName := range scriptNames {
		content, err := g.GenerateScript(scriptName)
		if err != nil {
			return nil, errors.Wrap(errors.ErrCodeScriptError,
				fmt.Sprintf("failed to generate script %s", scriptName), err)
		}
		scripts[scriptName] = content
	}

	return scripts, nil
}

// ValidateScriptType validates the script type is supported
func (g *Generator) ValidateScriptType() error {
	return ValidateScriptType(g.scriptType)
}
