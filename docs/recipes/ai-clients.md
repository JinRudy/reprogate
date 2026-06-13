# AI Client Recipes

ReproGate can run as a stdio MCP server for AI coding tools that support MCP
server configuration.

## Basic MCP Config

Install ReproGate first:

```bash
curl -fsSL https://raw.githubusercontent.com/JinRudy/reprogate/main/scripts/install.sh | sh
```

Then add this server entry to your MCP client config:

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

## What The AI Can Do

ReproGate exposes two MCP tools:

| Tool | Use case |
| --- | --- |
| `redact_text` | Remove likely secrets, credentials, and private paths before sharing logs. |
| `check_issue` | Check whether issue or PR text includes reproduction steps, environment details, and logs. |

## Prompt Examples

Redact a stack trace before posting it:

```text
Use ReproGate to redact this log before I paste it into a GitHub issue:

<paste log here>
```

Check whether an issue is ready for a maintainer:

```text
Use ReproGate to check whether this issue has enough reproduction evidence:

<paste issue body here>
```

Turn a vague failure report into a checklist:

```text
Use ReproGate to identify what evidence is missing from this bug report, then
rewrite the missing items as a short checklist for the reporter.
```

## Safety Boundary

The MCP server does not expose command execution. It only redacts provided text
and checks issue quality, so maintainers can connect it without giving an AI
tool permission to run arbitrary local commands through ReproGate.
