package setup

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const defaultGitHubActionPath = ".github/workflows/reprogate.yml"

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
        uses: JinRudy/reprogate@v0.1.5
      - run: echo "${{ steps.reprogate.outputs.summary }}"
`

func RunCLI(args []string, out io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: reprogate init github-action [--path %s] [--force]", defaultGitHubActionPath)
	}
	if args[0] != "github-action" {
		return fmt.Errorf("unknown init target %q", args[0])
	}
	return runGitHubAction(args[1:], out)
}

func runGitHubAction(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("init github-action", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	targetPath := fs.String("path", defaultGitHubActionPath, "workflow path")
	force := fs.Bool("force", false, "overwrite an existing workflow")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if _, err := os.Stat(*targetPath); err == nil && !*force {
		return fmt.Errorf("%s already exists; rerun with --force to overwrite it", *targetPath)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(*targetPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(*targetPath, []byte(githubActionWorkflow), 0o644); err != nil {
		return err
	}
	fmt.Fprintf(out, "Wrote %s\n", *targetPath)
	return nil
}
