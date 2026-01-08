---
description: Detect current AI model and load appropriate coding standards
---

# Model Detection and Standards Loading

## Purpose

Detect which AI model is running and load the correct coding standards for the current programming language.

## Steps

### 1. Identify Current Model

Determine which model Antigravity is using:
- Gemini 3 Pro (default)
- Claude Sonnet 4.5
- Claude Opus 4
- Other

**Method**: Check Antigravity settings or ask user to confirm

### 2. Identify Programming Language

From the current task or workspace, determine the language:
- Check file extensions (`.py`, `.rs`)
- Check user's request for language hints
- Check existing project structure

### 3. Load Appropriate Standards File

Based on model + language combination:

| Model | Python | Rust |
|-------|--------|------|
| Claude (Sonnet/Opus) | `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_CLAUDE_PYTHON.md` | `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_CLAUDE_RUST.md` |
| Gemini 3 Pro | `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_GEMINI_PYTHON.md` | `%USERPROFILE%\.gemini\docs\CODING_STANDARDS_GEMINI_RUST.md` |

**Note:** Standards files are global (in home directory), apply to all projects.

### 4. Confirm to User

Output to user:
```
✓ Detected [Model Name]
✓ Loading [Language] coding standards from %USERPROFILE%\.gemini\docs\[filename]
✓ Key rules active:
  - [Most important rule 1]
  - [Most important rule 2]
  - [Most important rule 3]
```

### 5. Set Verification Reminder

Remember to verify ALL standards compliance before marking any task complete.

## Example Output

```
✓ Detected Claude Sonnet 4.5
✓ Loading Python coding standards from %USERPROFILE%\.gemini\docs\CODING_STANDARDS_CLAUDE_PYTHON.md
✓ Key rules active:
  - All functions must have type hints (run mypy to verify)
  - Use orjson for JSON (not json module)
  - Use polars for dataframes (not pandas)
```