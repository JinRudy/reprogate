package checks

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Input struct {
	Body              string
	ChangedFiles      int
	ChangedLines      int
	LockfileChanged   bool
	GeneratedChanged  bool
	DependencyChanged bool
}

type Result struct {
	Labels  []string `json:"labels"`
	Missing []string `json:"missing"`
	Summary string   `json:"summary"`
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
	labels := []string{}
	missing := []string{}

	if !(strings.Contains(body, "steps to reproduce") || strings.Contains(body, "1.") || strings.Contains(body, "repro")) {
		labels = append(labels, "needs-repro")
		missing = append(missing, "reproduction steps")
	}
	if !(strings.Contains(body, "environment") || strings.Contains(body, "os:") || strings.Contains(body, "node") || strings.Contains(body, "go version")) {
		labels = append(labels, "missing-env")
		missing = append(missing, "environment details")
	}
	if !(strings.Contains(body, "```") || strings.Contains(body, "stack trace") || strings.Contains(body, "error:") || strings.Contains(body, "panic:") || strings.Contains(body, "traceback")) {
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
	opts, err := parseOptions(args)
	if err != nil {
		return err
	}
	body, err := readBody(opts.issueBodyPath, in)
	if err != nil {
		return err
	}
	result := Analyze(Input{Body: body})
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return err
	}
	if opts.failOnMissing && len(result.Missing) > 0 {
		return fmt.Errorf("missing required evidence: %s", strings.Join(result.Missing, ", "))
	}
	return nil
}

type options struct {
	issueBodyPath string
	failOnMissing bool
}

func parseOptions(args []string) (options, error) {
	var opts options
	fs := flag.NewFlagSet("ready-check", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.issueBodyPath, "issue-body", "", "path to issue or pull request body")
	fs.BoolVar(&opts.failOnMissing, "fail-on-missing", false, "exit non-zero when reproduction evidence is missing")
	if err := fs.Parse(args); err != nil {
		return options{}, err
	}
	return opts, nil
}

func readBody(issueBodyPath string, in io.Reader) (string, error) {
	if issueBodyPath != "" {
		data, err := os.ReadFile(issueBodyPath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	var b strings.Builder
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		b.WriteString(scanner.Text())
		b.WriteByte('\n')
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return extractBody(b.String()), nil
}

func extractBody(input string) string {
	var event struct {
		Issue struct {
			Body string `json:"body"`
		} `json:"issue"`
		PullRequest struct {
			Body string `json:"body"`
		} `json:"pull_request"`
	}
	if json.Unmarshal([]byte(input), &event) == nil {
		if event.Issue.Body != "" {
			return event.Issue.Body
		}
		if event.PullRequest.Body != "" {
			return event.PullRequest.Body
		}
	}
	return input
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
