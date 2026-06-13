package actionmetadata

import (
	"os"
	"strings"
	"testing"
)

func TestActionExposesCommentOnMissingInput(t *testing.T) {
	action := readAction(t)

	for _, want := range []string{
		"comment-on-missing:",
		"Post or update one issue or pull request comment when reproduction evidence is missing.",
		"default: \"false\"",
	} {
		if !strings.Contains(action, want) {
			t.Fatalf("expected action metadata to contain %q", want)
		}
	}
}

func TestActionCommentOnMissingUsesSingleGuardedIssueComment(t *testing.T) {
	action := readAction(t)

	for _, want := range []string{
		"if [ \"${{ inputs.comment-on-missing }}\" != \"true\" ]; then",
		"issue_number=\"$(jq -r '.issue.number // .pull_request.number // empty' \"$event_path\")\"",
		"<!-- reprogate-missing-evidence -->",
		"issues/comments/$previous_id",
	} {
		if !strings.Contains(action, want) {
			t.Fatalf("expected action metadata to contain %q", want)
		}
	}
}

func TestActionCommentHereDocStaysInsideYAMLBlock(t *testing.T) {
	action := readAction(t)

	for _, want := range []string{
		"\n        <!-- reprogate-missing-evidence -->",
		"\n        EOF\n",
		"\n        )\"\n",
	} {
		if !strings.Contains(action, want) {
			t.Fatalf("expected action metadata to contain indented here-doc line %q", want)
		}
	}
}

func readAction(t *testing.T) string {
	t.Helper()

	data, err := os.ReadFile("../../action.yml")
	if err != nil {
		t.Fatalf("read action metadata: %v", err)
	}
	return string(data)
}
