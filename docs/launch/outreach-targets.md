# Maintainer Outreach Targets

Snapshot date: 2026-06-13

This is a manual outreach list, not a comment-spam queue. Use it to find
maintainers who already show a public "needs reproduction" workflow, then ask
for feedback in a tailored way.

## Search Method

Commands used:

```bash
gh search issues '"please provide a reproduction"' --state open --archived=false
gh search issues '"minimal reproduction"' --state open --archived=false
gh search issues '"steps to reproduce" "environment"' --state open --archived=false
gh search issues '"needs reproduction"' --state open --archived=false
```

Selection rules:

- Prefer repositories that already use `needs reproduction`, `awaiting submitter`, or similar labels.
- Prefer active, non-archived repositories where maintainers visibly ask for reproduction evidence.
- Avoid commenting on a specific bug unless we can add real help to that bug.
- Avoid bot-managed issue queues that already warn about duplicate or low-signal proposals.
- Start with maintainers or discussions, not drive-by comments on unrelated issues.

## First Batch

| Priority | Repository | Stars at snapshot | Evidence | Why it fits | Outreach posture |
| --- | --- | ---: | --- | --- | --- |
| 1 | [opennextjs/opennextjs-aws](https://github.com/opennextjs/opennextjs-aws) | 5.0k | [Issue #1059](https://github.com/opennextjs/opennextjs-aws/issues/1059) has `need reproduction`. | Mid-size maintainer-led project; reproducible deployment context is painful and specific. | Ask if an Action checking for command/env/log fields would fit their issue workflow. |
| 2 | [getsentry/sentry-javascript](https://github.com/getsentry/sentry-javascript) | 8.7k | [Issue #20132](https://github.com/getsentry/sentry-javascript/issues/20132) asks for runtime version and a reproduction repo. | SDK maintainers repeatedly need env/runtime context. | Ask for feedback on ReproGate report fields, not a generic promotion. |
| 3 | [sveltejs/svelte](https://github.com/sveltejs/svelte) | 87.2k | [Issue #17610](https://github.com/sveltejs/svelte/issues/17610) has `awaiting submitter`, described as needing reproduction or clarification. | Large project with many low-signal bug reports. | Do not pitch on the bug; look for maintainer discussion or issue template improvement channel. |
| 4 | [cloudflare/workers-sdk](https://github.com/cloudflare/workers-sdk) | 4.2k | [Issue #14268](https://github.com/cloudflare/workers-sdk/issues/14268) has `needs-reproduction`. | Runtime/deployment bugs often need command output, versions, config, and sanitized logs. | Good feedback target after adding Cloudflare-style examples. |
| 5 | [vitejs/vite](https://github.com/vitejs/vite) | 81.4k | [Issue #22662](https://github.com/vitejs/vite/issues/22662) has `needs reproduction`. | Very high visibility and already has a reproduction bot, useful benchmark. | Study their workflow first; do not cold pitch unless we can improve a clear gap. |
| Watch | [Expensify/App](https://github.com/Expensify/App) | 4.9k | [Issue #93329](https://github.com/Expensify/App/issues/93329) has `Needs Reproduction` and visible anti-duplicate automation. | Shows the pain and the spam risk at scale. | Do not comment. Use as a cautionary example for human-first outreach. |

## Short Outreach

Use this only after checking the repository's contribution norms.

```text
Hi, I am building ReproGate, a small CLI/GitHub Action for reducing the
"please provide a reproduction" loop.

I noticed this repo already has a needs-reproduction workflow. ReproGate can
generate a sanitized report from a failing command and can check whether an
issue includes reproduction steps, environment details, and logs.

Would a small Action like this be useful for your issue intake, or are there
fields your maintainers usually need that I should support first?

Repo: https://github.com/JinRudy/reprogate
```

## Safer First Touch

For large projects, use a feedback request rather than a direct tool pitch:

```text
Quick maintainer question: when an issue is missing reproduction details, what
fields do you usually need before it becomes actionable?

I am building ReproGate to generate/check that evidence automatically, and I am
trying to align the report format with real maintainer workflows before asking
projects to try it.
```

## Weekly Outreach Loop

1. Pick 3 repositories from the list.
2. Read their issue template, contributing guide, and recent maintainer comments.
3. Decide whether the right channel is a Discussion, a tooling issue, email, or no contact.
4. Send at most one tailored message per repository.
5. Track replies and requested fields before changing ReproGate.

## Sent Log

2026-06-13:

| Repository | Issue | Comment | Status | Follow-up rule |
| --- | --- | --- | --- | --- |
| `cloudflare/workers-sdk` | [#14268](https://github.com/cloudflare/workers-sdk/issues/14268) | [comment](https://github.com/cloudflare/workers-sdk/issues/14268#issuecomment-4697811134) | Sent | Wait for maintainer/reporter response before commenting again. |
| `getsentry/sentry-javascript` | [#20132](https://github.com/getsentry/sentry-javascript/issues/20132) | [comment](https://github.com/getsentry/sentry-javascript/issues/20132#issuecomment-4697811533) | Sent | Wait for maintainer/reporter response before commenting again. |
| `opennextjs/opennextjs-aws` | [#1059](https://github.com/opennextjs/opennextjs-aws/issues/1059) | [comment](https://github.com/opennextjs/opennextjs-aws/issues/1059#issuecomment-4697811870) | Sent | Wait for maintainer/reporter response before commenting again. |
