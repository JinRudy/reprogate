package report

import (
	"strings"
	"testing"
)

func TestRenderMarkdownIncludesCoreSections(t *testing.T) {
	doc := Document{
		Command:  "npm test",
		ExitCode: 1,
		Duration: "2s",
		Environment: map[string]string{
			"node": "v22.0.0",
			"os":   "darwin",
		},
		Logs: "FAIL app.test.ts",
	}

	got := RenderMarkdown(doc)
	for _, want := range []string{"# Reproduction Report", "## Command", "npm test", "## Environment", "node", "## Logs", "FAIL app.test.ts", "## Expected Behavior"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q in markdown:\n%s", want, got)
		}
	}
}
