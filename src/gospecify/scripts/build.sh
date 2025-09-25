#!/bin/bash
# Build script for GoSpecify

set -uo pipefail

# Configuration
VERSION=${VERSION:-"dev"}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Build targets (focus on commonly used platforms)
PLATFORMS=(
    "linux/amd64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Main build function
build_binary() {
    local platform=$1
    local goos=${platform%/*}
    local goarch=${platform#*/}

    local output_name="gospecify-${goos}-${goarch}"
    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi

    log_info "Building for ${goos}/${goarch}..."

    # Set environment variables for cross-compilation
    export GOOS=$goos
    export GOARCH=$goarch
    export CGO_ENABLED=0  # Disable CGO for static binaries

    # Build the binary
    if go build -ldflags "$LDFLAGS" -o "dist/${output_name}" ./cmd/gospecify 2>/dev/null; then
        log_success "Built ${output_name}"
        return 0
    else
        log_error "Failed to build for ${goos}/${goarch}"
        return 1
    fi
}

# Clean previous builds
clean() {
    log_info "Cleaning previous builds..."
    rm -rf dist/
}

# Create distribution directory
setup_dist() {
    mkdir -p dist/
}

# Build for all platforms
build_all() {
    log_info "Starting build for all platforms..."
    log_info "Version: ${VERSION}"
    log_info "Commit: ${COMMIT}"
    log_info "Date: ${DATE}"

    local success_count=0
    local total_count=${#PLATFORMS[@]}

    for platform in "${PLATFORMS[@]}"; do
        if build_binary "$platform"; then
            ((success_count++))
        else
            log_warning "Failed to build for ${platform}"
        fi
    done

    log_info "Build complete: ${success_count}/${total_count} platforms successful"

    if [ $success_count -gt 0 ]; then
        log_success "${success_count} builds completed successfully!"
        list_binaries
        return 0
    else
        log_error "All builds failed"
        return 1
    fi
}

# List built binaries
list_binaries() {
    log_info "Built binaries:"
    ls -lh dist/
}

# Test build
test_build() {
    log_info "Testing build..."

    # Build for current platform
    if go build -ldflags "$LDFLAGS" -o "dist/gospecify-test" ./cmd/gospecify; then
        log_success "Test build successful"

        # Test basic functionality
        if ./dist/gospecify-test version >/dev/null 2>&1; then
            log_success "Binary functional test passed"
        else
            log_error "Binary functional test failed"
            return 1
        fi

        # Clean up test binary
        rm -f dist/gospecify-test
    else
        log_error "Test build failed"
        return 1
    fi
}

# Show usage
usage() {
    cat << EOF
GoSpecify Build Script

Usage: $0 [OPTIONS] [COMMAND]

Commands:
    all         Build for all platforms (default)
    clean       Clean build artifacts
    test        Test build for current platform
    help        Show this help

Options:
    VERSION=X   Set version for build (default: dev)

Examples:
    $0                          # Build for all platforms
    $0 test                     # Test build
    $0 clean                    # Clean builds
    VERSION=v1.0.0 $0           # Build with specific version

EOF
}

# Main script logic
main() {
    local command=${1:-all}

    case $command in
        all)
            clean
            setup_dist
            test_build
            build_all
            ;;
        clean)
            clean
            ;;
        test)
            setup_dist
            test_build
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            log_error "Unknown command: $command"
            usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"