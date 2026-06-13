# Reproduction Report

## Command

`npm test`

- Exit code: `1`
- Duration: `1.842s`

## Expected Behavior

The test suite should pass after installing dependencies from the committed lockfile.

## Actual Behavior

The command failed with a module resolution error.

## Environment

- arch: `arm64`
- node: `v22.2.0`
- npm: `10.8.1`
- os: `darwin`

## Dependency State

- package-lock.json: `48291 bytes`

## Logs

```text
FAIL src/app.test.ts
Error: Cannot find module '@example/missing-package'
Require stack:
- /workspace/project/src/app.ts
- /workspace/project/src/app.test.ts

Authorization: Bearer [REDACTED:bearer-token]
```

## Reproduction Confidence

This report includes command, environment, dependency state, and sanitized logs. A maintainer should be able to ask for the expected behavior only if it is still ambiguous.
