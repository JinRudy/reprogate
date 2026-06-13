# GitHub Actions Recipes

These recipes are copy-paste starting points for maintainers who want fewer
"please provide a reproduction" follow-ups.

To generate the default workflow from the CLI instead, run:

```bash
reprogate init github-action
```

To add a matching GitHub bug report form, run:

```bash
reprogate init issue-template
```

## Summary-Only Intake Check

Use this when you want a low-risk first rollout. It never fails the workflow and
only writes the ReproGate result to the job summary.

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

## Strict Intake Check

Use this when missing reproduction steps, environment details, or logs should
block the check.

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
      - uses: JinRudy/reprogate@v0.1.6
        with:
          fail-on-missing: "true"
```

## Route Existing Automation From Outputs

Use this when you already have a labeler, triage bot, or notification workflow.
ReproGate exposes stable outputs so your workflow can decide what to do next.

```yaml
name: reprogate
on:
  issues:
    types: [opened, edited]

jobs:
  ready-check:
    runs-on: ubuntu-latest
    outputs:
      ready: ${{ steps.reprogate.outputs.ready }}
      labels: ${{ steps.reprogate.outputs.labels }}
      missing_count: ${{ steps.reprogate.outputs.missing_count }}
    steps:
      - id: reprogate
        uses: JinRudy/reprogate@v0.1.6

  route:
    needs: ready-check
    if: needs.ready-check.outputs.ready != 'true'
    runs-on: ubuntu-latest
    steps:
      - run: |
          echo "Suggested labels: ${{ needs.ready-check.outputs.labels }}"
          echo "Missing fields: ${{ needs.ready-check.outputs.missing_count }}"
```

## Outputs

| Output | Description |
| --- | --- |
| `labels` | Comma-separated labels such as `needs-repro,missing-env,missing-log`. |
| `missing` | Comma-separated missing evidence fields. |
| `missing_count` | Number of missing evidence fields. |
| `ready` | `true` when the issue or pull request is review-ready. |
| `summary` | Human-readable readiness summary. |
| `result_json` | Full readiness result as JSON. |
