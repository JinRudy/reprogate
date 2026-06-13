package mcpserver

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRedactTextTool(t *testing.T) {
	got := redactText("password=hunter2")
	if strings.Contains(got, "hunter2") || !strings.Contains(got, "[REDACTED:secret-value]") {
		t.Fatalf("unexpected redaction: %q", got)
	}
}

func TestCheckIssueTool(t *testing.T) {
	got := checkIssue("it fails")
	if !strings.Contains(got, "needs-repro") {
		t.Fatalf("expected needs-repro, got %q", got)
	}
}

func TestServerExposesToolsOverMCP(t *testing.T) {
	ctx := context.Background()
	server := New()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverTransport, nil); err != nil {
		t.Fatalf("connect server: %v", err)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0.1.0"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("connect client: %v", err)
	}
	defer session.Close()

	seen := map[string]bool{}
	for tool, err := range session.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("list tools: %v", err)
		}
		seen[tool.Name] = true
	}
	for _, name := range []string{"redact_text", "check_issue"} {
		if !seen[name] {
			t.Fatalf("expected MCP tool %q, got %#v", name, seen)
		}
	}

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "redact_text",
		Arguments: map[string]any{"text": "password=hunter2"},
	})
	if err != nil {
		t.Fatalf("call redact_text: %v", err)
	}
	text := result.Content[0].(*mcp.TextContent).Text
	if strings.Contains(text, "hunter2") || !strings.Contains(text, "[REDACTED:secret-value]") {
		t.Fatalf("unexpected MCP tool result: %q", text)
	}
}
