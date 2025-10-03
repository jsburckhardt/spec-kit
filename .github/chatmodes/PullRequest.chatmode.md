---
description: Custom mode focusing on creating a Pull Request using GitHub CLI.
tools: ['runCommands']
---
GOAL
- Create a Pull Request using the GitHub CLI (gh).
- Compare the current working branch against main.
- Generate a temporary file with the PR description.
- Use the temporary file to create the PR.

ASSUMPTIONS
- You are in a git repository with a valid remote (origin).
- gh CLI is installed and authenticated.
- Current branch is different from main.

STEPS
1. Identify current branch
   - `git rev-parse --abbrev-ref HEAD`
   - Verify it is not "main".

2. Fetch and sync main
   - `git fetch origin`
   - Ensure local main exists and tracks origin/main.

3. Show diff with main
   - `git diff --stat main...HEAD`

4. Prepare PR description
   - Create TEMP file path.
   - Write PR title and body into TEMP file.

5. Create Pull Request
   - `gh pr create --base main --head <CURRENT_BRANCH> --title "<PR_TITLE>" --body-file TEMP`

6. Cleanup
   - Delete TEMP file after PR creation.

OUTPUTS
- pr_url: URL of created PR.
- pr_number: Number of created PR.
- diff_stat: Output from diff summary.

ERROR HANDLING
- If gh is not authenticated → instruct to run `gh auth login`.
- If current branch == main → fail with message.
- If no diff vs main → warn and abort.
