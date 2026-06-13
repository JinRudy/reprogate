# ReproGate

[![CI](https://github.com/JinRudy/reprogate/actions/workflows/reprogate.yml/badge.svg)](https://github.com/JinRudy/reprogate/actions/workflows/reprogate.yml)
[![Action self-test](https://github.com/JinRudy/reprogate/actions/workflows/action-self-test.yml/badge.svg)](https://github.com/JinRudy/reprogate/actions/workflows/action-self-test.yml)
[![Release](https://img.shields.io/github/v/release/JinRudy/reprogate?sort=semver)](https://github.com/JinRudy/reprogate/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

ReproGate turns "please provide a reproduction" into a one-command workflow.

It is a small Go CLI, GitHub Action, and MCP server for maintainers who are tired of asking for the same missing evidence in bug reports:

- reproduction steps
- environment details
- command output
- sanitized logs
- dependency and runtime context

```bash
reprogate capture -- npm test
```

That writes `.reprogate/repro.md`, ready to paste into a GitHub issue, Stack Overflow question, or maintainer discussion.

For maintainers, add the readiness check to a repository in one command:

```bash
reprogate init github-action
reprogate init issue-template
```

That writes `.github/workflows/reprogate.yml`, so new issues and PRs can be checked for reproduction steps, environment details, and logs.
It can also write `.github/ISSUE_TEMPLATE/bug_report.yml`, so reporters are prompted for the same evidence before opening an issue.

## Why It Exists

Maintainers lose time on issues that say "it fails" but omit the command, OS, runtime, logs, lockfile state, and expected behavior. ReproGate gives reporters a repeatable way to collect that evidence before the maintainer has to ask.

For maintainers, ReproGate can also check issue and PR text:

```bash
cat issue.md | reprogate ready-check
```

Example output:

```json
{
  "labels": ["needs-repro", "missing-env", "missing-log"],
  "missing": ["reproduction steps", "environment details", "logs or command output"],
  "summary": "labels: needs-repro, missing-env, missing-log; missing: reproduction steps, environment details, logs or command output"
}
```

## What It Does

| Surface | Command | Use case |
| --- | --- | --- |
| Capture | `reprogate capture -- <command>` | Run a failing command and generate a sanitized reproduction report. |
| Redact | `reprogate redact` | Remove likely secrets before pasting logs into an issue. |
| Ready check | `reprogate ready-check` | Check whether an issue or PR has enough evidence to review. |
| Init | `reprogate init github-action` | Generate a ready-to-use GitHub Actions workflow in the current repository. |
| Init | `reprogate init issue-template` | Generate a GitHub bug report form that asks for ReproGate evidence. |
| MCP | `reprogate mcp` | Let AI coding tools redact logs and check issue quality over stdio. |
| GitHub Action | `uses: JinRudy/reprogate@v0.1.6` | Add readiness checks to issue and PR workflows. |

## Install

macOS and Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/JinRudy/reprogate/main/scripts/install.sh | sh
```

Install a pinned version or custom directory:

```bash
curl -fsSL https://raw.githubusercontent.com/JinRudy/reprogate/main/scripts/install.sh | REPROGATE_VERSION=v0.1.6 BIN_DIR="$HOME/bin" sh
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
        uses: JinRudy/reprogate@v0.1.6
      - run: echo "${{ steps.reprogate.outputs.summary }}"
```

Strict mode fails the workflow when required evidence is missing:

```yaml
- uses: JinRudy/reprogate@v0.1.6
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
