package main

import (
	"context"
	"fmt"
	"os"

	"github.com/JinRudy/reprogate/internal/capture"
	"github.com/JinRudy/reprogate/internal/checks"
	"github.com/JinRudy/reprogate/internal/mcpserver"
	"github.com/JinRudy/reprogate/internal/redact"
	"github.com/JinRudy/reprogate/internal/report"
	"github.com/JinRudy/reprogate/internal/setup"
)

func main() {
	if err := run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHelp()
		return nil
	}

	switch args[0] {
	case "capture":
		return capture.RunCLI(ctx, args[1:], os.Stdout, os.Stderr)
	case "redact":
		return redact.RunCLI(args[1:], os.Stdin, os.Stdout)
	case "report":
		return report.RunCLI(args[1:], os.Stdout)
	case "ready-check":
		return checks.RunCLI(args[1:], os.Stdin, os.Stdout)
	case "init":
		return setup.RunCLI(args[1:], os.Stdout)
	case "mcp":
		return mcpserver.Run(ctx, os.Stdin, os.Stdout)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printHelp() {
	fmt.Fprintln(os.Stdout, `ReproGate collects reproducible bug evidence.

Usage:
  reprogate capture -- <command> [args...]
  reprogate redact < file.log
  reprogate ready-check --issue-body issue.md
  reprogate init github-action
  reprogate init issue-template
  reprogate mcp`)
}
