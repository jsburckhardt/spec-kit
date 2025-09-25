# GoSpecify Makefile

.PHONY: all build test clean install help goreleaser release snapshot

# Default target
all: test build

# Build for all platforms (legacy)
build:
	./scripts/build.sh all

# Test build (legacy)
test:
	./scripts/build.sh test

# Clean build artifacts
clean:
	./scripts/build.sh clean

# Install locally
install: test
	go install ./cmd/gospecify

# Run tests
test-go:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Build and show help
help:
	go run ./cmd/gospecify --help

# Development build
dev:
	go build -o gospecify-dev ./cmd/gospecify

# Show version
version:
	go run ./cmd/gospecify version

# Check tools
check:
	go run ./cmd/gospecify check

# Quick test init (dry run)
init-test:
	go run ./cmd/gospecify init --help

# GoReleaser targets

# Install GoReleaser
goreleaser-install:
	go install github.com/goreleaser/goreleaser/v2@latest

# Run GoReleaser in snapshot mode (for testing)
snapshot: goreleaser-install
	goreleaser release --snapshot --clean

# Run GoReleaser build only
goreleaser-build: goreleaser-install
	goreleaser build --clean

# Run GoReleaser check
goreleaser-check: goreleaser-install
	goreleaser check

# Run GoReleaser healthcheck
goreleaser-healthcheck: goreleaser-install
	goreleaser healthcheck

# Full release (requires GITHUB_TOKEN and proper git tag)
release: goreleaser-install
	goreleaser release --clean

# Release dry-run
release-dry-run: goreleaser-install
	goreleaser release --skip=publish --clean