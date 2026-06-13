# ReproGate Design

## Problem

Open source maintainers repeatedly lose time on issues and pull requests that cannot be reproduced or reviewed quickly. The recurring missing pieces are reproduction steps, environment details, command output, dependency state, logs, and a concise explanation of what changed.

ReproGate turns those missing pieces into machine-checkable artifacts:

- A local CLI captures a failing command and writes a sanitized Markdown reproduction report.
- A GitHub Action checks issues and pull requests for review readiness signals.
- The project stays deterministic and small; it does not try to detect whether text or code was AI-generated.

## Target Users

- Developers filing bug reports who want a one-command way to collect the right evidence.
- Open source maintainers who want fewer low-signal issues and pull requests.
- Contributors who want to prove that a PR is tested, scoped, and connected to a real problem.

## MVP Scope

The first release includes four user-visible surfaces.

1. `reprogate capture -- <command>`
   - Runs the command.
   - Captures exit code, stdout, stderr, duration, working directory name, Git commit, Git dirty status, OS, architecture, shell, detected runtimes, package manager versions, lockfile hashes, Docker version, Compose version, and listening port summary.
   - Writes `.reprogate/repro.md` by default.

2. Sanitization
   - Redacts likely secrets in logs, environment values, URLs, and paths.
   - Keeps environment key names but hides values.
   - Marks each redaction as `[REDACTED:<kind>]` so users can see what was removed.

3. Markdown report
   - Produces sections for summary, command, expected behavior, actual behavior, environment, dependency state, logs, and reproduction confidence.
   - Leaves expected behavior as an explicit user-editable prompt when it cannot be inferred.
   - Keeps output paste-ready for GitHub issues, Stack Overflow, and developer forums.

4. GitHub Action readiness check
   - Checks issue bodies and pull requests for reproduction steps, environment details, logs, linked issue, test evidence, large diff risk, lockfile changes, generated file changes, and dependency changes.
   - Emits a Markdown summary.
   - Optionally applies labels: `needs-repro`, `missing-env`, `missing-log`, `review-ready`, `risky-diff`.

## Non-Goals

- No AI PR detection.
- No hosted service.
- No dependency on OpenAI or any LLM for the MVP.
- No Docker deployment dashboard.
- No automatic issue closing.
- No uploading local files without an explicit user command.

## Architecture

ReproGate is a Go CLI with small packages.

- `cmd/reprogate`: command-line entrypoint.
- `internal/capture`: command execution and runtime probes.
- `internal/redact`: secret and path redaction.
- `internal/report`: Markdown report rendering.
- `internal/checks`: issue and pull request readiness checks.
- `internal/githubaction`: GitHub Actions input/output adapter.

The GitHub Action uses the same binary. This keeps CLI and CI behavior consistent and avoids duplicating readiness rules in YAML or JavaScript.

## Data Flow

### Local Capture

1. User runs `reprogate capture -- npm test`.
2. ReproGate records command metadata and starts the child process.
3. stdout and stderr are streamed to the terminal and captured in memory up to a configurable byte limit.
4. Environment and project probes run after the command completes.
5. Captured data is sanitized.
6. `.reprogate/repro.md` is written.
7. The CLI exits with the child command exit code by default.

### GitHub Action

1. Action receives issue or pull request event JSON.
2. ReproGate reads the event body, diff metadata, labels, and changed files.
3. Checks produce a readiness score and missing-evidence list.
4. Action writes a job summary.
5. If configured, Action applies labels through the GitHub token.

## Configuration

Optional `.reprogate.yml`:

```yaml
report:
  output: .reprogate/repro.md
  max_log_bytes: 200000
redaction:
  extra_patterns:
    - "company-internal-domain.example"
checks:
  large_diff_files: 20
  large_diff_lines: 800
  require_linked_issue: true
labels:
  enabled: true
```

If no config exists, ReproGate uses conservative defaults.

## Readiness Rules

The MVP readiness check is intentionally simple and explainable.

- `missing-repro`: no numbered steps, no command, and no reproduction repository link.
- `missing-env`: no OS/runtime/package manager/version details.
- `missing-log`: no log block, stack trace, command output, or screenshot link.
- `risky-diff`: changed files or lines exceed thresholds, lockfile changes are present, generated files changed, or dependency manifests changed.
- `review-ready`: issue or PR has reproduction details, environment details, logs or test output, and does not trigger high-risk thresholds.

Rules are heuristics. The tool reports evidence, not final merge decisions.

## Security And Privacy

- Redaction runs before reports are written.
- Raw captured logs are not persisted by default.
- Environment values are never printed unless explicitly allowed.
- Uploading a bundle is out of MVP scope.
- The GitHub Action only needs `issues: write` when labels are enabled. Without labels, read-only permissions are enough.

## Testing

MVP tests cover:

- Redaction of common secret forms: tokens, passwords, API keys, bearer headers, private URLs, and home paths.
- Command capture behavior for success, failure, timeout, and large logs.
- Report rendering with stable Markdown snapshots.
- Readiness checks for missing reproduction, missing environment, missing logs, risky diffs, and ready PRs.
- GitHub Action adapter using recorded event JSON fixtures.

## Release Plan

1. First commit: design spec only.
2. MVP implementation:
   - Go module and CLI skeleton.
   - `capture` command with Markdown report.
   - redaction package and tests.
   - readiness checks and tests.
   - GitHub Action wrapper.
3. Public README:
   - one-sentence value proposition.
   - install commands.
   - copy-paste examples.
   - sample `repro.md`.
   - sample GitHub Action workflow.
4. Dogfood:
   - enable ReproGate on its own repository.
   - open example issue and PR showing the generated summaries.

## Success Criteria

- A developer can run one command and paste a useful bug report in under one minute.
- A maintainer can glance at a ReproGate summary and know what evidence is missing.
- The project can run without a hosted backend or API key.
- The MVP is understandable from the README without reading implementation code.
