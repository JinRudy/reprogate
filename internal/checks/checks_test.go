package checks

import (
	"bytes"
	"strings"
	"testing"
)

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

func TestRunCLIEmitsEmptyMissingArrayWhenReady(t *testing.T) {
	body := "Steps to reproduce:\n1. run npm test\n\nEnvironment: macOS node 22\n\n```text\nFAIL app.test.ts\n```"
	var out bytes.Buffer
	if err := RunCLI(nil, strings.NewReader(body), &out); err != nil {
		t.Fatalf("run cli: %v", err)
	}
	if !strings.Contains(out.String(), `"missing": []`) {
		t.Fatalf("expected empty missing array, got %s", out.String())
	}
}

func TestRunCLIFailsWhenMissingEvidenceIsRequired(t *testing.T) {
	var out bytes.Buffer
	err := RunCLI([]string{"--fail-on-missing"}, strings.NewReader("it fails"), &out)
	if err == nil {
		t.Fatal("expected error when required evidence is missing")
	}
	if !strings.Contains(err.Error(), "missing required evidence") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), `"needs-repro"`) {
		t.Fatalf("expected JSON output before failure, got %s", out.String())
	}
}

func TestRunCLIDoesNotFailWhenRequiredEvidenceExists(t *testing.T) {
	body := "Steps to reproduce:\n1. run npm test\n\nEnvironment: macOS node 22\n\n```text\nFAIL app.test.ts\n```"
	var out bytes.Buffer
	if err := RunCLI([]string{"--fail-on-missing"}, strings.NewReader(body), &out); err != nil {
		t.Fatalf("expected ready issue to pass: %v", err)
	}
}
