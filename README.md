# ReproGate

ReproGate turns low-signal bug reports into paste-ready reproduction reports.

It is a small Go CLI plus MCP server for developers and open source maintainers:

- `reprogate capture -- <command>` runs a failing command and writes `.reprogate/repro.md`.
- `reprogate redact` removes likely secrets from logs.
- `reprogate ready-check` checks whether issue or PR text has reproduction steps, environment details, and logs.
- `reprogate mcp` exposes safe ReproGate tools to AI coding clients over stdio.

## Install

```bash
go install github.com/JinRudy/reprogate/cmd/reprogate@latest
```

For local development:

```bash
go run ./cmd/reprogate help
```

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

The output is JSON:

```json
{
  "labels": ["needs-repro", "missing-env", "missing-log"],
  "missing": ["reproduction steps", "environment details", "logs or command output"],
  "summary": "labels: needs-repro, missing-env, missing-log; missing: reproduction steps, environment details, logs or command output"
}
```

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
      - uses: JinRudy/reprogate@main
```

## Development

```bash
go test ./...
go run ./cmd/reprogate capture -- go env GOVERSION
go run ./cmd/reprogate ready-check < .reprogate/repro.md
```
