# ReproGate

[![CI](https://github.com/JinRudy/reprogate/actions/workflows/reprogate.yml/badge.svg)](https://github.com/JinRudy/reprogate/actions/workflows/reprogate.yml)
[![Action self-test](https://github.com/JinRudy/reprogate/actions/workflows/action-self-test.yml/badge.svg)](https://github.com/JinRudy/reprogate/actions/workflows/action-self-test.yml)
[![Release](https://img.shields.io/github/v/release/JinRudy/reprogate?sort=semver)](https://github.com/JinRudy/reprogate/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ReproGate checks issues and PRs for reproduction steps, environment details, and logs.

It helps maintainers stop repeating the same question:

> Can you provide reproduction steps, environment details, and logs?

It ships as a GitHub Action for issue and PR intake, plus a small CLI for reporters who need to generate a paste-ready reproduction report.

## Add The Action

```yaml
- uses: JinRudy/reprogate@v0.1.9
```

Use it on new or edited issues and pull requests:

```yaml
name: reprogate
on:
  issues:
    types: [opened, edited]
  pull_request:
    types: [opened, edited, synchronize]

jobs:
  ready-check:
    runs-on: ubuntu-latest
    steps:
      - id: reprogate
        uses: JinRudy/reprogate@v0.1.9
      - run: echo "${{ steps.reprogate.outputs.summary }}"
```

To have ReproGate update one issue or pull request comment when evidence is missing:

```yaml
- uses: JinRudy/reprogate@v0.1.9
  with:
    comment-on-missing: "true"
```

This requires `issues: write` permission in the caller workflow.

## What It Checks

ReproGate looks for the evidence maintainers usually need before a bug is actionable:

- reproduction steps
- environment details
- logs or command output

When evidence is missing, the Action exposes outputs such as `ready`, `labels`, `missing`, `missing_count`, and `summary` so your workflow can label, route, fail, or just report the result.

Need automated comments or labels? See the [GitHub Actions recipes](docs/recipes/github-actions.md).

## Fast Repository Setup

Generate the workflow and a matching bug report form:

```bash
reprogate init github-action
reprogate init issue-template
```

These commands write:

- `.github/workflows/reprogate.yml`
- `.github/ISSUE_TEMPLATE/bug_report.yml`

## Reporter CLI

Reporters can generate a sanitized Markdown report from the failing command:

```bash
reprogate capture -- npm test
```

That writes `.reprogate/repro.md`, ready to paste into a GitHub issue, Stack Overflow question, or maintainer discussion.

## Other Entry Points

- `reprogate ready-check`: check issue or PR text from a local file or stdin.
- `reprogate redact`: remove likely secrets before sharing logs.
- `reprogate mcp`: let AI coding tools redact text and check issue quality over stdio.

## Install

macOS and Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/JinRudy/reprogate/main/scripts/install.sh | sh
```

Install a pinned version or custom directory:

```bash
curl -fsSL https://raw.githubusercontent.com/JinRudy/reprogate/main/scripts/install.sh | REPROGATE_VERSION=v0.1.9 BIN_DIR="$HOME/bin" sh
```

Go users can also install from source:

```bash
go install github.com/JinRudy/reprogate/cmd/reprogate@latest
```

For local development:

```bash
go run ./cmd/reprogate help
```

## 60-Second Demo

```bash
curl -fsSL https://raw.githubusercontent.com/JinRudy/reprogate/main/scripts/install.sh | sh
reprogate capture -- go test ./...
cat .reprogate/repro.md
```

Example report: [docs/examples/repro.md](docs/examples/repro.md)

Demo issue: [#1 shows how a low-signal report is flagged](https://github.com/JinRudy/reprogate/issues/1).

## Initialize A Repository

Generate a ready-to-use workflow in the current repository:

```bash
reprogate init github-action
```

Generate a matching bug report issue form:

```bash
reprogate init issue-template
```

By default these write:

```text
.github/workflows/reprogate.yml
.github/ISSUE_TEMPLATE/bug_report.yml
```

If the target file already exists, ReproGate leaves it untouched unless you pass `--force`.

## Capture A Reproduction Report

```bash
reprogate capture -- npm test
```

By default this writes:

```text
.reprogate/repro.md
```

The report includes:

- command and exit code
- OS and architecture
- detected Go and Docker versions when available
- dependency lockfile summary
- sanitized logs
- an explicit expected-behavior section for the reporter to fill in

## Redact Logs

```bash
printf 'Authorization: Bearer abc123\npassword=hunter2\n' | reprogate redact
```

Example output:

```text
Authorization: Bearer [REDACTED:bearer-token]
password=[REDACTED:secret-value]
```

## Check Issue Or PR Readiness

```bash
reprogate ready-check --issue-body issue.md
```

or:

```bash
cat issue.md | reprogate ready-check
```

The output is JSON with labels and missing evidence.

## MCP

Start the stdio MCP server:

```bash
reprogate mcp
```

Example MCP client config:

```json
{
  "mcpServers": {
    "reprogate": {
      "command": "reprogate",
      "args": ["mcp"]
    }
  }
}
```

MCP tools:

- `redact_text`: redact likely secrets, credentials, and private paths from text.
- `check_issue`: check whether issue or PR text has reproduction steps, environment details, and logs.

AI client setup recipes: [docs/recipes/ai-clients.md](docs/recipes/ai-clients.md).

The MVP deliberately does not expose a command-execution MCP tool. Letting an AI run arbitrary local commands needs an explicit allowlist and confirmation model.

## GitHub Action

Minimal workflow:

```yaml
name: reprogate
on:
  issues:
    types: [opened, edited]
  pull_request:
    types: [opened, edited, synchronize]
jobs:
  ready-check:
    runs-on: ubuntu-latest
    steps:
      - id: reprogate
        uses: JinRudy/reprogate@v0.1.9
      - run: echo "${{ steps.reprogate.outputs.summary }}"
```

Strict mode fails the workflow when required evidence is missing:

```yaml
- uses: JinRudy/reprogate@v0.1.9
  with:
    fail-on-missing: "true"
```

Action outputs:

| Output | Description |
| --- | --- |
| `labels` | Comma-separated labels such as `needs-repro,missing-env,missing-log`. |
| `missing` | Comma-separated missing evidence fields. |
| `missing_count` | Number of missing evidence fields. |
| `ready` | `true` when the issue or pull request is review-ready. |
| `summary` | Human-readable readiness summary. |
| `result_json` | Full readiness result as JSON. |

Action inputs:

| Input | Default | Description |
| --- | --- | --- |
| `event-path` | `$GITHUB_EVENT_PATH` | Path to the GitHub event JSON file. Relative paths are resolved from the caller workspace. |
| `fail-on-missing` | `false` | Exit non-zero when reproduction steps, environment details, or logs are missing. |
| `comment-on-missing` | `false` | Post or update one issue or pull request comment when reproduction evidence is missing. Requires `issues: write`. |
| `go-version` | `1.25` | Go version used by the composite action to run ReproGate. |

More copy-paste workflows: [docs/recipes/github-actions.md](docs/recipes/github-actions.md).

## Development

```bash
go test ./...
go run ./cmd/reprogate capture -- go env GOVERSION
go run ./cmd/reprogate ready-check < .reprogate/repro.md
```

## License

MIT
