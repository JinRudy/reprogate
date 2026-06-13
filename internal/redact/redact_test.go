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
