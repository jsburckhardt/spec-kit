// Package github provides GitHub API integration
package github

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jsburckhardt/spec-kit/gospecify/pkg/errors"
)

// Extractor handles template extraction from zip archives
type Extractor struct {
	destDir string
}

// NewExtractor creates a new template extractor
func NewExtractor(destDir string) *Extractor {
	return &Extractor{
		destDir: destDir,
	}
}

// ExtractZip extracts a zip archive to the destination directory
func (e *Extractor) ExtractZip(zipPath string, progressFn func(int64, int64)) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to open zip file", err)
	}
	defer func() { _ = reader.Close() }()

	// Calculate total size for progress
	var totalSize int64
	for _, file := range reader.File {
		totalSize += int64(file.UncompressedSize64)
	}

	var extractedSize int64

	// Extract files
	for _, file := range reader.File {
		if err := e.extractFile(file, &extractedSize, totalSize, progressFn); err != nil {
			return err
		}
	}

	return nil
}

// extractFile extracts a single file from the zip archive
func (e *Extractor) extractFile(file *zip.File, extractedSize *int64, totalSize int64, progressFn func(int64, int64)) error {
	// Open the file in the zip
	src, err := file.Open()
	if err != nil {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to open file in zip", err)
	}
	defer func() { _ = src.Close() }()

	// Construct destination path
	destPath := filepath.Join(e.destDir, file.Name)

	// Create directory if needed
	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(destPath, file.Mode()); err != nil {
			return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create directory", err)
		}
		return nil
	}

	// Create parent directory
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create parent directory", err)
	}

	// Create destination file
	dest, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create destination file", err)
	}
	defer func() { _ = dest.Close() }()

	// Copy file contents
	written, err := io.Copy(dest, src)
	if err != nil {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to copy file contents", err)
	}

	*extractedSize += written
	if progressFn != nil {
		progressFn(*extractedSize, totalSize)
	}

	return nil
}

// FindTemplateAsset finds the appropriate template asset for the given AI assistant
func FindTemplateAsset(release *Release, aiAssistant string) (*ReleaseAsset, error) {
	// Look for asset matching the pattern: spec-kit-template-{aiAssistant}-{scriptType}-{version}.zip
	pattern := fmt.Sprintf("spec-kit-template-%s-", aiAssistant)

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, pattern) && strings.HasSuffix(asset.Name, ".zip") {
			return &asset, nil
		}
	}

	return nil, errors.NewAssetNotFound(fmt.Sprintf("template asset for %s", aiAssistant))
}

// Cleanup removes temporary files
func (e *Extractor) Cleanup(tempFile string) error {
	if err := os.Remove(tempFile); err != nil && !os.IsNotExist(err) {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to cleanup temporary file", err)
	}
	return nil
}
