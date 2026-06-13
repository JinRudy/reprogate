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

func RunCLI(_ []string, out io.Writer) error {
	_, err := io.WriteString(out, RenderMarkdown(Document{}))
	return err
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
