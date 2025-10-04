# Modifications Summary: Release Fix

## Date
October 4, 2025

## Branch
`fix/release` → `main`

## Changes Made

### 1. GoReleaser Workflow Configuration
**File**: `.github/workflows/goreleaser.yml`

**Changes**:
- Simplified tag pattern from `gospecify/v*` to `v*` across all workflow conditions:
  - Workflow trigger conditions
  - Production release job conditions  
  - Security scan job conditions
  - Provenance job conditions

**Rationale**: Standard `v*` tag pattern is more conventional and simplifies the release process.

### 2. Gitignore Update  
**File**: `.gitignore`

**Changes**:
- Added `flowspace` directory to ignore list
- Fixed file formatting (added missing newline at end)

**Rationale**: Prevents tracking of `flowspace` directory and improves file formatting.

## Pull Request
- **URL**: https://github.com/jsburckhardt/spec-kit/pull/5
- **Title**: "fix: simplify GoReleaser tag patterns and update gitignore"
- **Type**: Bug Fix
- **Files Modified**: 2

## Impact
- Simplified and standardized release process
- Consistent workflow behavior across all jobs
- Cleaner repository structure

## Commit Message Format
Following conventional commit format: `fix: update GoReleaser tag patterns and add flowspace to .gitignore`
