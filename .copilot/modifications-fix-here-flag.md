# Fix --here Flag Issue

## Problem
The `gospecify init --here` command was incorrectly throwing the error:
```
Directory already exists. Use --force to overwrite or specify a different name
```

This error occurred because the validation logic was checking if the current directory exists when using `--here`, which it always will (that's the point of `--here`).

## Root Cause
In `src/gospecify/cmd/init.go`, the `validateConfig()` function had incorrect logic:

```go
// OLD (incorrect) logic:
if _, err := os.Stat(cfg.Path); err == nil {
    if cfg.Here && !cfg.Force {
        return errors.NewValidationError(
            "Directory already exists. Use --force to overwrite or specify a different name")
    }
    if !cfg.Here {
        return errors.NewValidationError(
            fmt.Sprintf("Directory %s already exists", cfg.Path))
    }
}
```

The problem was that when using `--here`, `cfg.Path` is set to the current working directory, which always exists. The code would then error out saying "directory already exists" even though that's expected behavior for `--here`.

## Solution
Fixed the validation logic to only check for directory existence when NOT using `--here`:

```go
// NEW (correct) logic:
if !cfg.Here {
    // When not using --here, we're creating a new directory that shouldn't exist
    if _, err := os.Stat(cfg.Path); err == nil {
        return errors.NewValidationError(
            fmt.Sprintf("Directory %s already exists", cfg.Path))
    }
}
// When using --here, the current directory should exist and we don't need to check
```

## Behavior After Fix
- `gospecify init --here` works in empty directories without `--force`
- `gospecify init --here` requires `--force` in non-empty directories (correct behavior)
- `gospecify init project-name` still correctly errors if `project-name` directory already exists
- The `--force` flag now only applies to overwriting existing project files, not directory existence checks

## Files Changed
- `src/gospecify/cmd/init.go` - Fixed validation logic in `validateConfig()` function

## Testing
Tested the following scenarios successfully:
1. ✅ `gospecify init --here` in empty directory (works without --force)
2. ✅ `gospecify init --here` in non-empty directory (requires --force)
3. ✅ `gospecify init --here --force` in non-empty directory (works)
4. ✅ `gospecify init existing-dir` still errors correctly when directory exists
