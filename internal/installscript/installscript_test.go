package installscript

import (
	"os"
	"strings"
	"testing"
)

func TestInstallScriptCurlHasTimeoutsAndRetry(t *testing.T) {
	data, err := os.ReadFile("../../scripts/install.sh")
	if err != nil {
		t.Fatalf("read install script: %v", err)
	}
	script := string(data)

	for _, want := range []string{
		"--connect-timeout",
		"--max-time",
		"--retry",
	} {
		if !strings.Contains(script, want) {
			t.Fatalf("expected install script curl command to contain %q", want)
		}
	}
}
