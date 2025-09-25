# Security Documentation for gospecify

This document outlines the security features and setup requirements for the gospecify release pipeline.

## Overview

gospecify uses GoReleaser with enterprise-grade security features to ensure the integrity and authenticity of all releases:

- **Digital Signatures**: Keyless signing with cosign
- **Software Bill of Materials (SBOM)**: SPDX format dependency tracking
- **Cryptographic Checksums**: SHA256 integrity verification
- **Build Provenance**: Attestations proving build authenticity

## Security Features

### 1. Digital Signatures

All release artifacts are signed using [cosign](https://github.com/sigstore/cosign) with keyless signing. This provides:

- **Non-repudiation**: Proof that artifacts were created by the official build process
- **Integrity**: Detection of any tampering with released files
- **Transparency**: Public verification without secret keys

### 2. Software Bill of Materials (SBOM)

Each release includes an SPDX-formatted SBOM that catalogs:

- All Go module dependencies
- Version information for each component
- License information
- Build environment details

### 3. Cryptographic Checksums

SHA256 checksums are provided for all artifacts to enable:

- Download integrity verification
- Automated verification in CI/CD pipelines
- Manual verification by end users

### 4. Build Provenance

GitHub Attestations provide:

- Proof of build authenticity
- Build environment details
- Supply chain transparency

## Setup Requirements

### GitHub Secrets

The following secrets must be configured in your GitHub repository:

```bash
# Required for GitHub API access
GITHUB_TOKEN=your_github_token_with_repo_and_packages_permissions

# Required for cosign keyless signing
COSIGN_PASSWORD=your_cosign_key_password

# Optional: For traditional GPG signing (alternative to keyless)
COSIGN_PRIVATE_KEY=your_cosign_private_key
COSIGN_PUBLIC_KEY=your_cosign_public_key
```

### Required Permissions

The GitHub Actions workflow requires these permissions:

```yaml
permissions:
  contents: write        # For creating releases
  packages: write        # For publishing packages
  id-token: write        # For keyless signing
  attestations: write    # For build provenance
```

### Tool Installation

The release pipeline requires these tools:

- **GoReleaser**: `go install github.com/goreleaser/goreleaser/v2@latest`
- **cosign**: `go install github.com/sigstore/cosign/v2@latest`
- **Syft**: `curl -sSfL https://raw.githubusercontent.com/anchore/syft/main/install.sh | sh -s -- -b /usr/local/bin`

## Release Process

### Automated Releases

1. **Tag Creation**: Create a tag with the format `gospecify/vX.Y.Z`
   ```bash
   git tag gospecify/v1.0.0
   git push origin gospecify/v1.0.0
   ```

2. **Automated Build**: GitHub Actions triggers the release workflow

3. **Security Processing**:
   - Builds cross-platform binaries
   - Generates SBOM with Syft
   - Signs all artifacts with cosign
   - Creates checksums
   - Generates build attestations

4. **Release Publishing**: Creates GitHub release with all artifacts

### Manual Testing

Test the release process without publishing:

```bash
# Install GoReleaser
make goreleaser-install

# Test build locally
make goreleaser-build

# Test full pipeline (dry-run)
make release-dry-run

# Check configuration
make goreleaser-check
```

## Verification Instructions

### For End Users

1. **Download Checksums**:
   ```bash
   wget https://github.com/github/spec-kit/releases/latest/download/checksums.txt
   ```

2. **Verify Downloads**:
   ```bash
   sha256sum -c checksums.txt
   ```

3. **Verify Signatures**:
   ```bash
   # Download public key
   wget https://github.com/github/spec-kit/releases/latest/download/cosign.pub

   # Verify signature
   cosign verify-blob --key cosign.pub --signature gospecify.tar.gz.sig gospecify.tar.gz
   ```

4. **Check SBOM**:
   ```bash
   # Download SBOM
   wget https://github.com/github/spec-kit/releases/latest/download/gospecify.sbom.json

   # View dependencies
   jq '.packages[] | select(.SPDXID | contains("Package")) | .name' gospecify.sbom.json

   # Scan for vulnerabilities (requires Grype)
   grype sbom:gospecify.sbom.json
   ```

5. **Verify Build Provenance**:
   ```bash
   # Check attestation
   gh attestation verify gospecify -R github/spec-kit
   ```

### For CI/CD Pipelines

```bash
#!/bin/bash
set -e

# Download and verify
wget https://github.com/github/spec-kit/releases/latest/download/checksums.txt
wget https://github.com/github/spec-kit/releases/latest/download/gospecify-linux-amd64.tar.gz

# Verify checksum
sha256sum -c checksums.txt

# Verify signature
cosign verify-blob --key cosign.pub --signature gospecify-linux-amd64.tar.gz.sig gospecify-linux-amd64.tar.gz

echo "✅ All verifications passed"
```

## Security Considerations

### Key Management

- **Keyless Signing**: Recommended approach using OIDC identity
- **No Private Keys**: Eliminates key management and rotation concerns
- **Transparency Log**: All signatures are recorded in the public transparency log

### Supply Chain Security

- **Reproducible Builds**: Same inputs produce identical outputs
- **Dependency Scanning**: SBOM enables automated vulnerability detection
- **Build Environment**: Controlled GitHub Actions environment

### Threat Model

This security model protects against:

- **Artifact Tampering**: Digital signatures prevent modification
- **Supply Chain Attacks**: SBOM enables dependency analysis
- **Build Compromise**: Provenance attestations verify build authenticity
- **Download Corruption**: Checksums ensure file integrity

## Troubleshooting

### Common Issues

1. **Signing Fails**:
   - Ensure `COSIGN_PASSWORD` is set
   - Check GitHub Actions permissions
   - Verify OIDC token availability

2. **SBOM Generation Fails**:
   - Ensure Syft is installed
   - Check Go module integrity
   - Verify network connectivity

3. **Release Upload Fails**:
   - Check `GITHUB_TOKEN` permissions
   - Ensure repository allows releases
   - Verify tag format (`gospecify/vX.Y.Z`)

### Debug Mode

Enable verbose logging:

```bash
# Local testing
goreleaser release --snapshot --clean --verbose

# Check configuration
goreleaser check --verbose
```

## Compliance

This release pipeline supports:

- **NIST SP 800-161**: Supply chain risk management
- **SLSA Level 2**: Build provenance and integrity
- **EU Cybersecurity Act**: Software transparency requirements
- **OpenSSF Best Practices**: Security-focused development

## Contributing

When contributing to the security features:

1. Test changes with `make snapshot`
2. Verify signatures work correctly
3. Ensure SBOM generation succeeds
4. Update this documentation as needed

## References

- [GoReleaser Security Documentation](https://goreleaser.com/customization/sign/)
- [cosign Documentation](https://docs.sigstore.dev/cosign/overview/)
- [Syft SBOM Documentation](https://github.com/anchore/syft)
- [GitHub Attestations](https://docs.github.com/en/actions/security-guides/using-artifact-attestations-to-establish-provenance-for-builds)