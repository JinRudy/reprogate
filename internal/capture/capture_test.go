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
		Command:     []string{"go", "env", "GOVERSION"},
		WorkDir:     dir,
		OutputPath:  filepath.Join(dir, ".reprogate", "repro.md"),
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
