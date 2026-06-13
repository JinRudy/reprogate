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
	{regexp.MustCompile(`(?i)(password|passwd|pwd|secret|token|api[_-]?key|access[_-]?key)\s*[:=]\s*[^\s]+`), "$1=[REDACTED:secret-value]"},
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
