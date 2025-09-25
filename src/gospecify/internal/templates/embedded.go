// Package templates provides embedded template assets
package templates

import (
	"io/fs"
	"path/filepath"
	"strings"

	gospecify "github.com/github/spec-kit/gospecify"
)

// EmbeddedAssets holds all embedded template and script assets
type EmbeddedAssets struct {
	Templates map[string][]byte
	Scripts   map[string][]byte
}

// LoadEmbeddedAssets loads all embedded assets into memory
func LoadEmbeddedAssets() (*EmbeddedAssets, error) {
	assets := &EmbeddedAssets{
		Templates: make(map[string][]byte),
		Scripts:   make(map[string][]byte),
	}

	assetsFS := gospecify.GetAssetsFS()

	// Load all assets
	err := fs.WalkDir(assetsFS, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		content, err := assetsFS.ReadFile(path)
		if err != nil {
			return err
		}

		// Determine if this is a template or script based on path
		if strings.HasPrefix(path, "assets/templates/") {
			relativePath := strings.TrimPrefix(path, "assets/templates/")
			if relativePath != "" {
				assets.Templates[filepath.ToSlash(relativePath)] = content
			}
		} else if strings.HasPrefix(path, "assets/scripts/") {
			relativePath := strings.TrimPrefix(path, "assets/scripts/")
			if relativePath != "" {
				assets.Scripts[filepath.ToSlash(relativePath)] = content
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return assets, nil
}

// GetTemplate retrieves a template by name
func (ea *EmbeddedAssets) GetTemplate(name string) ([]byte, bool) {
	content, exists := ea.Templates[name]
	return content, exists
}

// GetScript retrieves a script by name
func (ea *EmbeddedAssets) GetScript(name string) ([]byte, bool) {
	content, exists := ea.Scripts[name]
	return content, exists
}

// ListTemplates returns a list of all available template names
func (ea *EmbeddedAssets) ListTemplates() []string {
	var names []string
	for name := range ea.Templates {
		names = append(names, name)
	}
	return names
}

// ListScripts returns a list of all available script names
func (ea *EmbeddedAssets) ListScripts() []string {
	var names []string
	for name := range ea.Scripts {
		names = append(names, name)
	}
	return names
}
