# Fix: Copilot Prompt Filenames

## Problem
When running `gospecify init --here` with GitHub Copilot as the AI assistant, the generated files in `.github/prompts/` were created with incorrect filenames. They were missing the `.prompt.` part in the middle of the filename.

**Expected**: `.github/prompts/specify.prompt.md`
**Actual**: `.github/prompts/specify.md`

## Root Cause
The `processTemplates` function in `/src/gospecify/cmd/init.go` was copying command template files to the assistant directory without considering the assistant's specific file format requirements. It was using the original template filename (`specify.md`) instead of applying the format-specific naming convention.

## Solution
Added a new helper function `generateCommandFileName()` that:

1. Extracts the base name from the original template filename (e.g., `specify` from `specify.md`)
2. Applies the assistant's format to generate the correct filename:
   - **FormatPrompt** (Copilot): `baseName + ".prompt.md"`
   - **FormatTOML** (Gemini): `baseName + ".toml"`
   - **FormatMarkdown** (Claude): `baseName + ".md"`

## Files Modified
- `/src/gospecify/cmd/init.go`:
  - Added `generateCommandFileName()` helper function
  - Modified template copying logic to use the helper function

## Testing Results
✅ **Copilot**: Files created as `.prompt.md` (e.g., `specify.prompt.md`)
✅ **Claude**: Files created as `.md` (e.g., `specify.md`)
✅ **Gemini**: Files created as `.toml` (e.g., `specify.toml`)

## Impact
- **Backward Compatibility**: ✅ Existing agents (Claude, Gemini, etc.) continue to work correctly
- **New Functionality**: ✅ Copilot now generates files with correct `.prompt.md` extension
- **Code Quality**: ✅ Added proper format handling for future AI assistants
