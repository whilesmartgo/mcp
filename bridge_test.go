package agentsmcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/whilesmartgo/agents"
)

func connect(t *testing.T, server *mcp.Server) *mcp.ClientSession {
	t.Helper()
	ctx := context.Background()
	serverT, clientT := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, serverT, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0"}, nil)
	cs, err := client.Connect(ctx, clientT, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = cs.Close() })
	return cs
}

func echoTool() agents.Tool {
	return agents.Tool{
		Name:        "echo",
		Description: "echo the text",
		Parameters: map[string]any{
			"type":       "object",
			"properties": map[string]any{"text": map[string]any{"type": "string"}},
		},
		Handler: func(_ context.Context, args json.RawMessage) (string, error) {
			var in struct {
				Text string `json:"text"`
			}
			_ = json.Unmarshal(args, &in)
			return "echo:" + in.Text, nil
		},
	}
}

func TestServerExposesRegistryTools(t *testing.T) {
	cs := connect(t, NewServer("test", "v0", agents.NewRegistry(echoTool())))
	ctx := context.Background()

	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "echo" {
		t.Fatalf("expected the echo tool advertised, got %+v", tools.Tools)
	}

	res, err := cs.CallTool(ctx, &mcp.CallToolParams{Name: "echo", Arguments: map[string]any{"text": "hi"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.IsError {
		t.Fatalf("unexpected error result: %+v", res)
	}
	if text := res.Content[0].(*mcp.TextContent).Text; text != "echo:hi" {
		t.Errorf("tool output = %q", text)
	}
}

func TestToolErrorBecomesErrorResult(t *testing.T) {
	failing := agents.Tool{
		Name:       "boom",
		Parameters: map[string]any{"type": "object"},
		Handler:    func(context.Context, json.RawMessage) (string, error) { return "", errors.New("nope") },
	}
	cs := connect(t, NewServer("test", "v0", agents.NewRegistry(failing)))

	// A tool error must surface as an error result, not a protocol-level error,
	// so a model on the other end can see and adapt to it.
	res, err := cs.CallTool(context.Background(), &mcp.CallToolParams{Name: "boom"})
	if err != nil {
		t.Fatalf("a tool error should not be a protocol error: %v", err)
	}
	if !res.IsError {
		t.Fatal("a tool error should be reported as an error result")
	}
	if text := res.Content[0].(*mcp.TextContent).Text; text != "nope" {
		t.Errorf("error text = %q", text)
	}
}
