// Package github provides GitHub API integration
package github

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/github/spec-kit/gospecify/internal/config"
	"github.com/github/spec-kit/gospecify/pkg/errors"
)

// Release represents a GitHub release
type Release struct {
	TagName string         `json:"tag_name"`
	Assets  []ReleaseAsset `json:"assets"`
}

// ReleaseAsset represents a GitHub release asset
type ReleaseAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// Client handles GitHub API interactions
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// NewClient creates a new GitHub API client
func NewClient(token string, skipTLS bool) *Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if skipTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &Client{
		httpClient: client,
		token:      token,
		baseURL:    config.GitHubAPI,
	}
}

// GetLatestRelease gets the latest release for the spec-kit repository
func (c *Client) GetLatestRelease(ctx context.Context) (*Release, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest",
		c.baseURL, config.GitHubOwner, config.GitHubRepo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkError, "failed to create request", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(errors.ErrCodeNetworkError, "failed to get latest release", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewGitHubAPIError(
			fmt.Sprintf("GitHub API returned %d", resp.StatusCode), nil)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, errors.Wrap(errors.ErrCodeGitHubAPIError, "failed to decode release", err)
	}

	return &release, nil
}

// DownloadAsset downloads a release asset to the specified path
func (c *Client) DownloadAsset(ctx context.Context, asset ReleaseAsset, destPath string, progressFn func(int64, int64)) error {
	req, err := http.NewRequestWithContext(ctx, "GET", asset.BrowserDownloadURL, nil)
	if err != nil {
		return errors.Wrap(errors.ErrCodeNetworkError, "failed to create download request", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	req.Header.Set("User-Agent", config.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(errors.ErrCodeNetworkError, "failed to download asset", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewGitHubAPIError(
			fmt.Sprintf("download failed with status %d", resp.StatusCode), nil)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return errors.Wrap(errors.ErrCodeFileSystemError, "failed to create destination file", err)
	}
	defer file.Close()

	var written int64
	buffer := make([]byte, 32*1024) // 32KB buffer

	for {
		n, readErr := resp.Body.Read(buffer)
		if n > 0 {
			written += int64(n)
			if _, writeErr := file.Write(buffer[:n]); writeErr != nil {
				return errors.Wrap(errors.ErrCodeFileSystemError, "failed to write to file", writeErr)
			}

			if progressFn != nil {
				progressFn(written, resp.ContentLength)
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return errors.Wrap(errors.ErrCodeNetworkError, "failed to read response", readErr)
		}
	}

	return nil
}

// GetGitHubToken retrieves the GitHub token from environment or CLI
func GetGitHubToken(cliToken string) string {
	if cliToken != "" {
		return cliToken
	}

	// Check environment variables
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	return ""
}
