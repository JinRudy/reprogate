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
        uses: JinRudy/reprogate@v0.1.8
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
      - uses: JinRudy/reprogate@v0.1.8
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
        uses: JinRudy/reprogate@v0.1.8

  route:
    needs: ready-check
    if: needs.ready-check.outputs.ready != 'true'
    runs-on: ubuntu-latest
    steps:
      - run: |
          echo "Suggested labels: ${{ needs.ready-check.outputs.labels }}"
          echo "Missing fields: ${{ needs.ready-check.outputs.missing_count }}"
```

## Comment And Label Missing Evidence

Use this when you want ReproGate to leave one maintainer-style comment on issues
that are missing reproduction details. The workflow updates its previous comment
instead of adding a new one on every edit.

It only applies labels that already exist in the repository, so the workflow does
not fail if `needs-repro`, `missing-env`, or `missing-log` has not been created.

```yaml
name: reprogate-intake
on:
  issues:
    types: [opened, edited]

permissions:
  contents: read
  issues: write

jobs:
  intake:
    runs-on: ubuntu-latest
    steps:
      - id: reprogate
        uses: JinRudy/reprogate@v0.1.8

      - name: Apply existing evidence labels
        if: steps.reprogate.outputs.ready != 'true'
        uses: actions/github-script@v9
        env:
          REPROGATE_LABELS: ${{ steps.reprogate.outputs.labels }}
        with:
          script: |
            const suggested = process.env.REPROGATE_LABELS
              .split(',')
              .map((label) => label.trim())
              .filter(Boolean);

            const existing = await github.paginate(github.rest.issues.listLabelsForRepo, {
              owner: context.repo.owner,
              repo: context.repo.repo,
              per_page: 100,
            });
            const existingNames = new Set(existing.map((label) => label.name));
            const labels = suggested.filter((label) => existingNames.has(label));

            if (labels.length === 0) {
              return;
            }
            await github.rest.issues.addLabels({
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              labels,
            });

      - name: Ask for missing evidence
        if: steps.reprogate.outputs.ready != 'true'
        uses: actions/github-script@v9
        env:
          REPROGATE_MISSING: ${{ steps.reprogate.outputs.missing }}
        with:
          script: |
            const marker = '<!-- reprogate-missing-evidence -->';
            const missing = process.env.REPROGATE_MISSING
              .split(',')
              .map((item) => item.trim())
              .filter(Boolean);
            const body = [
              marker,
              'Thanks for the report. ReproGate found that this issue is missing:',
              '',
              ...missing.map((item) => `- ${item}`),
              '',
              'Please add the missing details, or generate a paste-ready report with:',
              '',
              '```bash',
              'reprogate capture -- <failing command>',
              '```',
            ].join('\n');

            const comments = await github.paginate(github.rest.issues.listComments, {
              owner: context.repo.owner,
              repo: context.repo.repo,
              issue_number: context.issue.number,
              per_page: 100,
            });
            const previous = comments.find((comment) =>
              comment.user?.type === 'Bot' && comment.body?.includes(marker)
            );

            if (previous) {
              await github.rest.issues.updateComment({
                owner: context.repo.owner,
                repo: context.repo.repo,
                comment_id: previous.id,
                body,
              });
            } else {
              await github.rest.issues.createComment({
                owner: context.repo.owner,
                repo: context.repo.repo,
                issue_number: context.issue.number,
                body,
              });
            }
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
