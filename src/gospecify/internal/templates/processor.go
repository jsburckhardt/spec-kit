// Package templates provides template processing functionality
package templates

import (
	"fmt"
	"strings"

	"github.com/jsburckhardt/spec-kit/gospecify/internal/config"
	"github.com/jsburckhardt/spec-kit/gospecify/pkg/errors"
)

// Processor handles template processing for different AI assistants
type Processor struct {
	assets     *EmbeddedAssets
	assistant  *config.AIAssistant
	scriptType string
}

// NewProcessor creates a new template processor
func NewProcessor(assets *EmbeddedAssets, assistant *config.AIAssistant, scriptType string) *Processor {
	return &Processor{
		assets:     assets,
		assistant:  assistant,
		scriptType: scriptType,
	}
}

// ProcessTemplate processes a template and returns the processed content
func (p *Processor) ProcessTemplate(templateName string) ([]byte, error) {
	template, exists := p.assets.GetTemplate(templateName)
	if !exists {
		return nil, errors.NewAssetNotFound(fmt.Sprintf("template %s", templateName))
	}

	content := string(template)

	// Apply replacements
	content = p.applyReplacements(content)

	// Apply format-specific processing
	switch p.assistant.Format {
	case config.FormatMarkdown:
		return p.processMarkdownTemplate(content)
	case config.FormatTOML:
		return p.processTOMLTemplate(content)
	case config.FormatPrompt:
		return p.processPromptTemplate(content)
	default:
		return []byte(content), nil
	}
}

// applyReplacements applies common placeholder replacements
func (p *Processor) applyReplacements(content string) string {
	replacements := map[string]string{
		"__AGENT__":  p.assistant.Key,
		"$ARGUMENTS": p.assistant.ArgFormat,
		"{{args}}":   p.assistant.ArgFormat,
		"{SCRIPT}":   "", // Will be set per template
	}

	for placeholder, replacement := range replacements {
		content = strings.ReplaceAll(content, placeholder, replacement)
	}

	return content
}

// processMarkdownTemplate processes Markdown format templates
func (p *Processor) processMarkdownTemplate(content string) ([]byte, error) {
	lines := strings.Split(content, "\n")
	var processed []string

	for _, line := range lines {
		// Handle script placeholders in Markdown
		if strings.Contains(line, "{SCRIPT}") {
			line = strings.ReplaceAll(line, "{SCRIPT}", p.getScriptPath())
		}
		processed = append(processed, line)
	}

	return []byte(strings.Join(processed, "\n")), nil
}

// processTOMLTemplate processes TOML format templates
func (p *Processor) processTOMLTemplate(content string) ([]byte, error) {
	lines := strings.Split(content, "\n")
	var processed []string

	inPrompt := false
	for _, line := range lines {
		// Track prompt sections
		if strings.Contains(line, `prompt = """`) {
			inPrompt = true
		} else if inPrompt && strings.Contains(line, `"""`) {
			inPrompt = false
		}

		// Apply replacements in prompt sections
		if inPrompt {
			line = strings.ReplaceAll(line, "$ARGUMENTS", p.assistant.ArgFormat)
			line = strings.ReplaceAll(line, "{{args}}", p.assistant.ArgFormat)
		}

		processed = append(processed, line)
	}

	return []byte(strings.Join(processed, "\n")), nil
}

// processPromptTemplate processes prompt.md format templates
func (p *Processor) processPromptTemplate(content string) ([]byte, error) {
	// Prompt templates are simpler, just apply basic replacements
	content = strings.ReplaceAll(content, "{SCRIPT}", p.getScriptPath())
	return []byte(content), nil
}

// getScriptPath returns the appropriate script path for the current configuration
func (p *Processor) getScriptPath() string {
	return fmt.Sprintf(".specify/scripts/%s", p.getScriptName())
}

// getScriptName returns the script name for the current configuration
func (p *Processor) getScriptName() string {
	switch p.scriptType {
	case config.ScriptTypeBash:
		return "setup.sh"
	case config.ScriptTypePowerShell:
		return "setup.ps1"
	default:
		return "setup.sh"
	}
}

// ProcessAllTemplates processes all templates for the current assistant
func (p *Processor) ProcessAllTemplates() (map[string][]byte, error) {
	processed := make(map[string][]byte)

	templates := p.assets.ListTemplates()
	for _, templateName := range templates {
		content, err := p.ProcessTemplate(templateName)
		if err != nil {
			return nil, errors.Wrap(errors.ErrCodeTemplateError,
				fmt.Sprintf("failed to process template %s", templateName), err)
		}
		processed[templateName] = content
	}

	return processed, nil
}
