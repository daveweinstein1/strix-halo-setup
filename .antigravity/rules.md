# Project Coding Rules

## Standards Auto-Loading

This project uses model and language-specific coding standards that are stored globally.

At the start of each coding session, the agent will:
1. Detect which AI model is running (Gemini 3 Pro, Claude Sonnet 4.5, or Claude Opus 4)
2. Detect the programming language (Python or Rust)
3. Automatically load the appropriate standards file from **global location**: `%USERPROFILE%\.gemini\docs\`

## Standards Files Location

**IMPORTANT:** Standards files are in your home directory (global, apply to all projects):

- **Python + Claude**: `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_CLAUDE_PYTHON.md`
- **Python + Gemini**: `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_GEMINI_PYTHON.md`
- **Rust + Claude**: `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_CLAUDE_RUST.md`
- **Rust + Gemini**: `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_GEMINI_RUST.md`

On Windows, `%USERPROFILE%` expands to `C:\Users\YourUsername\`

## Verification Required

Before marking any task complete, the agent must verify:
- Code passes all standards checks (mypy, ruff, clippy, etc.)
- Tests are written and passing
- Documentation is complete
- No debug statements remain

## Project-Specific Overrides

Add any project-specific rules below that override or extend the global standards:

<!-- Your project-specific rules here -->