// Package gospecify provides embedded assets
package gospecify

import "embed"

//go:embed assets/*
var AssetsFS embed.FS

// GetAssetsFS returns the embedded assets filesystem
func GetAssetsFS() embed.FS {
	return AssetsFS
}
