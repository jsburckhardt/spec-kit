# Fixed gospecify Embedded Assets

## Issues Fixed

### 1. Script Path Bug
The gospecify tool was failing during project initialization with:
```
Error: failed to generate script check-prerequisites: asset not found: script template check-prerequisites
```

**Root Cause**: The `getScriptPath` method was using script type constants ("sh"/"ps") as directory names, but embedded assets use "bash"/"powershell" directories.

**Fix**: Updated path construction to use correct directory names.

### 2. Unnecessary GitHub Downloads
The tool was **downloading templates from GitHub** on every `init` command instead of using embedded assets.

**Root Cause**: The init process was:
1. Download GitHub release ❌ (unnecessary)
2. Extract template ❌ (unnecessary)
3. Load embedded assets ✅ (should be only step)
4. Process templates using embedded assets

**Fix**: Completely removed GitHub download/extract and use embedded assets directly.

## Changes Made

### `/src/gospecify/internal/scripts/generator.go`
```go
// Before (incorrect paths)
return fmt.Sprintf("%s/%s%s", g.scriptType, scriptName, extension)

// After (correct directory mapping)
var directory, extension string
switch g.scriptType {
case config.ScriptTypeBash:
    directory = "bash"      // was "sh"
    extension = ".sh"
case config.ScriptTypePowerShell:
    directory = "powershell" // was "ps"
    extension = ".ps1"
}
return fmt.Sprintf("%s/%s%s", directory, scriptName, extension)
```

### `/src/gospecify/cmd/init.go`
- **Removed**: `downloadTemplate()` and `extractTemplate()` functions
- **Added**: `prepareProjectDirectory()` function
- **Modified**: `processTemplates()` to create full project structure from embedded assets
- **Removed**: Unused imports (context, github, progressbar)

## Benefits

1. **Much faster initialization** - No network requests
2. **Works offline** - No GitHub dependency
3. **More reliable** - No network failures
4. **Cleaner code** - Removed 60+ lines of unused download/extract code
5. **Fixed original bug** - Script generation now works

## Testing Results

```bash
# Before: Downloaded from GitHub + script errors + ugly output
❌ gospecify init test --ai claude
Error: failed to generate script check-prerequisites: asset not found...
Output: Consider adding [cyan].claude/[/cyan] to [cyan].gitignore[/cyan]

# After: Pure embedded assets, fast, working, and clean output
✅ gospecify init test --ai claude --ignore-agent-tools --no-git
Successfully initialized Specify project (no downloads!)
Clean output: Consider adding .claude/ (or parts of it) to .gitignore
```

Project structure created correctly with all templates and scripts.

## Additional Fix: Terminal Output Formatting

**Issue**: Raw `[cyan]` markup was showing in terminal instead of actual colors
**Fix**: Replaced manual `[cyan]` strings with proper `ui.CyanStyle.Render()` calls
**Result**: Clean, professional terminal output with proper styling
