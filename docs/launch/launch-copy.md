# ReproGate Launch Copy

## One-Line Positioning

ReproGate turns "please provide a reproduction" into a one-command workflow.

## Short Description

ReproGate is a Go CLI, GitHub Action, and MCP server that helps developers generate sanitized reproduction reports and helps maintainers check whether issues or PRs include enough evidence to review.

## Show HN Draft

Title:

```text
Show HN: ReproGate - one command to generate reproducible bug reports
```

Post:

```text
Hi HN,

I built ReproGate because maintainers often have to ask the same questions on bug reports: what command failed, what OS/runtime was used, where are the logs, and can you provide a minimal reproduction?

ReproGate is a small Go CLI:

  reprogate capture -- npm test

It runs the command, captures exit code/stdout/stderr/runtime/dependency context, redacts likely secrets, and writes a paste-ready Markdown report.

It also has:
- reprogate redact
- reprogate ready-check
- reprogate mcp, exposing redact_text and check_issue over stdio
- a GitHub Action for issue/PR readiness checks

The first release intentionally does not let the MCP server execute arbitrary local commands. For now it only exposes safe text tools.

I am looking for feedback from maintainers: what evidence do you ask for most often when an issue is not reproducible?
```

## V2EX Draft

Title:

```text
做了个 CLI：一条命令生成可复现 bug 报告，减少维护者反复追问
```

Post:

```text
最近看了很多 GitHub issue，发现维护者经常卡在同一个问题：用户说“跑不起来 / 报错了”，但没有复现步骤、环境、命令输出、日志和依赖状态。

我做了一个小工具 ReproGate：

  reprogate capture -- npm test

它会运行命令，收集退出码、stdout/stderr、OS/架构、runtime、lockfile 摘要，并把可能的 token/password/path 做脱敏，生成一个可以直接贴到 issue 里的 Markdown 报告。

还有几个入口：
- reprogate redact：单独脱敏日志
- reprogate ready-check：检查 issue/PR 是否缺复现信息、环境、日志
- reprogate mcp：给 AI coding 工具接入，当前只提供 redact_text 和 check_issue 两个安全工具
- GitHub Action：可放到开源仓库做 issue/PR readiness check

我现在更想找维护者反馈：你们处理不可复现 issue 时，最想自动收集哪些信息？
```

## Maintainer Outreach Template

```text
Hi, I am building ReproGate, a small CLI/GitHub Action for reducing back-and-forth on unreproducible issues.

It generates a sanitized reproduction report from a failing command and can check whether issue/PR text includes reproduction steps, environment details, and logs.

I noticed your project has to ask for reproductions on some issues. If I adapt the output format to your issue template, would this be useful for your maintainers?

Repo: https://github.com/JinRudy/reprogate
```

## Submission Checklist

- GitHub topics are set.
- `v0.1.2` release includes the Marketplace-ready action and downloadable binaries.
- README first screen shows the problem, one command, and sample output.
- Demo issue exists: https://github.com/JinRudy/reprogate/issues/1
- Example reproduction report is linked from README.
- GitHub Actions is green on the release commit.
