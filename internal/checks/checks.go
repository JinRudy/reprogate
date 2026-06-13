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
	body, err := readBody(args, in)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(Analyze(Input{Body: body}))
}

func readBody(args []string, in io.Reader) (string, error) {
	if len(args) >= 2 && args[0] == "--issue-body" {
		data, err := os.ReadFile(args[1])
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
