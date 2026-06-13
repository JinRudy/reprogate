package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/JinRudy/reprogate/internal/checks"
	"github.com/JinRudy/reprogate/internal/redact"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Run(ctx context.Context, in io.Reader, out io.Writer) error {
	server := New()
	transport := &mcp.IOTransport{
		Reader: io.NopCloser(in),
		Writer: nopWriteCloser{Writer: out},
	}
	return server.Run(ctx, transport)
}

func New() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "reprogate", Version: "0.1.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "redact_text",
		Description: "Redact likely secrets, credentials, and private paths from text.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, input redactTextArgs) (*mcp.CallToolResult, any, error) {
		return textResult(redactText(input.Text)), nil, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "check_issue",
		Description: "Check whether an issue or pull request body has reproduction steps, environment details, and logs.",
	}, func(_ context.Context, _ *mcp.CallToolRequest, input checkIssueArgs) (*mcp.CallToolResult, any, error) {
		return textResult(checkIssue(input.Body)), nil, nil
	})

	return server
}

type redactTextArgs struct {
	Text string `json:"text" jsonschema:"text to redact"`
}

type checkIssueArgs struct {
	Body string `json:"body" jsonschema:"issue or pull request body"`
}

type nopWriteCloser struct {
	io.Writer
}

func (w nopWriteCloser) Close() error {
	return nil
}

func redactText(input string) string {
	return redact.Text(input)
}

func checkIssue(body string) string {
	result := checks.Analyze(checks.Input{Body: body})
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error":%q}`, err.Error())
	}
	return string(data)
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
