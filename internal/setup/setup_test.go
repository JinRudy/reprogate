package setup

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCLIWritesGitHubActionWorkflow(t *testing.T) {
	target := filepath.Join(t.TempDir(), ".github", "workflows", "reprogate.yml")
	var out bytes.Buffer

	if err := RunCLI([]string{"github-action", "--path", target}, &out); err != nil {
		t.Fatalf("run cli: %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	workflow := string(data)
	for _, want := range []string{
		"name: reprogate",
		"issues:",
		"pull_request:",
		"uses: JinRudy/reprogate@v0.1.4",
	} {
		if !strings.Contains(workflow, want) {
			t.Fatalf("expected workflow to contain %q, got %s", want, workflow)
		}
	}
	if !strings.Contains(out.String(), "Wrote "+target) {
		t.Fatalf("expected output to mention written path, got %q", out.String())
	}
}

func TestRunCLIRefusesToOverwriteExistingWorkflow(t *testing.T) {
	target := filepath.Join(t.TempDir(), "reprogate.yml")
	if err := os.WriteFile(target, []byte("existing"), 0o600); err != nil {
		t.Fatalf("write existing workflow: %v", err)
	}

	err := RunCLI([]string{"github-action", "--path", target}, ioDiscard{})
	if err == nil {
		t.Fatal("expected overwrite error")
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	if string(data) != "existing" {
		t.Fatalf("expected existing workflow to stay unchanged, got %q", string(data))
	}
}

func TestRunCLIForceOverwritesExistingWorkflow(t *testing.T) {
	target := filepath.Join(t.TempDir(), "reprogate.yml")
	if err := os.WriteFile(target, []byte("existing"), 0o600); err != nil {
		t.Fatalf("write existing workflow: %v", err)
	}

	if err := RunCLI([]string{"github-action", "--path", target, "--force"}, ioDiscard{}); err != nil {
		t.Fatalf("run cli: %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	if !strings.Contains(string(data), "uses: JinRudy/reprogate@v0.1.4") {
		t.Fatalf("expected workflow to be overwritten, got %s", string(data))
	}
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
