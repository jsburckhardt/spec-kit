// Package errors provides custom error types for gospecify
package errors

import "fmt"

// Error represents a gospecify error with additional context
type Error struct {
	Code    string
	Message string
	Cause   error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// Error codes
const (
	ErrCodeInvalidConfig   = "INVALID_CONFIG"
	ErrCodeNetworkError    = "NETWORK_ERROR"
	ErrCodeGitError        = "GIT_ERROR"
	ErrCodeTemplateError   = "TEMPLATE_ERROR"
	ErrCodeScriptError     = "SCRIPT_ERROR"
	ErrCodeValidationError = "VALIDATION_ERROR"
	ErrCodeFileSystemError = "FILESYSTEM_ERROR"
	ErrCodeGitHubAPIError  = "GITHUB_API_ERROR"
	ErrCodeAssetNotFound   = "ASSET_NOT_FOUND"
	ErrCodeToolNotFound    = "TOOL_NOT_FOUND"
)

// New creates a new Error with the given code and message
func New(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewInvalidConfig creates an invalid configuration error
func NewInvalidConfig(message string) *Error {
	return New(ErrCodeInvalidConfig, message)
}

// NewNetworkError creates a network error
func NewNetworkError(message string, cause error) *Error {
	return Wrap(ErrCodeNetworkError, message, cause)
}

// NewGitError creates a git operation error
func NewGitError(message string, cause error) *Error {
	return Wrap(ErrCodeGitError, message, cause)
}

// NewTemplateError creates a template processing error
func NewTemplateError(message string, cause error) *Error {
	return Wrap(ErrCodeTemplateError, message, cause)
}

// NewScriptError creates a script execution error
func NewScriptError(message string, cause error) *Error {
	return Wrap(ErrCodeScriptError, message, cause)
}

// NewValidationError creates a validation error
func NewValidationError(message string) *Error {
	return New(ErrCodeValidationError, message)
}

// NewFileSystemError creates a filesystem error
func NewFileSystemError(message string, cause error) *Error {
	return Wrap(ErrCodeFileSystemError, message, cause)
}

// NewGitHubAPIError creates a GitHub API error
func NewGitHubAPIError(message string, cause error) *Error {
	return Wrap(ErrCodeGitHubAPIError, message, cause)
}

// NewAssetNotFound creates an asset not found error
func NewAssetNotFound(assetName string) *Error {
	return New(ErrCodeAssetNotFound, fmt.Sprintf("asset not found: %s", assetName))
}

// NewToolNotFound creates a tool not found error
func NewToolNotFound(toolName string) *Error {
	return New(ErrCodeToolNotFound, fmt.Sprintf("required tool not found: %s", toolName))
}
