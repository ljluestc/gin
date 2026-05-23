# docs(readme): add Go module troubleshooting for "module found but does not contain package" errors (#2283)

## Related Issue

Closes #2283

## Problem

Users reported startup/build failures like:

- `module github.com/gin-gonic/gin@latest found (...), but does not contain package github.com/gin-gonic/gin`
- `module golang.org/x/sync@latest found (...), but does not contain package golang.org/x/sync/errgroup`

The issue report is environment-specific (module cache/proxy/module-path setup), but currently there is no direct troubleshooting guidance in the README installation flow for this error pattern.

## Solution

Add a focused **Go modules troubleshooting** subsection under `README.md` installation.

The new section:

1. Shows the exact error pattern users are likely to see.
2. Provides immediate recovery commands:
   - `go clean -modcache`
   - `go mod tidy`
   - `go mod download`
3. Adds checks for common root causes:
   - incorrect `module` path in consumer `go.mod`
   - unsupported Go version
   - unreachable/misconfigured `GOPROXY`
4. Tells users what diagnostic data to provide when opening a bug report.

## Files Changed

- `README.md`
- `PR_DESCRIPTION_2283_EN.md` (this local PR draft)

## Why this approach

- Keeps fix low risk and immediately useful for users blocked by module resolution errors.
- Avoids introducing framework code changes for an issue that is not reproducible as a Gin runtime bug.
- Improves triage quality by asking for actionable environment/module details.

## Testing

- Documentation change only (no runtime behavior changes).
- Verified markdown structure renders correctly and commands are valid Go module commands.

## Backward Compatibility

- Fully backward compatible.
- No API changes, no behavior changes, no dependency changes.

## Security Impact

- None.

## Summary

- **Problem**: Users can hit confusing module-resolution errors and do not have targeted guidance in README.
- **What changed**: Added a troubleshooting subsection in installation docs with recovery commands and validation checks.
- **What did NOT change**: Gin framework code, public APIs, and runtime behavior.
