# Proposal: Add Build Timestamp File

## Summary

Generate a `build.json` file during each site build containing the build timestamp in RFC 3339 format. This file enables external monitoring of rclone sync status by providing a machine-readable timestamp that can be compared against current time.

## Motivation

Dee needs to verify that rclone syncing to Fastmail is working correctly. A build timestamp file that changes with every build will:
1. Always be synced (since it's new each build)
2. Provide a parseable timestamp for comparison
3. Enable an external application to detect stale syncs

## Format Decision

**JSON chosen over plain text because:**
- Standardized parsing in all languages
- Unambiguous RFC 3339 datetime format
- Extensible for future metadata needs
- Simpler consumption by monitoring application

## Output

File: `build.json` in output root

```json
{
  "buildTime": "2026-01-25T14:30:00Z"
}
```

## Scope

- **In scope**: Generate `build.json` with UTC build timestamp
- **Out of scope**: Monitoring/comparison logic (handled by external application)

## Impact

- Adds one new file to build output
- Extends existing site-building capability
- No breaking changes

## Related Capabilities

- `site-building` - Extends build output with new metadata file
