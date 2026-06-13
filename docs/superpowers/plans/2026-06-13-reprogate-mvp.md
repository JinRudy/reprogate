# ReproGate MVP Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI plus stdio MCP server that generates sanitized reproduction reports and checks issue/PR readiness.

**Architecture:** The CLI and MCP server share the same internal packages. `capture` gathers command and environment evidence, `redact` sanitizes sensitive text, `report` renders Markdown, `checks` evaluates review readiness, and `mcpserver` exposes the same capabilities to AI tools over stdio.

**Tech Stack:** Go 1.25+, standard library, `github.com/modelcontextprotocol/go-sdk/mcp`, Markdown fixtures, GitHub Actions YAML.

---

## File Structure

- `go.mod`: Go module definition.
- `cmd/reprogate/main.go`: CLI entrypoint and command dispatch.
- `internal/redact/redact.go`: text and path sanitization.
- `internal/redact/redact_test.go`: redaction tests.
- `internal/report/report.go`: reproduction report data model and Markdown renderer.
- `internal/report/report_test.go`: report snapshot-style tests.
- `internal/capture/capture.go`: command runner and environment probe orchestration.
- `internal/capture/capture_test.go`: capture behavior tests using Go subprocess fixtures.
- `internal/checks/checks.go`: issue/PR readiness scoring.
- `internal/checks/checks_test.go`: readiness rule tests.
- `internal/mcpserver/server.go`: stdio MCP server tools.
- `internal/mcpserver/server_test.go`: tool handler tests without spawning a long-running MCP process.
- `.github/workflows/reprogate.yml`: dogfood workflow that runs tests.
- `action.yml`: GitHub Action metadata for `reprogate ready-check`.
- `README.md`: value proposition, install, CLI examples, MCP config, and Action example.

## Task 1: Module And CLI Skeleton

**Files:**
- Create: `go.mod`
- Create: `cmd/reprogate/main.go`

- [ ] **Step 1: Write the failing CLI smoke test command**

Run:

```bash
go test ./...
```

Expected: FAIL because the module and packages do not exist yet.

- [ ] **Step 2: Create module and CLI entrypoint**

Create `go.mod`:

```go
module github.com/JinRudy/reprogate

go 1.25

require github.com/modelcontextprotocol/go-sdk v0.9.0
```

Create `cmd/reprogate/main.go`:

```go
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
  reprogate mcp`)
}
```

- [ ] **Step 3: Run the test command**

Run:

```bash
go test ./...
```

Expected: FAIL because referenced internal packages are not implemented.

- [ ] **Step 4: Commit after green in later tasks**

Do not commit this task alone until the referenced packages compile.

## Task 2: Redaction Package

**Files:**
- Create: `internal/redact/redact_test.go`
- Create: `internal/redact/redact.go`

- [ ] **Step 1: Write failing redaction tests**

Create `internal/redact/redact_test.go`:

```go
package redact

import (
	"strings"
	"testing"
)

func TestTextRedactsCommonSecrets(t *testing.T) {
	input := "Authorization: Bearer abc123\npassword=my-secret\nOPENAI_API_KEY=sk-test\nurl=https://user:pass@example.com/db"
	got := Text(input)

	for _, secret := range []string{"abc123", "my-secret", "sk-test", "user:pass"} {
		if strings.Contains(got, secret) {
			t.Fatalf("expected %q to be redacted from %q", secret, got)
		}
	}
	for _, marker := range []string{"[REDACTED:bearer-token]", "[REDACTED:secret-value]", "[REDACTED:url-credentials]"} {
		if !strings.Contains(got, marker) {
			t.Fatalf("expected marker %q in %q", marker, got)
		}
	}
}

func TestTextRedactsHomePaths(t *testing.T) {
	input := "/Users/alice/projects/app failed and /home/bob/app failed"
	got := Text(input)
	if strings.Contains(got, "/Users/alice") || strings.Contains(got, "/home/bob") {
		t.Fatalf("expected home paths to be redacted, got %q", got)
	}
	if !strings.Contains(got, "[REDACTED:home-path]") {
		t.Fatalf("expected home path marker, got %q", got)
	}
}
```

- [ ] **Step 2: Verify red**

Run:

```bash
go test ./internal/redact
```

Expected: FAIL because `Text` is undefined.

- [ ] **Step 3: Implement minimal redaction**

Create `internal/redact/redact.go`:

```go
package redact

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var patterns = []struct {
	re   *regexp.Regexp
	repl string
}{
	{regexp.MustCompile(`(?i)Bearer\s+[A-Za-z0-9._~+/=-]+`), "Bearer [REDACTED:bearer-token]"},
	{regexp.MustCompile(`(?i)(password|passwd|pwd|secret|token|api[_-]?key|access[_-]?key)\s*[:=]\s*[^\\s]+`), "$1=[REDACTED:secret-value]"},
	{regexp.MustCompile(`https?://([^:/\s]+):([^@\s]+)@`), "https://[REDACTED:url-credentials]@"},
	{regexp.MustCompile(`/Users/[^/\s]+`), "[REDACTED:home-path]"},
	{regexp.MustCompile(`/home/[^/\s]+`), "[REDACTED:home-path]"},
}

func Text(input string) string {
	out := input
	for _, pattern := range patterns {
		out = pattern.re.ReplaceAllString(out, pattern.repl)
	}
	return out
}

func RunCLI(_ []string, in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		fmt.Fprintln(out, Text(scanner.Text()))
	}
	return scanner.Err()
}

func LooksRedacted(input string) bool {
	return strings.Contains(input, "[REDACTED:")
}
```

- [ ] **Step 4: Verify green**

Run:

```bash
go test ./internal/redact
```

Expected: PASS.

## Task 3: Report Renderer

**Files:**
- Create: `internal/report/report_test.go`
- Create: `internal/report/report.go`

- [ ] **Step 1: Write failing report tests**

Create `internal/report/report_test.go`:

```go
package report

import (
	"strings"
	"testing"
)

func TestRenderMarkdownIncludesCoreSections(t *testing.T) {
	doc := Document{
		Command: "npm test",
		ExitCode: 1,
		Duration: "2s",
		Environment: map[string]string{"os": "darwin", "node": "v22.0.0"},
		Logs: "FAIL app.test.ts",
	}

	got := RenderMarkdown(doc)
	for _, want := range []string{"# Reproduction Report", "## Command", "npm test", "## Environment", "node", "## Logs", "FAIL app.test.ts", "## Expected Behavior"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in markdown:\n%s", want, got)
		}
	}
}
```

- [ ] **Step 2: Verify red**

Run:

```bash
go test ./internal/report
```

Expected: FAIL because `Document` and `RenderMarkdown` are undefined.

- [ ] **Step 3: Implement renderer**

Create `internal/report/report.go`:

```go
package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Document struct {
	Command     string
	ExitCode    int
	Duration    string
	Environment map[string]string
	Dependency  map[string]string
	Logs        string
}

func RenderMarkdown(doc Document) string {
	var b strings.Builder
	fmt.Fprintln(&b, "# Reproduction Report")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Command")
	fmt.Fprintf(&b, "`%s`\n\n", doc.Command)
	fmt.Fprintf(&b, "- Exit code: `%d`\n", doc.ExitCode)
	if doc.Duration != "" {
		fmt.Fprintf(&b, "- Duration: `%s`\n", doc.Duration)
	}
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Expected Behavior")
	fmt.Fprintln(&b, "Describe what you expected to happen.")
	fmt.Fprintln(&b)
	fmt.Fprintln(&b, "## Actual Behavior")
	fmt.Fprintln(&b, "The command failed or produced the logs below.")
	fmt.Fprintln(&b)
	writeMap(&b, "Environment", doc.Environment)
	writeMap(&b, "Dependency State", doc.Dependency)
	fmt.Fprintln(&b, "## Logs")
	fmt.Fprintln(&b, "```text")
	fmt.Fprintln(&b, doc.Logs)
	fmt.Fprintln(&b, "```")
	return b.String()
}

func writeMap(b *strings.Builder, title string, values map[string]string) {
	fmt.Fprintf(b, "## %s\n", title)
	if len(values) == 0 {
		fmt.Fprintln(b, "- none detected")
		fmt.Fprintln(b)
		return
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(b, "- %s: `%s`\n", key, values[key])
	}
	fmt.Fprintln(b)
}

func RunCLI(_ []string, out io.Writer) error {
	_, err := io.WriteString(out, RenderMarkdown(Document{}))
	return err
}
```

- [ ] **Step 4: Verify green**

Run:

```bash
go test ./internal/report
```

Expected: PASS.

## Task 4: Capture Command

**Files:**
- Create: `internal/capture/capture_test.go`
- Create: `internal/capture/capture.go`

- [ ] **Step 1: Write failing capture test**

Create `internal/capture/capture_test.go`:

```go
package capture

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWritesSanitizedReport(t *testing.T) {
	dir := t.TempDir()
	result, err := Run(context.Background(), Options{
		Command: []string{"go", "env", "GOVERSION"},
		WorkDir: dir,
		OutputPath: filepath.Join(dir, ".reprogate", "repro.md"),
		MaxLogBytes: 200000,
	})
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if result.ReportPath == "" {
		t.Fatal("expected report path")
	}
	data, err := os.ReadFile(result.ReportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, "# Reproduction Report") || !strings.Contains(text, "go env GOVERSION") {
		t.Fatalf("unexpected report:\n%s", text)
	}
}
```

- [ ] **Step 2: Verify red**

Run:

```bash
go test ./internal/capture
```

Expected: FAIL because capture package is missing.

- [ ] **Step 3: Implement capture**

Create `internal/capture/capture.go`:

```go
package capture

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/JinRudy/reprogate/internal/redact"
	"github.com/JinRudy/reprogate/internal/report"
)

type Options struct {
	Command     []string
	WorkDir     string
	OutputPath  string
	MaxLogBytes int
}

type Result struct {
	ExitCode   int
	ReportPath string
}

func Run(ctx context.Context, opts Options) (Result, error) {
	if len(opts.Command) == 0 {
		return Result{}, fmt.Errorf("missing command")
	}
	if opts.WorkDir == "" {
		opts.WorkDir, _ = os.Getwd()
	}
	if opts.OutputPath == "" {
		opts.OutputPath = filepath.Join(opts.WorkDir, ".reprogate", "repro.md")
	}
	if opts.MaxLogBytes <= 0 {
		opts.MaxLogBytes = 200000
	}

	start := time.Now()
	cmd := exec.CommandContext(ctx, opts.Command[0], opts.Command[1:]...)
	cmd.Dir = opts.WorkDir
	var logs limitedBuffer
	logs.limit = opts.MaxLogBytes
	cmd.Stdout = &logs
	cmd.Stderr = &logs
	err := cmd.Run()
	exitCode := 0
	if err != nil {
		exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	doc := report.Document{
		Command: strings.Join(opts.Command, " "),
		ExitCode: exitCode,
		Duration: time.Since(start).Round(time.Millisecond).String(),
		Environment: probeEnvironment(),
		Dependency: probeDependency(opts.WorkDir),
		Logs: redact.Text(logs.String()),
	}

	if err := os.MkdirAll(filepath.Dir(opts.OutputPath), 0o755); err != nil {
		return Result{}, err
	}
	if err := os.WriteFile(opts.OutputPath, []byte(report.RenderMarkdown(doc)), 0o644); err != nil {
		return Result{}, err
	}
	return Result{ExitCode: exitCode, ReportPath: opts.OutputPath}, nil
}

func RunCLI(ctx context.Context, args []string, out io.Writer, _ io.Writer) error {
	if len(args) == 0 || args[0] != "--" || len(args) == 1 {
		return fmt.Errorf("usage: reprogate capture -- <command> [args...]")
	}
	result, err := Run(ctx, Options{Command: args[1:]})
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "wrote %s\n", result.ReportPath)
	if result.ExitCode != 0 {
		return fmt.Errorf("command exited with code %d", result.ExitCode)
	}
	return nil
}

type limitedBuffer struct {
	bytes.Buffer
	limit int
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.Len() < b.limit {
		remaining := b.limit - b.Len()
		if len(p) > remaining {
			b.Buffer.Write(p[:remaining])
		} else {
			b.Buffer.Write(p)
		}
	}
	return len(p), nil
}

func probeEnvironment() map[string]string {
	values := map[string]string{
		"os": runtime.GOOS,
		"arch": runtime.GOARCH,
	}
	if version, err := exec.Command("go", "version").Output(); err == nil {
		values["go"] = strings.TrimSpace(string(version))
	}
	return values
}

func probeDependency(workDir string) map[string]string {
	values := map[string]string{}
	for _, name := range []string{"go.sum", "package-lock.json", "pnpm-lock.yaml", "yarn.lock"} {
		path := filepath.Join(workDir, name)
		if data, err := os.ReadFile(path); err == nil {
			values[name] = fmt.Sprintf("%d bytes", len(data))
		}
	}
	return values
}
```

- [ ] **Step 4: Verify green**

Run:

```bash
go test ./internal/capture
```

Expected: PASS.

## Task 5: Readiness Checks

**Files:**
- Create: `internal/checks/checks_test.go`
- Create: `internal/checks/checks.go`

- [ ] **Step 1: Write failing readiness tests**

Create `internal/checks/checks_test.go`:

```go
package checks

import "testing"

func TestAnalyzeFlagsMissingEvidence(t *testing.T) {
	result := Analyze(Input{Body: "it fails"})
	if !result.HasLabel("needs-repro") || !result.HasLabel("missing-env") || !result.HasLabel("missing-log") {
		t.Fatalf("expected missing evidence labels, got %#v", result.Labels)
	}
}

func TestAnalyzeMarksReadyWhenEvidenceExists(t *testing.T) {
	body := "Steps to reproduce:\n1. run npm test\n\nEnvironment: macOS node 22\n\n```text\nFAIL app.test.ts\n```"
	result := Analyze(Input{Body: body, ChangedFiles: 2, ChangedLines: 20})
	if !result.HasLabel("review-ready") {
		t.Fatalf("expected review-ready, got %#v", result.Labels)
	}
}
```

- [ ] **Step 2: Verify red**

Run:

```bash
go test ./internal/checks
```

Expected: FAIL because `Analyze` is undefined.

- [ ] **Step 3: Implement checks**

Create `internal/checks/checks.go`:

```go
package checks

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Input struct {
	Body string
	ChangedFiles int
	ChangedLines int
	LockfileChanged bool
	GeneratedChanged bool
	DependencyChanged bool
}

type Result struct {
	Labels []string `json:"labels"`
	Missing []string `json:"missing"`
	Summary string `json:"summary"`
}

func (r Result) HasLabel(label string) bool {
	for _, existing := range r.Labels {
		if existing == label {
			return true
		}
	}
	return false
}

func Analyze(input Input) Result {
	body := strings.ToLower(input.Body)
	var labels []string
	var missing []string

	if !(strings.Contains(body, "steps to reproduce") || strings.Contains(body, "1.") || strings.Contains(body, "repro")) {
		labels = append(labels, "needs-repro")
		missing = append(missing, "reproduction steps")
	}
	if !(strings.Contains(body, "environment") || strings.Contains(body, "os:") || strings.Contains(body, "node") || strings.Contains(body, "go version")) {
		labels = append(labels, "missing-env")
		missing = append(missing, "environment details")
	}
	if !(strings.Contains(body, "```") || strings.Contains(body, "stack trace") || strings.Contains(body, "error:") || strings.Contains(body, "fail")) {
		labels = append(labels, "missing-log")
		missing = append(missing, "logs or command output")
	}
	if input.ChangedFiles > 20 || input.ChangedLines > 800 || input.LockfileChanged || input.GeneratedChanged || input.DependencyChanged {
		labels = append(labels, "risky-diff")
	}
	if len(missing) == 0 && !contains(labels, "risky-diff") {
		labels = append(labels, "review-ready")
	}
	return Result{Labels: labels, Missing: missing, Summary: renderSummary(labels, missing)}
}

func RunCLI(args []string, in io.Reader, out io.Writer) error {
	var body string
	if len(args) >= 2 && args[0] == "--issue-body" {
		data, err := os.ReadFile(args[1])
		if err != nil {
			return err
		}
		body = string(data)
	} else {
		var b strings.Builder
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			b.WriteString(scanner.Text())
			b.WriteByte('\n')
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		body = b.String()
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(Analyze(Input{Body: body}))
}

func renderSummary(labels []string, missing []string) string {
	return fmt.Sprintf("labels: %s; missing: %s", strings.Join(labels, ", "), strings.Join(missing, ", "))
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Verify green**

Run:

```bash
go test ./internal/checks
```

Expected: PASS.

## Task 6: MCP Server

**Files:**
- Create: `internal/mcpserver/server_test.go`
- Create: `internal/mcpserver/server.go`

- [ ] **Step 1: Write failing tool handler tests**

Create `internal/mcpserver/server_test.go`:

```go
package mcpserver

import (
	"strings"
	"testing"
)

func TestRedactTextTool(t *testing.T) {
	got := redactText("password=secret")
	if strings.Contains(got, "secret") || !strings.Contains(got, "[REDACTED:secret-value]") {
		t.Fatalf("unexpected redaction: %q", got)
	}
}

func TestCheckIssueTool(t *testing.T) {
	got := checkIssue("it fails")
	if !strings.Contains(got, "needs-repro") {
		t.Fatalf("expected needs-repro, got %q", got)
	}
}
```

- [ ] **Step 2: Verify red**

Run:

```bash
go test ./internal/mcpserver
```

Expected: FAIL because tool helpers are undefined.

- [ ] **Step 3: Implement MCP server wrapper**

Create `internal/mcpserver/server.go`:

```go
package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/JinRudy/reprogate/internal/checks"
	"github.com/JinRudy/reprogate/internal/redact"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Run(ctx context.Context, in io.Reader, out io.Writer) error {
	server := mcp.NewServer(&mcp.Implementation{Name: "reprogate", Version: "0.1.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name: "redact_text",
		Description: "Redact likely secrets, credentials, and private paths from text.",
	}, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[redactTextArgs]) (*mcp.CallToolResultFor[any], error) {
		return textResult(redactText(params.Arguments.Text)), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name: "check_issue",
		Description: "Check whether an issue or PR body has reproduction steps, environment details, and logs.",
	}, func(_ context.Context, _ *mcp.ServerSession, params *mcp.CallToolParamsFor[checkIssueArgs]) (*mcp.CallToolResultFor[any], error) {
		return textResult(checkIssue(params.Arguments.Body)), nil
	})

	transport := mcp.NewStdioTransport(in, out)
	return server.Run(ctx, transport)
}

type redactTextArgs struct {
	Text string `json:"text" jsonschema:"text to redact"`
}

type checkIssueArgs struct {
	Body string `json:"body" jsonschema:"issue or pull request body"`
}

func redactText(input string) string {
	return redact.Text(input)
}

func checkIssue(body string) string {
	result := checks.Analyze(checks.Input{Body: body})
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error":%q}`, err.Error())
	}
	return string(data)
}

func textResult(text string) *mcp.CallToolResultFor[any] {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
```

- [ ] **Step 4: Verify green**

Run:

```bash
go test ./internal/mcpserver
```

Expected: PASS.

## Task 7: Repository Docs And Action Metadata

**Files:**
- Create: `README.md`
- Create: `action.yml`
- Create: `.github/workflows/reprogate.yml`

- [ ] **Step 1: Write README and metadata**

Create `README.md` with install, CLI, MCP, and GitHub Action examples. Include this MCP client snippet:

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

Create `action.yml`:

```yaml
name: ReproGate Ready Check
description: Check issue and pull request readiness signals.
runs:
  using: composite
  steps:
    - shell: bash
      run: |
        go run ./cmd/reprogate ready-check < "$GITHUB_EVENT_PATH"
```

Create `.github/workflows/reprogate.yml`:

```yaml
name: reprogate
on:
  push:
  pull_request:
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: go test ./...
```

- [ ] **Step 2: Verify docs mention key commands**

Run:

```bash
rg -n "reprogate capture|reprogate mcp|ready-check|mcpServers" README.md action.yml .github/workflows/reprogate.yml
```

Expected: each key command appears.

## Task 8: Full Verification And Commit

**Files:**
- Modify: all touched files as needed for compile and test fixes.

- [ ] **Step 1: Run formatting**

Run:

```bash
gofmt -w cmd internal
```

Expected: no output.

- [ ] **Step 2: Run tests**

Run:

```bash
go test ./...
```

Expected: PASS.

- [ ] **Step 3: Run CLI smoke checks**

Run:

```bash
go run ./cmd/reprogate help
printf 'password=secret\n' | go run ./cmd/reprogate redact
go run ./cmd/reprogate capture -- go env GOVERSION
go run ./cmd/reprogate ready-check < .reprogate/repro.md
```

Expected:
- help prints usage.
- redaction output contains `[REDACTED:secret-value]`.
- capture writes `.reprogate/repro.md`.
- ready-check prints JSON labels.

- [ ] **Step 4: Commit and push**

Run:

```bash
git add .
git commit -m "feat(reprogate): 实现复现报告 MVP" -m $'- 新增 Go CLI 的采集、脱敏和报告能力\n- 新增 issue/PR readiness 检查和 MCP stdio 入口\n- 补充 README、Action 元数据和 CI 工作流'
git push
```

Expected: commit succeeds and pushes to `origin/main`.
