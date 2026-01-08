---
description: Verify code follows all coding standards before presenting
---

# Standards Verification Checklist

## Purpose

Verify that all code follows the loaded coding standards before marking any task complete or presenting code to the user.

## Python Code Verification

If working with Python code:

- [ ] Uses `orjson` for JSON operations (not `json` module)
- [ ] Uses `polars` for dataframes (not `pandas`)
- [ ] All functions have complete type hints
- [ ] All functions have docstrings with Args/Returns/Raises
- [ ] Proper error handling (no bare `except:`)
- [ ] No hardcoded secrets or API keys
- [ ] No debug `print()` statements
- [ ] Code formatted with `ruff format`
- [ ] `ruff check` passes with no errors
- [ ] `mypy` passes with no errors
- [ ] Tests are written and passing
- [ ] No mutable default arguments

## Rust Code Verification

If working with Rust code:

- [ ] No `.unwrap()` in production code paths
- [ ] All public items have doc comments
- [ ] Uses `thiserror` for library errors or `anyhow` for application errors
- [ ] Code compiles with zero warnings
- [ ] No unnecessary `unsafe` blocks
- [ ] No hardcoded secrets or API keys
- [ ] No debug `println!` or `dbg!` macros
- [ ] Code formatted with `rustfmt` (`cargo fmt`)
- [ ] `cargo clippy -- -D warnings` passes
- [ ] `cargo test` passes
- [ ] If PyO3 project: rebuilt with `maturin develop`
- [ ] If WASM project: rebuilt with `wasm-pack build`

## Universal Verification

For all languages:

- [ ] Tests written for all new functions
- [ ] Documentation complete and up-to-date
- [ ] Security requirements met (no secrets committed)
- [ ] No commented-out code
- [ ] No TODO/FIXME without tracking issues
- [ ] `.env` file in `.gitignore`

## Action Required

**If ANY checkbox is unchecked:**
1. STOP immediately
2. Fix the violation
3. Re-run verification
4. Only then mark task complete

## Usage

Run this workflow manually before finalizing code:
```
/verify-standards
```

Or use as a mental checklist before claiming any task is done.

## Evidence Required

Don't just claim standards are met—provide evidence:
- Show `mypy` output (zero errors)
- Show `ruff check` output (all passed)
- Show `cargo clippy` output (no warnings)
- Show test results (all passed)
- Provide screenshot if UI changes

**"I believe it's correct" < "Here's proof it's correct"**