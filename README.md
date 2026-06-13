# ReproGate

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
| MCP | `reprogate mcp` | Let AI coding tools redact logs and check issue quality over stdio. |
| GitHub Action | `uses: JinRudy/reprogate@v0.1.0` | Add readiness checks to issue and PR workflows. |

## Install

```bash
go install github.com/JinRudy/reprogate/cmd/reprogate@latest
```

For local development:

```bash
go run ./cmd/reprogate help
```

## 60-Second Demo

```bash
go install github.com/JinRudy/reprogate/cmd/reprogate@latest
reprogate capture -- go test ./...
cat .reprogate/repro.md
```

Example report: [docs/examples/repro.md](docs/examples/repro.md)

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
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - uses: JinRudy/reprogate@v0.1.0
```

## Development

```bash
go test ./...
go run ./cmd/reprogate capture -- go env GOVERSION
go run ./cmd/reprogate ready-check < .reprogate/repro.md
```
