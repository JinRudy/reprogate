package setup

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const defaultGitHubActionPath = ".github/workflows/reprogate.yml"
const defaultIssueTemplatePath = ".github/ISSUE_TEMPLATE/bug_report.yml"

const githubActionWorkflow = `name: reprogate
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
`

const issueTemplate = `name: Bug report
description: Report a reproducible bug with ReproGate evidence.
title: "[Bug]: "
body:
  - type: markdown
    attributes:
      value: |
        Please include a ReproGate report when possible:

        reprogate capture -- <failing command>

        Then paste .reprogate/repro.md below.
  - type: textarea
    id: reprogate-report
    attributes:
      label: ReproGate report
      description: Paste .reprogate/repro.md here, or explain why you could not generate one.
      placeholder: |
        ## Reproduction Report
        Command: ...
    validations:
      required: false
  - type: textarea
    id: steps
    attributes:
      label: Steps to reproduce
      description: List the smallest steps that trigger the bug.
      placeholder: |
        1. ...
        2. ...
    validations:
      required: true
  - type: input
    id: environment
    attributes:
      label: Environment
      description: OS, runtime, package manager, framework, and relevant versions.
      placeholder: macOS 15, Node 22, pnpm 10
    validations:
      required: true
  - type: textarea
    id: logs
    attributes:
      label: Logs or command output
      description: Paste sanitized logs, stack traces, or command output.
      render: text
    validations:
      required: true
`

func RunCLI(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: reprogate init <github-action|issue-template> [--path PATH] [--force]")
	}
	switch args[0] {
	case "github-action":
		return runGitHubAction(args[1:], out)
	case "issue-template":
		return runIssueTemplate(args[1:], out)
	default:
		return fmt.Errorf("unknown init target %q", args[0])
	}
}

func runGitHubAction(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("init github-action", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	targetPath := fs.String("path", defaultGitHubActionPath, "workflow path")
	force := fs.Bool("force", false, "overwrite an existing workflow")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return writeFile(*targetPath, githubActionWorkflow, *force, out)
}

func runIssueTemplate(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("init issue-template", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	targetPath := fs.String("path", defaultIssueTemplatePath, "issue template path")
	force := fs.Bool("force", false, "overwrite an existing issue template")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return writeFile(*targetPath, issueTemplate, *force, out)
}

func writeFile(targetPath string, content string, force bool, out io.Writer) error {
	if _, err := os.Stat(targetPath); err == nil && !force {
		return fmt.Errorf("%s already exists; rerun with --force to overwrite it", targetPath)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(targetPath, []byte(content), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(out, "Wrote %s\n", targetPath)
	return nil
}
